package entity

import (
	"database/sql"
	"fmt"
	"html"

	. "github.com/sergeykochiev/curs/backend/types"
	. "github.com/sergeykochiev/curs/backend/util"
	. "maragu.dev/gomponents"
	_ "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

type ResourceEntity struct {
	Id                int
	Name              string
	Date_last_updated string
	Cost_by_one       float32
	Quantity          int
}

func (e *ResourceEntity) ScanRow(r Scanner) error {
	return r.Scan(&e.Id, &e.Name, &e.Date_last_updated, &e.Cost_by_one, &e.Quantity)
}

func (e *ResourceEntity) GetSelectWhereQuery(where string) string {
	return "select * from \"resource\"" + where
}

func (e *ResourceEntity) Insert(db QueryExecutor) (sql.Result, error) {
	return db.Exec("insert into public.resource (name, date_last_updated, cost_by_one, quantity) values ($1, $2, $3, $4)", e.Name, GetCurrentTime(), e.Cost_by_one, e.Quantity)
}

func (e *ResourceEntity) Update(db QueryExecutor) (sql.Result, error) {
	return db.Exec("update public.resource set name = $1, date_last_updated = $2, cost_by_one = $3, quantity = $4 where public.resource_entity.id = $5", e.Name, GetCurrentTime(), e.Cost_by_one, e.Quantity, e.Id)
}

func (e *ResourceEntity) ToHtmlDataRow() Group {
	return Group{
		Div(Class("w-full grid place-items-center"), Text(html.EscapeString(e.Name))),
		Div(Class("w-full grid place-items-center"), Text(html.EscapeString(e.Date_last_updated))),
		Div(Class("w-full grid place-items-center"), Text(html.EscapeString(fmt.Sprintf("%f", e.Cost_by_one)))),
		Div(Class("w-full grid place-items-center"), Text(html.EscapeString(fmt.Sprintf("%d", e.Quantity)))),
	}
}

func (e *ResourceEntity) GetTableHeader() Group {
	return Group{
		Div(Class("w-full grid place-items-center"), Text("Наименование")),
		Div(Class("w-full grid place-items-center"), Text("Дата обновления")),
		Div(Class("w-full grid place-items-center"), Text("Цена за единицу")),
		Div(Class("w-full grid place-items-center"), Text("Количество")),
	}
}

func (e *ResourceEntity) GetReadableName() string {
	return "Ресурс"
}

func (e *ResourceEntity) GetId() int {
	return e.Id
}

func (e *ResourceEntity) Validate() bool {
	return true
}

func (e *ResourceEntity) GetName() string {
	return "resource"
}
