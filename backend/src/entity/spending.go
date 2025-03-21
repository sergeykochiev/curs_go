package entity

import (
	"database/sql"
	"errors"
	"fmt"
	"html"

	. "github.com/sergeykochiev/curs/backend/types"
	. "maragu.dev/gomponents"
	_ "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

type ResourceSpendingEntity struct {
	Id             int
	Order_id       int
	Resource_id    int
	Quantity_spent int
	Date           string
	_Order         OrderEntity
	_Resource      ResourceEntity
}

func (e *ResourceSpendingEntity) ScanRow(r Scanner) error {
	return r.Scan(&e.Id, &e.Order_id, &e.Resource_id, &e.Quantity_spent, &e.Date, &e.Id, &e._Order.Id, &e._Order.Name, &e._Order.Client_name, &e._Order.Client_phone, &e._Order.Date_created, &e._Order.Date_ended, &e._Order.Ended, &e._Order.Creator_id, &e._Resource.Id, &e._Resource.Name, &e._Resource.Date_last_updated)
}

func (e *ResourceSpendingEntity) GetSelectWhereQuery(where string) string {
	return "select * from public.resource_spending left join public.order on public.order.id = public.resource_spending.order_id left join public.resource on public.resource.id = public.resource_spending.resource_id" + where
}

func (e *ResourceSpendingEntity) Insert(db QueryExecutor) (sql.Result, error) {
	if err := e._Resource.ScanRow(db.QueryRow(e._Resource.GetSelectWhereQuery("where public.resource.id = $1"), e.Resource_id)); err != nil {
		return nil, err
	}
	if e.Quantity_spent > e._Resource.Quantity {
		return nil, errors.New("Invalid quantity")
	}
	e._Resource.Quantity -= e.Quantity_spent
	if _, err := e._Resource.Update(db); err != nil {
		return nil, err
	}
	return db.Exec("insert into public.resource_spending (order_id, resource_id, quantity_spent, date) values ($1, $2, $3, $4)", e.Order_id, e.Resource_id, e.Quantity_spent, e.Date)
}

// TODO implement me
func (e *ResourceSpendingEntity) Update(db QueryExecutor) (sql.Result, error) {
	return db.Exec("")
}

func (e *ResourceSpendingEntity) ToHtmlDataRow() Group {
	return Group{
		Div(Class("w-full grid place-items-center"), Text(html.EscapeString(e._Order.Name))),
		Div(Class("w-full grid place-items-center"), Text(html.EscapeString(fmt.Sprintf("%d", e.Quantity_spent)))),
		Div(Class("w-full grid place-items-center"), Text(html.EscapeString(fmt.Sprintf("%d", e.Quantity_spent)))),
		Div(Class("w-full grid place-items-center"), Text(html.EscapeString(e.Date))),
	}
}

func (e *ResourceSpendingEntity) GetTableHeader() Group {
	return Group{
		Div(Class("w-full grid place-items-center"), Text("Название заказа")),
		Div(Class("w-full grid place-items-center"), Text("Наименование ресурса")),
		Div(Class("w-full grid place-items-center"), Text("Количество потрачено (единиц)")),
		Div(Class("w-full grid place-items-center"), Text("Дата траты")),
	}
}

func (e *ResourceSpendingEntity) GetReadableName() string {
	return "Трата ресурса"
}

func (e *ResourceSpendingEntity) Validate() bool {
	return true
}

func (e *ResourceSpendingEntity) GetId() int {
	return e.Id
}

func (e *ResourceSpendingEntity) GetName() string {
	return "resource_spending"
}
