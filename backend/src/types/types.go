package types

import (
	"database/sql"

	. "maragu.dev/gomponents"
)

type TStateType struct {
	DB *sql.DB
}

type Scanner interface {
	Scan(dest ...any) error
}

type QueryExecutor interface {
	Exec(query string, args ...any) (sql.Result, error)
	QueryRow(query string, args ...any) *sql.Row
	Query(query string, args ...any) (*sql.Rows, error)
}

type ActiveRecorder interface {
	ScanRow(r Scanner) error
	GetSelectWhereQuery(where string) string
	Insert(db QueryExecutor) (sql.Result, error)
}

type HtmlTemplater interface {
	GetTableHeader() Group
	ToHtmlDataRow() Group
	GetReadableName() string
}

type Validator interface {
	Validate() bool
}

type Identifier interface {
	GetName() string
	GetId() int
}
