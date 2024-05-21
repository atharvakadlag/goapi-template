package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	PublicHost             string
	Port                   string
	DBUser                 string
	DBPasswd               string
	DBAddr                 string
	DBName                 string
	JWTExpirationInSeconds int64
	JWTSecret              string
}

var Envs = initConfig()

func initConfig() Config {
	return Config{
		PublicHost: getEnvStr("HOST", "127.0.0.1"),
		Port:       getEnvStr("PORT", "8001"),

		DBUser:                 getEnvStr("DB_USER", "root"),
		DBPasswd:               getEnvStr("DB_PASSWD", "toor"),
		DBAddr:                 getEnvStr("DB_ADDR", fmt.Sprintf("%s:%s", getEnvStr("DB_HOST", "127.0.0.1"), getEnvStr("DB_PORT", "3306"))),
		DBName:                 getEnvStr("DB_NAME", "splitfree"),
		JWTExpirationInSeconds: getEnvInt("JWT_EXP", 60*60*24*7),
		JWTSecret:              getEnvStr("JWT_SECRET", "super-long-secretive-string-that-no-one-can-find"),
	}
}

func getEnvStr(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}

	return fallback
}
