package config

import (
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

type ConfigMap map[string]interface{}

var (
	Conf ConfigMap
	mu   sync.RWMutex
)

func init() {
	Conf = make(ConfigMap)
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func Set(key string, value interface{}) {
	mu.Lock()
	defer mu.Unlock()

	keys := strings.Split(key, ".")
	setRecursive(Conf, keys, value)
}

func setRecursive(current ConfigMap, keys []string, value interface{}) {
	if len(keys) == 1 {
		if nestedMap, ok := value.(map[string]interface{}); ok {
			// If the value is a map, convert it to ConfigMap
			configMap := make(ConfigMap)
			for k, v := range nestedMap {
				setRecursive(configMap, []string{k}, v)
			}
			current[keys[0]] = configMap
		} else {
			current[keys[0]] = value
		}
	} else {
		if _, ok := current[keys[0]]; !ok {
			current[keys[0]] = make(ConfigMap)
		}
		setRecursive(current[keys[0]].(ConfigMap), keys[1:], value)
	}
}

func Get[T any](key string, fallback ...T) T {
	mu.RLock()
	defer mu.RUnlock()

	keys := strings.Split(key, ".")
	current := Conf

	for _, k := range keys {
		if val, exists := current[k]; exists {
			// Try to type assert the value to the expected type
			if typedVal, ok := val.(T); ok {
				return typedVal
			}
			// Handle nested maps as ConfigMap
			if m, ok := val.(ConfigMap); ok {
				current = m
			} else {
				break // Exit the loop if we can't type assert
			}
		} else {
			// Return fallback if provided
			if len(fallback) > 0 {
				return fallback[0]
			}
			// Zero value of T if no fallback
			var zero T
			return zero
		}
	}

	// If the key points to a map, return the map if it matches T
	if len(keys) == 1 {
		if typedVal, ok := current[keys[0]].(T); ok {
			return typedVal
		}
	}

	// Return fallback if nothing was found
	if len(fallback) > 0 {
		return fallback[0]
	}

	var zero T
	return zero
}

func GetAll() ConfigMap {
	mu.RLock()
	defer mu.RUnlock()

	return Conf
}

func Reset() {
	mu.Lock()
	defer mu.Unlock()

	Conf = make(ConfigMap)
}

// MustEnv returns the value of the environment variable or panics if the variable is not set or if the type is unsupported.
func MustEnv[T any](key string, fallback T) T {
	value, ok := os.LookupEnv(key)
	if !ok {
		slog.Info(fmt.Sprintf("Using fallback value for key: %s", key), key, fallback)
		return fallback
	}

	fallbackType := reflect.TypeOf(fallback)

	// Check if the fallback type is supported
	switch fallbackType.Kind() {
	case reflect.Int, reflect.Int64:
		result, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			panic(err)
		}
		return reflect.ValueOf(result).Convert(fallbackType).Interface().(T)

	case reflect.Float64:
		result, err := strconv.ParseFloat(value, 64)
		if err != nil {
			panic(err)
		}
		return reflect.ValueOf(result).Convert(fallbackType).Interface().(T)

	case reflect.Bool:
		result, err := strconv.ParseBool(value)
		if err != nil {
			panic(err)
		}
		return reflect.ValueOf(result).Convert(fallbackType).Interface().(T)

	case reflect.String:
		return reflect.ValueOf(value).Convert(fallbackType).Interface().(T)

	default:
		panic(fmt.Sprintf("unsupported type: %v", fallbackType))
	}
}
