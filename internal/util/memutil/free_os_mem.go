package memutil

var defaultBgFreer = newBgFreer()

func FreeOSMemoryPeriodically() {
	go defaultBgFreer.freeOSMemoryPeriodically()
}

func StopFreeOSMemory() {
	defaultBgFreer.stop()
}
