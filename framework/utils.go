package framework

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"time"
)

// MustEnv returns the value of the environment variable or panics if the variable is not set or if the type is unsupported.
func MustEnv[T any](key string, fallback T) T {
	value, ok := os.LookupEnv(key)
	if !ok {
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

// GenerateRandomString generates a random string of a given length using the characters provided.
func GenerateRandomString(length int) string {
	characters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = characters[rand.Intn(len(characters))]
	}

	return string(result)
}

// PrettyPrint converts a map to a pretty-printed JSON string
func PrettyPrint(data map[string]interface{}) (string, error) {
	// Marshal the map to JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}

	// Convert the JSON byte slice to a string
	return string(jsonData), nil
}
