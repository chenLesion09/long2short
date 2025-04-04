package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SequenceModel = (*customSequenceModel)(nil)

type (
	// SequenceModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSequenceModel.
	SequenceModel interface {
		sequenceModel
		withSession(session sqlx.Session) SequenceModel
	}

	customSequenceModel struct {
		*defaultSequenceModel
	}
)

// NewSequenceModel returns a model for the database table.
func NewSequenceModel(conn sqlx.SqlConn) SequenceModel {
	return &customSequenceModel{
		defaultSequenceModel: newSequenceModel(conn),
	}
}

func (m *customSequenceModel) withSession(session sqlx.Session) SequenceModel {
	return NewSequenceModel(sqlx.NewSqlConnFromSession(session))
}
