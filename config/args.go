package config

import (
	"flag"
	"os"
	"strings"
)

const (
	envDstPort       = "MGC_DST_PORT"
	envNumHeartbeats = "MGC_NUM_HEARTBEATS"
	envDelayInt      = "MGC_DELAY_INTERVAL"
	envTransmit      = "MGC_TRANSMIT"
	envTrace         = "MGC_TRACE"

	leaderFlag     = "L"
	followerFlog   = "F"
	traceFlag      = "YES"
	noTransmitFlag = "NO"

	defaultAddr          = "localhost"
	defaultDstAddr       = "localhost"
	defaultDstPort       = "22221"
	defaultNumHeartbeats = 10
	defaultDelayInt      = 1000
)

var (
	Role     string
	Addr     string
	DstAddrs []string

	dsts string
)

func IsLeader() bool {
	return Role == leaderFlag
}

func Trace() bool {
	t := os.Getenv(envTrace)
	if len(t) > 0 && strings.ToUpper(t) == traceFlag {
		return true
	} else {
		return false
	}
}

func init() {
	flag.StringVar(&Role, "role", leaderFlag, "Node role [L,F]")
	flag.StringVar(&Addr, "addr", defaultAddr, "Node IP address")
	flag.StringVar(&dsts, "dsts", defaultDstAddr, "Destination IP addresses")
	flag.Parse()
	validate()
}

func validate() {
	parseDsts()
}

func parseDsts() {
	DstAddrs = strings.Split(dsts, ",")
}
