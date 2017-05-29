package heartbeat

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bahadley/mgc/log"
)

const (
	envAddr          = "MGC_ADDR"
	envDstAddr       = "MGC_DST_ADDR"
	envDstPort       = "MGC_DST_PORT"
	envNumHeartbeats = "MGC_NUM_HEARTBEATS"
	envDelayInt      = "MGC_DELAY_INTERVAL"
	envTransmit      = "MGC_TRANSMIT"
	envTrace         = "MGC_TRACE"

	defaultAddr          = "localhost"
	defaultDstAddr       = "localhost"
	defaultDstPort       = "22221"
	defaultNumHeartbeats = 10
	defaultDelayInt      = 1000
	traceFlag            = "YES"
	noTransmitFlag       = "NO"
)

func Addr() string {
	addr := os.Getenv(envAddr)
	if len(addr) == 0 {
		return defaultAddr
	} else {
		return addr
	}
}

func DstAddr() []string {
	addr := os.Getenv(envDstAddr)
	if len(addr) == 0 {
		return []string{defaultDstAddr}
	} else {
		return strings.Split(addr, ",")
	}
}

func DstPort() string {
	port := os.Getenv(envDstPort)
	if len(port) == 0 {
		return defaultDstPort
	} else {
		return port
	}
}

func NumHeartbeats() int {
	var numHeartbeats int

	env := os.Getenv(envNumHeartbeats)
	if len(env) == 0 {
		numHeartbeats = defaultNumHeartbeats
	} else {
		val, err := strconv.Atoi(env)
		if err != nil {
			log.Error.Fatalf("Invalid environment variable: %s",
				envNumHeartbeats)
		}

		if val <= 0 {
			log.Error.Fatalf("Invalid environment variable value: %s",
				envNumHeartbeats)
		}
		numHeartbeats = val
	}

	return numHeartbeats
}

func DelayInterval() time.Duration {
	var delayInterval time.Duration

	env := os.Getenv(envDelayInt)
	if len(env) == 0 {
		delayInterval = defaultDelayInt
	} else {
		val, err := strconv.Atoi(env)
		if err != nil {
			log.Error.Fatalf("Invalid environment variable: %s",
				envDelayInt)
		}

		if val < 0 {
			log.Error.Fatalf("Invalid environment variable value: %s",
				envDelayInt)
		}
		delayInterval = time.Duration(val)
	}

	return delayInterval
}

func Trace() bool {
	t := os.Getenv(envTrace)
	if len(t) > 0 && strings.ToUpper(t) == traceFlag {
		return true
	} else {
		return false
	}
}

func Send() bool {
	t := os.Getenv(envTransmit)
	if len(t) > 0 && strings.ToUpper(t) == noTransmitFlag {
		return false
	} else {
		return true
	}
}
