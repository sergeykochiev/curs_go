package types

import (
	"database/sql"

	. "maragu.dev/gomponents"
)

type StateType struct {
	DB     *sql.DB
	MainDB *sql.DB
}

type DatabaseRecord struct {
	Id             int
	Name           string
	Filepath       string
	Is_initialized int
}

type Scanner interface {
	Scan(dest ...any) error
}

type HtmlEntity interface {
	HtmlTemplater
	ActiveRecorder
	Identifier
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
	GetDataRow() Group
	GetReadableName() string
	GetEntityPage(recursive bool) Group
}

type HtmlCreatable interface {
	GetCreateForm(arg ...HtmlEntity) Group
}

type Validator interface {
	Validate() bool
}

type Identifier interface {
	GetName() string
	GetId() int
}
