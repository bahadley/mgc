package config

import (
	"flag"
	"os"
	"strings"

	"github.com/bahadley/mgc/log"
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

	dsts          string
	dstPort       string
	numHeartbeats int
)

func IsLeader() bool {
	return Role == leaderFlag
}

func DstPort() string {
	return dstPort
}

func NumHeartbeats() int {
	return numHeartbeats
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
	flag.StringVar(&dsts, "dsts", defaultDstAddr, "Peer IP addresses")
	flag.StringVar(&dstPort, "port", defaultDstPort, "Peer port number")
	numHeartbeats = *(flag.Int("hbts", defaultNumHeartbeats, "Number of heartbeats"))
	flag.Parse()
	validateAll()
}

func validateAll() {
	parseDsts()
	validateNumHeartbeats()
}

func parseDsts() {
	DstAddrs = strings.Split(dsts, ",")
}

func validateNumHeartbeats() {
	if numHeartbeats <= 0 {
		log.Error.Fatalf("Invalid environment variable value: %s",
			"hbts")
	}
}
