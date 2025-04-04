package logic

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"shortener/internal/svc"
	"shortener/internal/types"
	"shortener/model"
	"shortener/pkg/base62"
	"shortener/pkg/connect"
	"shortener/pkg/md5"
	"shortener/pkg/urltool"
)

type ConvertLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConvertLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConvertLogic {
	return &ConvertLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConvertLogic) Convert(req *types.ConvertRequest) (resp *types.ConvertResponse, err error) {
	// 1. 校验输入的数据
	// 1.1 输入的数据不得为空，这一步骤在handler一层解决
	// 1.2 输入的长链地址必须是可以ping通过的
	if ok := connect.Get(req.LongUrl); !ok {
		return nil, errors.New("long url is not reachable")
	}
	// 1.3 判断输入的链接是否已被转链，防止重复转链
	// 1.3.1 使用长连接获取到MD5
	md5Val := md5.Encrypt([]byte(req.LongUrl))
	// 1.3.2 利用MD5去查找数据库中是否有对应的数据
	result, err := l.svcCtx.ShortUrlModel.FindOneByMd5(l.ctx, sql.NullString{String: md5Val, Valid: true})
	if !errors.Is(err, sqlx.ErrNotFound) {
		if err == nil {
			return nil, fmt.Errorf("该链接已经被转还为%s", result.Surl.String)
		}
		logx.Errorw("ShortUrlModel FindOneByMd5() error", logx.LogField{Key: "error", Value: err.Error()})
		return nil, err
	}
	// 1.4 判断输入的链接是否已为锻炼，防止循环转链
	// 1.4.1 确保输入的是一个完整的URL，比如https://www.baidu.com/books/?id:123456789
	basePath, err := urltool.GetBasePath(req.LongUrl)
	if err != nil {
		logx.Errorw("urltool GetBasePath failed", logx.LogField{Key: "err", Value: err.Error()})
		return nil, err
	}

	// 1.4.2 利用短链接去查找数据库中是否有对应的数据
	_, err = l.svcCtx.ShortUrlModel.FindOneBySurl(l.ctx, sql.NullString{String: basePath, Valid: true})
	if !errors.Is(err, sqlx.ErrNotFound) {
		if err == nil {
			return nil, errors.New("该链接已经是短链，不必再转")
		}
		logx.Errorw("ShortUrlModel FindOneBySurl() error", logx.LogField{Key: "error", Value: err.Error()})
		return nil, err
	}

	var short string

	for {
		// 2. 取号，基于MySQL实现的发号器
		// 每接收到一个转链请求，就使用replace into语句往sequence里面更新数据，然后取出主键进行62位转化
		seq, err := l.svcCtx.Sequence.Next()
		if err != nil {
			logx.Errorw("Sequence Next() error", logx.LogField{Key: "error", Value: err.Error()})
			return nil, err
		}

		// 3. 将取出来的号码转化为62进制
		short = base62.Int2String(seq)
		// 3.1 通过for循环来检测短链是不是在黑名单中，避免特殊词汇比如：fuck, api, admin, etc
		if _, ok := l.svcCtx.ShortUrlBlackList[short]; !ok {
			break // 生成不在黑名单中的短链，跳出for循环
		}
	}

	fmt.Printf("short:%v\n", short)

	// 4. 存储长短链的映射关系
	_, err = l.svcCtx.ShortUrlModel.Insert(
		l.ctx,
		&model.ShortUrlMap{
			Lurl: sql.NullString{String: req.LongUrl, Valid: true},
			Md5:  sql.NullString{String: md5Val, Valid: true},
			Surl: sql.NullString{String: short, Valid: true},
		},
	)

	if err != nil {
		logx.Errorw("ShortUrlModel Insert() error", logx.LogField{Key: "error", Value: err.Error()})
		return nil, err
	}

	// 将短链返回给用户
	shortUrl := l.svcCtx.Config.ShortUrlDomain + "/" + short
	fmt.Printf("shortUrl:%v\n", shortUrl)
	return &types.ConvertResponse{ShortUrl: shortUrl}, nil
}
