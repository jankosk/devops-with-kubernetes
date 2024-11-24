package common

import "os"

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func CheckErr(e error) {
	if e != nil {
		panic(e)
	}
}
