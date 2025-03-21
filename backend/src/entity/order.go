package entity

import (
	"database/sql"
	"html"

	. "github.com/sergeykochiev/curs/backend/types"
	. "github.com/sergeykochiev/curs/backend/util"
	. "maragu.dev/gomponents"
	_ "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

type OrderEntity struct {
	Id           int
	Name         string
	Client_name  string
	Client_phone string
	Date_created string
	Date_ended   sql.NullString
	Ended        bool
	Creator_id   int
	_Creator     UserEntity
}

func (e *OrderEntity) ScanRow(r Scanner) error {
	return r.Scan(&e.Id, &e.Name, &e.Client_name, &e.Client_phone, &e.Date_created, &e.Creator_id, &e.Date_ended, &e.Ended, &e._Creator.Id, &e._Creator.Name, &e._Creator.Password, &e._Creator.Is_admin)
}

func (e *OrderEntity) GetSelectWhereQuery(where string) string {
	return "select * from public.order left join public.user on public.order.creator_id = public.user.id " + where
}

func (e *OrderEntity) Insert(db QueryExecutor) (sql.Result, error) {
	return db.Exec("insert into public.order (name, client_name, client_phone, date_created, creator_id) values ($1, $2, $3, $4, $5)", e.Name, e.Client_name, e.Client_phone, e.Date_created, e.Creator_id)
}

func (e OrderEntity) ToHtmlDataRow() Group {
	return Group{
		Div(Class("w-full grid place-items-center"), Text(html.EscapeString(e.Name))),
		Div(Class("w-full grid place-items-center"), Text(html.EscapeString(e.Client_name))),
		Div(Class("w-full grid place-items-center"), Text(html.EscapeString(e.Client_phone))),
		Div(Class("w-full grid place-items-center"), Text(html.EscapeString(e.Date_created))),
		Div(Class("w-full grid place-items-center"), Text(html.EscapeString(ConditionalArg(e.Ended, "Да", "Нет")))),
		Div(Class("w-full grid place-items-center"), Text(html.EscapeString(e.Date_ended.String))),
		Div(Class("w-full grid place-items-center"), Text(html.EscapeString(e._Creator.Name))),
	}
}

func (e OrderEntity) GetTableHeader() Group {
	return Group{
		Div(Class("w-full grid place-items-center"), Text("Название")),
		Div(Class("w-full grid place-items-center"), Text("Имя клиента")),
		Div(Class("w-full grid place-items-center"), Text("Телефон клиента")),
		Div(Class("w-full grid place-items-center"), Text("Дата создания")),
		Div(Class("w-full grid place-items-center"), Text("Завершен")),
		Div(Class("w-full grid place-items-center"), Text("Дата завершения")),
		Div(Class("w-full grid place-items-center"), Text("Создатель")),
	}
}

func (e OrderEntity) GetReadableName() string {
	return "Заказ"
}

func (e OrderEntity) GetId() int {
	return e.Id
}

func (e OrderEntity) GetName() string {
	return "order"
}

func (e *OrderEntity) Validate() bool {
	return len(e.Client_phone) == 11
}
