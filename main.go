package main

import (
	"flag"
	cmd "github.com/milanaleksic/flowdock_stats/cmdcolors"
	"github.com/milanaleksic/flowdock_stats/flowdock"
	"time"
)

func main() {
	days := flag.Int("days", 1, "number of days to look in the history")
	flowdockAPIToken := flag.String("flowdockApiToken", "", "Flowdock API token (from https://www.flowdock.com/account/tokens)")
	companyToAnalyze := flag.String("companyToAnalyze", "", "Company whose Flowdock flows are to be analyzed")
	flowToAnalyze := flag.String("flowToAnalyze", "", "Company's Flow to analyze")
	flag.Parse()

	if *flowdockAPIToken == "" || *companyToAnalyze == "" || *flowToAnalyze == "" {
		cmd.Warn("flowdockAPIToken, companyToAnalyze and flowToAnalyze are mandatory program arguments. Please use --help to see command line help")
		return
	}

	context := Context{
		timeToLookInto:   time.Hour * 24 * time.Duration(*days),
		companyToAnalyze: *companyToAnalyze,
		flowToAnalyze:    *flowToAnalyze,
		api:              flowdock.FlowdockApi{APIToken: *flowdockAPIToken},
	}
	context.fetchMessages()
	context.enrichStatisticsWithRealUserNames()
	context.presentStatistics()
}
