package config

import (
	"time"
)

const (
	defaultTupleBufLen = 128
	defaultTupleBufCap = 1024
	defaultChanBufSz   = 1000

	// Number of observations to maintain in sliding window.
	defaultWindowSz = 4

	// Duration in milliseconds of constant safety margin.
	defaultSafetyMargin = 15

	// Deadline duration for initial bootstraping of failure
	// detection (i.e., no observations are available).
	defaultDeadline = 20
)

func TupleBufLen() uint32 {
	return defaultTupleBufLen
}

func TupleBufCap() uint32 {
	return defaultTupleBufCap
}

func ChannelBufSz() int {
	return defaultChanBufSz
}

func DefaultWindowSz() uint32 {
	return defaultWindowSz
}

func DefaultSafetyMargin() time.Duration {
	return time.Duration(defaultSafetyMargin)
}

func DefaultDeadline() time.Duration {
	return time.Duration(defaultDeadline)
}
