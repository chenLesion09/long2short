package sequence

import (
	"database/sql"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const replaceSQL = "replace into sequence(stub) values ('a')"

type MySQL struct {
	conn sqlx.SqlConn
}

func NewMySQL(dsn string) Sequence {
	return &MySQL{
		conn: sqlx.NewMysql(dsn),
	}
}

func (m *MySQL) Next() (seq uint64, err error) {
	// SQL预处理
	var stmt sqlx.StmtSession

	stmt, err = m.conn.Prepare(replaceSQL)

	if err != nil {
		logx.Errorw("replaceSQL prepare error", logx.LogField{Key: "error", Value: err.Error()})
		return 0, err
	}

	defer func(stmt sqlx.StmtSession) {
		err := stmt.Close()
		if err != nil {
			logx.Errorw("stmt Close() error", logx.LogField{Key: "error", Value: err.Error()})
			return
		}
	}(stmt)

	// 执行SQL
	var rest sql.Result
	rest, err = stmt.Exec()
	if err != nil {
		logx.Errorw("replaceSQL exec error", logx.LogField{Key: "error", Value: err.Error()})
		return 0, err
	}

	// 获取最新的主键id
	var lid int64
	lid, err = rest.LastInsertId()
	if err != nil {
		logx.Errorw("get LastInsertId() error", logx.LogField{Key: "error", Value: err.Error()})
		return 0, err
	}

	return uint64(lid), nil
}
