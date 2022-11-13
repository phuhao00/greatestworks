package chat

import (
	"container/ring"
	"go.mongodb.org/mongo-driver/bson"
)

type System struct {
	latestOnlineMessages    *ring.Ring
	latestCrossZoneMessages *ring.Ring
	latestZoneMessages      *ring.Ring
	latestCrossSrvMessages  *ring.Ring
	Owner
}

type Chat struct {
	Id      uint64
	Content string
}

func (c *Chat) ToDB() *Model {
	return &Model{
		Id:      c.Id,
		Content: c.Content,
	}
}

func (c *System) SetOwner(owner Owner) {
	c.Owner = owner
}

func (c *System) ToDB() bson.M {

	return nil
}
