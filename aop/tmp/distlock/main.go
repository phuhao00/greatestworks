package main

func main() {
	//// Create a pool with go-redis (or redigo) which is the pool redisync will
	//// use while communicating with Redis. This can also be any pool that
	//// implements the `redis.Pool` interface.
	//client := goredislib.NewClient(&goredislib.Options{
	//	Addr: "localhost:6379",
	//})
	//pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)
	//
	//// Create an instance of redisync to be used to obtain a mutual exclusion
	//// lock.
	//rs := redsync.New(pool)
	//
	//// Obtain a new mutex by using the same name for all instances wanting the
	//// same lock.
	//mutexname := "my-global-mutex"
	//mutex := rs.NewMutex(mutexname)
	//
	//// Obtain a lock for our given mutex. After this is successful, no one else
	//// can obtain the same lock (the same mutex name) until we unlock it.
	//if err := mutex.Lock(); err != nil {
	//	panic(err)
	//}
	//
	//// Do your work that requires the lock.
	//
	//// Release the lock so other processes or threads can obtain a lock.
	//if ok, err := mutex.Unlock(); !ok || err != nil {
	//	panic("unlock failed")
	//}
}
