package main

import (
	"time"
	"flag"
	"fmt"
	"github.com/milanaleksic/flowdock_stats/flowdock"
)

func main() {
	days := flag.Int("days", 1, "number of days to look in the history")
	flowdockApiToken := flag.String("flowdockApiToken", "", "Flowdock API token")
	companyToAnalyze := flag.String("companyToAnalyze", "", "Company whose Flowdock flows are to be analyzed")
	flowToAnalyze := flag.String("flowToAnalyze", "", "Company's Flow to analyze")
	flag.Parse()

	if *flowdockApiToken == "" || *companyToAnalyze == "" || *flowToAnalyze == "" {
		fmt.Errorf("flowdockApiToken, companyToAnalyze and flowToAnalyze are mandatory program arguments. Please use --help to see them")
		return
	}

	context := Context{
		timeToLookInto: time.Hour * 24 * time.Duration(*days),
		companyToAnalyze: *companyToAnalyze,
		flowToAnalyze: *flowToAnalyze,
		api: flowdock.FlowdockApi{ApiToken:*flowdockApiToken},
	}
	context.fetchMessages()
	context.enrichStatisticsWithRealUserNames()
	context.presentStatistics()
}
