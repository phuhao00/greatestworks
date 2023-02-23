package pubsub

type ExchangeKind string

const (
	ExchangeKindDirect  ExchangeKind = "direct"
	ExchangeKindFanout  ExchangeKind = "fanout"
	ExchangeKindTopic   ExchangeKind = "topic"
	ExchangeKindHeaders ExchangeKind = "headers"
)
