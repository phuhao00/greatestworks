package chat

type Model struct {
	Id      uint64 `bson:"id"`
	Content string `bson:"content"`
}
