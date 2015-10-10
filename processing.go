package main

import (
	"fmt"
	"regexp"
	"time"
	"sort"
	"github.com/milanaleksic/flowdock_stats/serialization"
	cmd "github.com/milanaleksic/flowdock_stats/cmd_colors"
	"github.com/milanaleksic/flowdock_stats/flowdock"
)

var wordsRegex = regexp.MustCompile("\\w+")

type Context struct {
	timeToLookInto   time.Duration
	companyToAnalyze string
	flowToAnalyze    string
	statistics       map[string]Stat
	api              flowdock.FlowdockApi
}

func (context *Context) calculateStatsForAMessage(msg flowdock.Message) {
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
	var lastId int = -1
	var begin = time.Now().Add(-context.timeToLookInto)
	for {
		messages, err := context.api.GetMessages(context.companyToAnalyze, context.flowToAnalyze, lastId)
		if err != nil {
			cmd.Warn(fmt.Sprintf("Giving up from messages fetching: %v", err))
			break
		}
		for _, msg := range messages {
			if msg.User != "0" {
				context.calculateStatsForAMessage(msg)
			}
		}
		if len(messages) > 0 {
			if messages[0].CreatedAt.Before(begin) {
				break
			} else {
				cmd.InfoInline(fmt.Sprintf("Fetching messages: %d%%", int64(100 - 100 * (messages[0].CreatedAt.UnixNano() - begin.UnixNano()) / context.timeToLookInto.Nanoseconds())))
			}
			lastId = messages[0].Id
		}
	}
	cmd.Info(fmt.Sprintf("%-30s", "Messages downloaded"))
}

func (context *Context) enrichStatisticsWithRealUserNames() {
	catalog := serialization.GetKnownUsers()
	var count = len(context.statistics)
	var iter = 0
	for user, stat := range context.statistics {
		cmd.InfoInline(fmt.Sprintf("Fetching users: %d%%", 100 * iter / count))
		result, err := context.getUserName(&catalog.Users, user)
		if err != nil {
			cmd.Warn(fmt.Sprintf("Error encountered while trying to get user name: %v", err))
			continue
		}
		stat.name = result
		context.statistics[user] = stat
		iter++
	}
	serialization.SaveUsers(catalog)
	cmd.Info(fmt.Sprintf("%-30s", "Users downloaded"))
}

func (context *Context) getUserName(catalog *map[string]*serialization.Catalog_User, userId string) (userNick string, err error) {
	if knownUser, ok := (*catalog)[userId]; ok {
		userNick = knownUser.Username
		return
	}
	userNick, err = context.api.GetUser(userId)
	if err != nil {
		return
	}
	catalogUser := serialization.Catalog_User{UserId:userId, Username: userNick}
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
	cmd.Info("\nProcessing finished, statistics are: ")
	for _, stat := range statsOnly {
		cmd.Info(fmt.Sprintf("%v has made %d comments, %d words (%d words per comment), %d%% corrected",
			stat.name,
			stat.numberOfAppearances,
			stat.words,
			stat.words / stat.numberOfAppearances,
			100 * stat.numberOfEdits / stat.numberOfAppearances))
	}
}