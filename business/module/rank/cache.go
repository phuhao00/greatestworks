package rank

import "sync"

type Cache struct {
	size          int64
	lastDataIndex int64
	rankID        uint32
	sortType      uint32
	rankData      []KV
	dataMutex     sync.Mutex
}

type KV struct {
	PlayerId uint64
	Score    int64
	SetTM    int64 //set time
}

func newCacheRank(rankID uint32, size int64, sortType uint32) *Cache {
	cr := &Cache{
		rankID:   rankID,
		size:     size,
		sortType: sortType,
		rankData: make([]KV, size, size),
	}
	return cr
}

func (c *Cache) getRank(begin, end int64, blackPlayerIds []uint64) ([]KV, bool) {

	if end > c.lastDataIndex {
		end = c.lastDataIndex
	}

	getLen := end - begin
	if begin < 0 || getLen <= 0 || end >= c.size {
		return nil, false // 参数错误
	}

	kvs := make([]KV, 0, getLen)

	c.dataMutex.Lock()
	for idx := begin; idx <= end; idx++ {
		kv := KV{
			PlayerId: c.rankData[idx].PlayerId,
			Score:    c.rankData[idx].Score,
			SetTM:    c.rankData[idx].SetTM,
		}
		for _, id := range blackPlayerIds {
			if kv.PlayerId == id {
				continue
			}
		}

		kvs = append(kvs, kv)
	}
	c.dataMutex.Unlock()

	return kvs, true
}

func (c *Cache) refresh(kvData []KV) {

	var (
		getCnt        int64 = 500
		curStartIndex int64 = 0
		curEndIndex   int64 = 499
	)
	for curEndIndex < c.size {
		dataLen := len(kvData)
		if dataLen <= 0 {
			curEndIndex = 0
			break
		}

		c.setCacheData(curStartIndex, curStartIndex+int64(dataLen), kvData)

		if dataLen < int(getCnt) {
			curEndIndex = curStartIndex + int64(dataLen)
			break
		}

		curStartIndex += getCnt
		curEndIndex += getCnt
	}
	c.lastDataIndex = curEndIndex
}

func (c *Cache) setCacheData(begin, end int64, data []KV) {
	getLen := end - begin
	if begin < 0 || getLen <= 0 || end > c.size {
		return
	}

	c.dataMutex.Lock()
	for idx := begin; idx < end; idx++ {
		c.rankData[idx].PlayerId = data[idx-begin].PlayerId
		c.rankData[idx].Score = data[idx-begin].Score
		c.rankData[idx].SetTM = data[idx-begin].SetTM
	}
	c.dataMutex.Unlock()
}
