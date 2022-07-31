package mongo

import (
	"context"
	"sync"

	"github.com/phuhao00/broker"
	mongobrocker "github.com/phuhao00/broker/mongo"
)

var (
	Client        *mongobrocker.Client
	onceInitMongo sync.Once
)

func init() {
	onceInitMongo.Do(func() {
		ctx := context.Background()
		tc := &mongobrocker.Client{
			BaseComponent: broker.NewBaseComponent(),
			RealCli: mongobrocker.NewClient(ctx, &mongobrocker.Config{
				URI:         "mongodb://localhost:27017",
				MinPoolSize: 3,
				MaxPoolSize: 3000,
			}),
		}
		tc.Launch()
	})
}
