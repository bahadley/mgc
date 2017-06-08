package config

import (
	"flag"
	"strings"
	"time"

	"github.com/bahadley/mgc/log"
)

const (
	flagRole     = "role"
	flagAddr     = "addr"
	flagDsts     = "dsts"
	flagPort     = "port"
	flagDelayInt = "hbdelay"
	flagStart    = "start"
	flagTrace    = "trace"

	leaderFlag   = "l"
	followerFlag = "f"

	defaultAddr     = "localhost"
	defaultDstAddr  = "localhost"
	defaultPort     = "22221"
	defaultDelayInt = 1000
	defaultStart    = 0
	defaultTrace    = true
)

var (
	role     *string
	addr     *string
	dstAddrs *string
	port     *string
	delayInt *int
	start    *int64

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

func Start() time.Time {
	return time.Unix(*start, 0)
}

func DurationToRegimeStart() time.Duration {
	return (time.Unix(*start, 0)).Sub(time.Now())
}

func DurationOfHeartbeatInterval() time.Duration {
	return time.Duration(*delayInt) * time.Millisecond
}

func init() {
	setFlags()
	flag.Parse()
	validateAll()
	log.SetTrace(*trace)
}

func setFlags() {
	role = flag.String(flagRole, leaderFlag, "Node role [(l)eader,(f)ollower]")
	addr = flag.String(flagAddr, defaultAddr, "Node IP address")
	dstAddrs = flag.String(flagDsts, defaultDstAddr, "Peer IP addresses")
	port = flag.String(flagPort, defaultPort, "Peer port number")
	delayInt = flag.Int(flagDelayInt, defaultDelayInt, "Interval (ms) between heartbeats")
	start = flag.Int64(flagStart, defaultStart, "Unix epoch start time for heartbeat regime")
	trace = flag.Bool("trace", false, "Turn on tracing")
}

func validateAll() {
	validateRole()
	validateDelayInterval()
	validateStart()
}

func validateRole() {
	if *role != leaderFlag && *role != followerFlag {
		log.Error.Fatalf("Invalid environment variable value: %s",
			flagRole)
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
