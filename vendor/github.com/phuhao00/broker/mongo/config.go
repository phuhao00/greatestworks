package mongobrocker

type Config struct {
	URI                      string
	MinPoolSize, MaxPoolSize uint64
}
