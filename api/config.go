package api

type ConfigMap map[string]interface{}

var Conf ConfigMap = ConfigMap{}

func (c ConfigMap) set(key string, value interface{}) {
	c[key] = value
}

func (c ConfigMap) get(key string) interface{} {
	return c[key]
}

func SetConfig(key string, value interface{}) {
	Conf.set(key, value)
}

func Config(key string) interface{} {
	return Conf.get(key)
}
