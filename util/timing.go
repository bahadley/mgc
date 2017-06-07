package util

import (
	"fmt"
	"time"

	"github.com/bahadley/mgc/config"
	"github.com/bahadley/mgc/log"
)

func DurationToRegimeStart() time.Duration {
	return (config.Start()).Sub(time.Now())
}

func DurationOfHeartbeatInterval() time.Duration {
	d, err := time.ParseDuration(fmt.Sprintf("%dms", config.DelayInterval()))
	if err != nil {
		log.Error.Fatal(err.Error())
	}
	return d
}
