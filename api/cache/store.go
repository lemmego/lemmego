package cache

type Store interface {
	Get(key string) interface{}
	Many(keys []string) map[string]interface{}
	Put(key string, value interface{}, seconds int)
	PutMany(values map[string]interface{}, seconds int)
	Increment(key string, value int) int
	Decrement(key string, value int) int
	Forever(key string, value interface{})
	Forget(key string) bool
	Flush() bool
	GetPrefix() string
}
