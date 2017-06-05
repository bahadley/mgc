package config

import (
	"flag"
	"strings"
	"time"

	"github.com/bahadley/mgc/log"
)

const (
	flagRole          = "role"
	flagAddr          = "addr"
	flagDsts          = "dsts"
	flagPort          = "port"
	flagNumHeartbeats = "hbts"
	flagDelayInt      = "hbdelay"
	flagStart         = "start"
	flagTrace         = "trace"

	leaderFlag   = "L"
	followerFlag = "F"

	defaultAddr          = "localhost"
	defaultDstAddr       = "localhost"
	defaultPort          = "22221"
	defaultNumHeartbeats = 10
	defaultDelayInt      = 1000
	defaultStart         = 0
	defaultTrace         = true
)

var (
	role          *string
	addr          *string
	dstAddrs      *string
	port          *string
	numHeartbeats *int
	delayInt      *int
	start         *int64

	trace *bool
)

func IsLeader() bool {
	return *role == leaderFlag
}

func IsFollower() bool {
	return *role == followerFlag
}

func Addr() string {
	return *addr
}

func DstAddrs() []string {
	return strings.Split(*dstAddrs, ",")
}

func Port() string {
	return *port
}

func NumHeartbeats() int {
	return *numHeartbeats
}

func DelayInterval() time.Duration {
	return time.Duration(*delayInt)
}

func Start() time.Time {
	return time.Unix(*start, 0)
}

func init() {
	setFlags()
	flag.Parse()
	validateAll()
	log.SetTrace(*trace)
}

func setFlags() {
	role = flag.String(flagRole, leaderFlag, "Node role [L,F]")
	addr = flag.String(flagAddr, defaultAddr, "Node IP address")
	dstAddrs = flag.String(flagDsts, defaultDstAddr, "Peer IP addresses")
	port = flag.String(flagPort, defaultPort, "Peer port number")
	numHeartbeats = flag.Int(flagNumHeartbeats, defaultNumHeartbeats, "Number of heartbeats to transmit")
	delayInt = flag.Int(flagDelayInt, defaultDelayInt, "Interval (ms) between heartbeats")
	start = flag.Int64(flagStart, defaultStart, "Unix epoch start time for heartbeat regime")
	trace = flag.Bool("trace", false, "Turn on tracing")
}

func validateAll() {
	validateNumHeartbeats()
	validateDelayInterval()
}

func validateNumHeartbeats() {
	if *numHeartbeats <= 0 {
		log.Error.Fatalf("Invalid environment variable value: %s",
			flagNumHeartbeats)
	}
}

func validateDelayInterval() {
	if *delayInt < 0 {
		log.Error.Fatalf("Invalid environment variable value: %s",
			flagDelayInt)
	}
}

func validateStart() {
	if *start < 0 {
		log.Error.Fatalf("Invalid environment variable value: %s",
			flagStart)
	}
}
