package svc

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"shortener/internal/config"
	"shortener/model"
	"shortener/sequence"
)

type ServiceContext struct {
	Config config.Config

	ShortUrlModel model.ShortUrlMapModel // short_url_map

	Sequence sequence.Sequence // sequence

	ShortUrlBlackList map[string]struct{}
}

func NewServiceContext(c config.Config) *ServiceContext {
	dsn1 := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		c.ShortUrlDB.User,
		c.ShortUrlDB.Password,
		c.ShortUrlDB.Host,
		c.ShortUrlDB.Port,
		c.ShortUrlDB.DBName,
	)
	conn1 := sqlx.NewMysql(dsn1)

	dsn2 := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		c.Sequence.User,
		c.Sequence.Password,
		c.Sequence.Host,
		c.Sequence.Port,
		c.Sequence.DBName,
	)

	shortUrlBlackList := make(map[string]struct{}, len(c.ShortUrlBlackList))

	for _, v := range c.ShortUrlBlackList {
		shortUrlBlackList[v] = struct{}{}
	}

	return &ServiceContext{
		Config:        c,
		ShortUrlModel: model.NewShortUrlMapModel(conn1),
		Sequence:      sequence.NewMySQL(dsn2),
	}
}
