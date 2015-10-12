package main

import (
	"flag"
	cmd "github.com/milanaleksic/flowdock_stats/cmd_colors"
	"github.com/milanaleksic/flowdock_stats/flowdock"
	"time"
)

func main() {
	days := flag.Int("days", 1, "number of days to look in the history")
	flowdockApiToken := flag.String("flowdockApiToken", "", "Flowdock API token")
	companyToAnalyze := flag.String("companyToAnalyze", "", "Company whose Flowdock flows are to be analyzed")
	flowToAnalyze := flag.String("flowToAnalyze", "", "Company's Flow to analyze")
	flag.Parse()

	if *flowdockApiToken == "" || *companyToAnalyze == "" || *flowToAnalyze == "" {
		cmd.Warn("flowdockApiToken, companyToAnalyze and flowToAnalyze are mandatory program arguments. Please use --help to see command line help")
		return
	}

	context := Context{
		timeToLookInto:   time.Hour * 24 * time.Duration(*days),
		companyToAnalyze: *companyToAnalyze,
		flowToAnalyze:    *flowToAnalyze,
		api:              flowdock.FlowdockApi{ApiToken: *flowdockApiToken},
	}
	context.fetchMessages()
	context.enrichStatisticsWithRealUserNames()
	context.presentStatistics()
}
