package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"
	"sort"
	"github.com/milanaleksic/flowdock_stats/serialization"
	stats "github.com/milanaleksic/flowdock_stats/cmd_colors"
)

var wordsRegex = regexp.MustCompile("\\w+")

type Context struct {
	flowdockApiToken string
	timeToLookInto time.Duration
	companyToAnalyze string
	flowToAnalyze string
	statistics map[string]Stat
}

func printOutContents(msg Message) {
	var content string
	err := json.Unmarshal(msg.Content, &content)
	if err != nil {
		stats.Warn(fmt.Sprint("Could not parse contents: ", err))
	}
	stats.Info(fmt.Sprintf(msg.User, " said: ", content, " id=", msg.Id, " at ", msg.CreatedAt))
}

func (context *Context) processMessage(msg Message) {
	stat, ok := context.statistics[msg.User]
	if !ok {
		stat = Stat{}
	}
	stat.numberOfAppearances += 1
	if msg.Edited != 0 {
		stat.numberOfEdits += 1
	}
	stat.words += len(wordsRegex.FindAllString(string(msg.Content), -1))
	context.statistics[msg.User] = stat
	//printOutContents(msg)
}

func (context *Context) fetchMessages() {
	context.statistics = make(map[string]Stat)
	flowLocation := fmt.Sprintf("https://api.flowdock.com/flows/%s/%s/messages", context.companyToAnalyze, context.flowToAnalyze)
	var lastId int = -1
	var begin = time.Now().Add(-context.timeToLookInto)
	for {
		request, err := http.NewRequest("GET", flowLocation, nil)
		if err != nil {
			stats.Warn(fmt.Sprintf("Error encountered while fetching: %v", err))
			break
		}

		request.SetBasicAuth(context.flowdockApiToken, "")
		values := request.URL.Query()
		if lastId != -1 {
			values.Add("until_id", strconv.Itoa(lastId))
		}
		values.Add("limit", "100")
		request.URL.RawQuery = values.Encode()

		client := &http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			stats.Warn(fmt.Sprintf("Error encountered while fetching: %v", err))
			break
		}
		defer resp.Body.Close()
		messages := make([]Message, 100)
		err = json.NewDecoder(resp.Body).Decode(&messages)
		if err != nil {
			stats.Warn(fmt.Sprintf("Error encountered while parsing: %v", err))
			break
		}
		if len(messages) > 0 {
			if messages[0].CreatedAt.Before(begin) {
				break
			} else {
				stats.InfoInline(fmt.Sprintf("Fetching messages: %d%%", int64(100-100*(messages[0].CreatedAt.UnixNano()-begin.UnixNano())/context.timeToLookInto.Nanoseconds())))
			}
			lastId = messages[0].Id
		}
		for _, msg := range messages {
			if msg.User != "0" {
				context.processMessage(msg)
			}
		}
	}
	stats.Info(fmt.Sprintf("%-30s", "Messages downloaded"))
}

func (context *Context) enrichStatisticsWithRealUserNames() {
	catalog := serialization.GetKnownUsers()
	var count = len(context.statistics)
	var iter = 0
	for user, stat := range context.statistics {
		stats.InfoInline(fmt.Sprintf("Fetching users: %d%%", 100*iter/count))
		result,err := context.getUserName(&catalog.Users, user)
		if err != nil {
			stats.Warn(fmt.Sprintf("Error encountered while trying to get user name: %v", err))
			continue
		}
		stat.name = result.Nick
		context.statistics[user] = stat
		iter++
	}
	serialization.SaveUsers(catalog)
	stats.Info(fmt.Sprintf("%-30s", "Users downloaded"))
}

func (context *Context) getUserName(catalog *map[string]*serialization.Catalog_User, userId string) (user User, err error) {
	if knownUser, ok := (*catalog)[userId]; ok {
		user = User{ Nick: knownUser.Username }
		return
	}
	request, err := http.NewRequest("GET", "https://api.flowdock.com/users/"+userId, nil)
	if err != nil {
		stats.Warn(fmt.Sprintf("Error encountered while fetching: %v", err))
		return
	}
	request.SetBasicAuth(context.flowdockApiToken, "")
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		stats.Warn(fmt.Sprintf("Error encountered while fetching: %v", err))
		return
	}
	defer resp.Body.Close()
	result := User{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		stats.Warn(fmt.Sprintf("Error encountered while parsing: %v", err))
		return
	}
	user = result
	catalogUser := serialization.Catalog_User{UserId:userId, Username:result.Nick}
	(*catalog)[userId] = &catalogUser
	return
}

func (context *Context) asSortedStatsOnly() (statsOnly []Stat) {
	statsOnly = make([]Stat, 0)
	for _, p := range context.statistics {
		statsOnly = append(statsOnly, p)
	}
	sort.Sort(StatByNumberOfAppearances(statsOnly))
	return statsOnly
}

func (context *Context) presentStatistics() {
	statsOnly := context.asSortedStatsOnly()
	stats.Info("\nProcessing finished, statistics are: ")
	for _, stat := range statsOnly {
		stats.Info(fmt.Sprintf("%v has made %d comments, %d words (%d words per comment), %d%% corrected",
			stat.name,
			stat.numberOfAppearances,
			stat.words,
			stat.words/stat.numberOfAppearances,
			100*stat.numberOfEdits/stat.numberOfAppearances))
	}
}