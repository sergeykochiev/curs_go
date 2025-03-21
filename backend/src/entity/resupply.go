package entity

import (
	"database/sql"
	"fmt"
	"html"

	. "github.com/sergeykochiev/curs/backend/types"
	. "maragu.dev/gomponents"
	_ "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

type ResourceResupplyEntity struct {
	Id             int
	Resource_id    int
	Quantity_added int
	Date           string
	_Resource      ResourceEntity
}

func (e *ResourceResupplyEntity) ScanRow(r Scanner) error {
	return r.Scan(&e.Id, &e.Resource_id, &e.Quantity_added, &e.Date, &e._Resource.Id, &e._Resource.Name, &e._Resource.Date_last_updated, &e._Resource.Cost_by_one, &e._Resource.Quantity)
}

func (e *ResourceResupplyEntity) GetSelectWhereQuery(where string) string {
	return "select * from public.resource_resupply  left join public.resource on public.resource_resupply.resource_id = public.resource.id" + where
}

func (e *ResourceResupplyEntity) Insert(db QueryExecutor) (sql.Result, error) {
	if err := e._Resource.ScanRow(db.QueryRow(e._Resource.GetSelectWhereQuery("where public.resource.id = $1"), e._Resource.Id)); err != nil {
		return nil, err
	}
	e._Resource.Quantity += e.Quantity_added
	if _, err := e._Resource.Update(db); err != nil {
		return nil, err
	}
	return db.Exec("insert into public.resource_resupply (resource_id, quantity_added, date) values ($1, $2, $3)", e.Resource_id, e.Quantity_added, e.Date)
}

// TODO implement me
func (e *ResourceResupplyEntity) Update(db QueryExecutor) (sql.Result, error) {
	return db.Exec("")
}

func (e *ResourceResupplyEntity) ToHtmlDataRow() Group {
	return Group{
		Div(Class("w-full grid place-items-center"), Text(html.EscapeString(e._Resource.Name))),
		Div(Class("w-full grid place-items-center"), Text(html.EscapeString(fmt.Sprintf("%d", e.Quantity_added)))),
		Div(Class("w-full grid place-items-center"), Text(html.EscapeString(e.Date))),
		Div(Class("w-full grid place-items-center"), Text(html.EscapeString(fmt.Sprintf("%f", e._Resource.Cost_by_one)))),
	}
}

func (e *ResourceResupplyEntity) GetTableHeader() Group {
	return Group{
		Div(Class("w-full grid place-items-center"), Text("Ресурс")),
		Div(Class("w-full grid place-items-center"), Text("Количество добавлено (единиц)")),
		Div(Class("w-full grid place-items-center"), Text("Дата поставки")),
		Div(Class("w-full grid place-items-center"), Text("Цена за один")),
	}
}

func (e *ResourceResupplyEntity) GetReadableName() string {
	return "Поставка ресурса"
}

func (e *ResourceResupplyEntity) Validate() bool {
	return true
}

func (e *ResourceResupplyEntity) GetId() int {
	return e.Id
}

func (e *ResourceResupplyEntity) GetName() string {
	return "resource_resupply"
}
