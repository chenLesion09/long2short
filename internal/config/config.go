package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf

	ShortUrlDB DBString

	Sequence DBString

	BaseString string

	ShortUrlBlackList []string

	ShortUrlDomain string
}

type DBString struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}
