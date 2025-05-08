package entity

import (
	. "github.com/sergeykochiev/curs/backend/gui"
	"github.com/shopspring/decimal"
	. "maragu.dev/gomponents"
	_ "maragu.dev/gomponents/components"
)

type UserEntity struct {
	Id       decimal.Decimal `gorm:"primaryKey"`
	Name     string
	Password string
	Is_admin bool
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

func (e UserEntity) GetEntityPage(recursive bool) Group {
	return Group{
		LabeledFieldComponent("Имя", e.Name),
	}
}

func (e UserEntity) Validate() bool {
	return len(e.Password) >= 8
}

func (e UserEntity) CheckPassword(password string) bool {
	return password == e.Password
}

func (e UserEntity) GetId() int64 {
	return e.Id.IntPart()
}

func (e UserEntity) TableName() string {
	return "user"
}
