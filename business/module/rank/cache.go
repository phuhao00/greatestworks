package rank

import "sync"

type Cache struct {
	size          int64
	lastDataIndex int64
	rankID        uint32
	sortType      uint32
	rankData      []KV
	dataMutex     sync.Mutex
	isInitial     bool
}

type KV struct {
	PlayerId uint64
	Score    int64
	SetTM    int64 //set time
}
