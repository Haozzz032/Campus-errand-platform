package config

import "time"

const defaultDSN = "root:200732@tcp(127.0.0.1:3306)/campus?charset=utf8mb4&parseTime=True&loc=Local"

type Config struct {
	DBDSN     string
	Host      string
	Port      string
	JWTSecret string
	JWTTTL    time.Duration
}

func (cfg *Config) Addr() string {
	return cfg.Host + ":" + cfg.Port
}

func Loadconfig() Config {
	return Config{
		DBDSN:     defaultDSN,
		Host:      "0.0.0.0",
		Port:      "8080",
		JWTSecret: "TEST_TEST_TEST_TEST_TEST_TEST_TEST_TEST",
		JWTTTL:    7 * 24 * time.Hour,
	}
}
