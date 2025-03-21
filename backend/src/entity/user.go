package entity

import (
	"database/sql"

	. "github.com/sergeykochiev/curs/backend/types"
)

type UserEntity struct {
	Id       int
	Name     string
	Password string
	Is_admin bool
}

func (e *UserEntity) ScanRow(r Scanner) error {
	return r.Scan(&e.Id, &e.Name, &e.Password, &e.Is_admin)
}

func (e *UserEntity) InsertRow(db *sql.DB) (sql.Result, error) {
	return db.Exec("insert into public.user (name, password, is_admin) values ($1, $2, $3)", e.Name, e.Password, e.Is_admin)
}

// func (e *UserEntity) getHtmlCreateForm() Group {
// 	return Group{
// 		InputComponent("text", "Имя клиента"),
// 		InputComponent("number", "Телефон клиента"),
// 		InputComponent("date", "Дата заказа"),
// 	}
// }

// func (e *UserEntity) toHtmlDataRow() Group {
// 	return Group{
// 		Div(Class("w-full grid place-items-center"), Text(html.EscapeString(fmt.Sprintf("%d", e.Id)))),
// 		Div(Class("w-full grid place-items-center"), Text(html.EscapeString(e.Client_name))),
// 		Div(Class("w-full grid place-items-center"), Text(html.EscapeString(e.Client_phone))),
// 		Div(Class("w-full grid place-items-center"), Text(html.EscapeString(e.Date))),
// 		Div(Class("w-full grid place-items-center"), Text(html.EscapeString(fmt.Sprintf("%d", e.Creator_id)))),
// 	}
// }

func (e *UserEntity) Validate() bool {
	return len(e.Password) >= 8
}

func (e *UserEntity) CheckPassword(password string) bool {
	return password == e.Password
}

func (e *UserEntity) GetId() int {
	return e.Id
}
