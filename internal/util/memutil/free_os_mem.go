package memutil

var defaultBgFreer = newBgFreer()

func FreeOSMemoryPeriodically() {
	defaultBgFreer.freeOSMemoryPeriodically()
}

func StopFreeOSMemory() {
	defaultBgFreer.stop()
}
