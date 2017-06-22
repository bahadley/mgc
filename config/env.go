package config

const (
	defaultTupleBufLen = 128
	defaultTupleBufCap = 1024
	defaultChanBufSz   = 1000

	defaultWindowSz = 4
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
