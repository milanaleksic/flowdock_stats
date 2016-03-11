package main

import (
	"flag"
	"fmt"
	"time"

	cmd "github.com/milanaleksic/flowdock_stats/cmdcolors"
	"github.com/milanaleksic/flowdock_stats/flowdock"
)

// Version holds the main version string which should be updated externally when building release
var Version = "undefined"

func main() {
	days := flag.Int("days", 1, "number of days to look in the history")
	showVersion := flag.Bool("version", false, "Get application version")
	flowdockAPIToken := flag.String("flowdockApiToken", "", "Flowdock API token (from https://www.flowdock.com/account/tokens)")
	companyToAnalyze := flag.String("companyToAnalyze", "", "Company whose Flowdock flows are to be analyzed")
	flowToAnalyze := flag.String("flowToAnalyze", "", "Company's Flow to analyze")
	flag.Parse()

	if *showVersion {
		cmd.Info(fmt.Sprintf("flowdock_stats version: %v\n", Version))
		return
	}

	if *flowdockAPIToken == "" || *companyToAnalyze == "" || *flowToAnalyze == "" {
		cmd.Warn("flowdockAPIToken, companyToAnalyze and flowToAnalyze are mandatory program arguments. Please use --help to see command line help")
		return
	}

	mainContext := context{
		timeToLookInto:   time.Hour * 24 * time.Duration(*days),
		companyToAnalyze: *companyToAnalyze,
		flowToAnalyze:    *flowToAnalyze,
		api:              flowdock.API{APIToken: *flowdockAPIToken},
	}
	mainContext.fetchMessages()
	mainContext.enrichStatisticsWithRealUserNames()
	mainContext.presentStatistics()
}
