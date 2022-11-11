package chat

import "container/ring"

type Chat struct {
	latestOnlineMessages    *ring.Ring
	latestCrossZoneMessages *ring.Ring
	latestZoneMessages      *ring.Ring
	latestCrossSrvMessages  *ring.Ring
	Owner
}
