package util

import (
	"database/sql"
	"fmt"
	"time"

	. "github.com/sergeykochiev/curs/backend/types"
)

func ConditionalArg[T any](condition bool, arg T, notarg T) T {
	if condition {
		return arg
	}
	return notarg
}

func GetCurrentTime() string {
	return time.Now().Format(time.DateTime)
}

func GetRows[T ActiveRecorder](db *sql.DB, ent T, where string) ([]T, error) {
	var rows *sql.Rows
	var err error
	rows, err = db.Query(ent.GetSelectWhereQuery(where))
	if err != nil {
		println("Failed to query rows", err.Error())
		return nil, err
	}
	var arr []T
	for rows.Next() {
		if err := ent.ScanRow(rows); err != nil {
			return arr, err
		}
		arr = append(arr, ent)
	}
	return arr, nil
}

func GetSingleRow[T interface {
	ActiveRecorder
	Identifier
}, D QueryExecutor](db D, ent T, id string) error {
	rows, err := db.Query(fmt.Sprintf(ent.GetSelectWhereQuery("where public.%s.id = $1"), ent.GetName()), id)
	if err != nil {
		println("Failed to query rows")
		return err
	}
	if !rows.Next() {
		return sql.ErrNoRows
	}
	return ent.ScanRow(rows)
}
