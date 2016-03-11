package main

import (
	"fmt"
	"regexp"
	"sort"
	"time"

	cmd "github.com/milanaleksic/flowdock_stats/cmdcolors"
	"github.com/milanaleksic/flowdock_stats/flowdock"
	"github.com/milanaleksic/flowdock_stats/serialization"
)

var wordsRegex = regexp.MustCompile("\\w+")

type context struct {
	timeToLookInto   time.Duration
	companyToAnalyze string
	flowToAnalyze    string
	statistics       map[string]stat
	api              flowdock.API
}

func (context *context) calculatestatsForAMessage(msg flowdock.Message) {
	aStat, ok := context.statistics[msg.User]
	if !ok {
		aStat = stat{}
	}
	aStat.numberOfAppearances++
	if msg.Edited != 0 {
		aStat.numberOfEdits++
	}
	aStat.words += len(wordsRegex.FindAllString(string(msg.Content), -1))
	context.statistics[msg.User] = aStat
}

func (context *context) fetchMessages() {
	context.statistics = make(map[string]stat)
	var lastID = -1
	var begin = time.Now().Add(-context.timeToLookInto)
	for {
		messages, err := context.api.GetMessages(context.companyToAnalyze, context.flowToAnalyze, lastID)
		if err != nil {
			cmd.Warn(fmt.Sprintf("Giving up from messages fetching: %v", err))
			break
		}
		for _, msg := range messages {
			if msg.User != "0" {
				context.calculatestatsForAMessage(msg)
			}
		}
		if len(messages) > 0 {
			if messages[0].CreatedAt.Before(begin) {
				break
			} else {
				cmd.InfoInline(fmt.Sprintf("Fetching messages: %d%%", 100-100*(messages[0].CreatedAt.UnixNano()-begin.UnixNano())/context.timeToLookInto.Nanoseconds()))
			}
			lastID = messages[0].ID
		}
	}
	cmd.Info(fmt.Sprintf("%-30s", "Messages downloaded"))
}

func (context *context) enrichStatisticsWithRealUserNames() {
	catalog := serialization.GetKnownUsers()
	var count = len(context.statistics)
	var iter = 0
	for user, stat := range context.statistics {
		cmd.InfoInline(fmt.Sprintf("Fetching users: %d%%", 100*iter/count))
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

func (context *context) getUserName(catalog *map[string]*serialization.Catalog_User, userID string) (userNick string, err error) {
	if knownUser, ok := (*catalog)[userID]; ok {
		userNick = knownUser.Username
		return
	}
	userNick, err = context.api.GetUser(userID)
	if err != nil {
		return
	}
	catalogUser := serialization.Catalog_User{UserId: userID, Username: userNick}
	(*catalog)[userID] = &catalogUser
	return
}

func (context *context) asSortedStatsOnly() (statsOnly []stat) {
	statsOnly = make([]stat, 0)
	for _, p := range context.statistics {
		statsOnly = append(statsOnly, p)
	}
	sort.Sort(statByNumberOfAppearances(statsOnly))
	return statsOnly
}

func (context *context) presentStatistics() {
	statsOnly := context.asSortedStatsOnly()
	cmd.Info("\nProcessing finished, statistics are: ")
	for _, stat := range statsOnly {
		cmd.Info(fmt.Sprintf("%v has made %d comments, %d words (%d words per comment), %d%% corrected",
			stat.name,
			stat.numberOfAppearances,
			stat.words,
			stat.words/stat.numberOfAppearances,
			100*stat.numberOfEdits/stat.numberOfAppearances))
	}
}
