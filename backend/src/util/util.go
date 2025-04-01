package util

import (
	"database/sql"
	"errors"
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

func GetRows[T ActiveRecorder](db QueryExecutor, ent T, where string) ([]T, error) {
	var rows *sql.Rows
	var err error
	rows, err = db.Query(ent.GetSelectWhereQuery(where))
	if err != nil {
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
}](db QueryExecutor, ent T, id int) error {
	println(ent.GetSelectWhereQuery(fmt.Sprintf("where \"%s\".id = $1", ent.GetName())))
	rows, err := db.Query(ent.GetSelectWhereQuery(fmt.Sprintf("where \"%s\".id = $1", ent.GetName())), id)
	if err != nil {
		return err
	}
	if !rows.Next() {
		return sql.ErrNoRows
	}
	return ent.ScanRow(rows)
}

func PreappendError(prefix string, err error) error {
	return errors.New(fmt.Sprintf("%s: %s", prefix, err.Error()))
}
