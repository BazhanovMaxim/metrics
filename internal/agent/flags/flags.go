package flags

import "flag"

var RunAddress string
var ReportInterval int
var PollInterval int

func ParseAgentFlags() {
	flag.StringVar(&RunAddress, "a", "localhost:8080", "address and port to run agent")
	flag.IntVar(&ReportInterval, "r", 10, "report interval")
	flag.IntVar(&PollInterval, "p", 2, "pollInterval interval")
	flag.Parse()
}
