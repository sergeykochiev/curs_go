package entity

import (
	"database/sql"
	"fmt"
	"html"
	"net/url"
	"strconv"

	. "github.com/sergeykochiev/curs/backend/gui"
	. "github.com/sergeykochiev/curs/backend/util"
	"gorm.io/gorm"
	. "maragu.dev/gomponents"
	_ "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

type OrderEntity struct {
	ID           int
	Name         string
	Client_name  string
	Client_phone string
	Date_created string
	Date_ended   sql.NullString
	Ended        int
	Creator_id   int
	_Creator     UserEntity
}

func (e OrderEntity) GetDataRow() Group {
	return Group{
		Div(Class("px-[2px] grid place-items-center"), Text(html.EscapeString(fmt.Sprintf("%d", e.ID)))),
		TableCell(e.Name),
		TableCell(e.Client_name),
		TableCell(e.Client_phone),
		TableCell(e.Date_created),
		TableCell(ConditionalArg(e.Ended == 1, "Да", "Нет")),
		TableCell(e.Date_ended.String),
		TableCell(e._Creator.Name),
	}
}

func (e OrderEntity) GetTableHeader() Group {
	return Group{
		Div(Class("px-[2px] grid place-items-center"), Text("ID")),
		TableCell("Название"),
		TableCell("Имя клиента"),
		TableCell("Телефон клиента"),
		TableCell("Дата создания"),
		TableCell("Завершен"),
		TableCell("Дата завершения"),
		TableCell("Создатель"),
	}
}

func (e OrderEntity) GetEntityPage(recursive bool) Group {
	return Group{
		LabeledField("Название", e.Name),
		LabeledField("Имя клиента", e.Client_name),
		LabeledField("Телефон клиента", e.Client_phone),
		LabeledField("Дата создания", e.Date_created),
		LabeledField("Завершен", ConditionalArg(e.Ended == 1, "ДА", "НЕТ")),
		LabeledField("Дата завершения", ConditionalArg(e.Date_ended.Valid, e.Date_ended.String, "-")),
		If(recursive, Div(
			Class("bg-gray-100 flex flex-col gap-[8px] p-[8px]"),
			H2(Text(fmt.Sprintf("Создал пользователь #%d", e.Creator_id))),
			LabeledField("Имя", e._Creator.Name),
		)),
	}
}

func (e OrderEntity) GetCreateForm(db *gorm.DB) Group {
	return Group{
		InputComponent("text", "", "name", "Название заказа", "", true),
		InputComponent("text", "", "client_name", "Имя клиента", "", true),
		InputComponent("number", "", "client_phone", "Телефон клиента", "", true),
		InputComponent("date", "", "date_created", "Дата создания", "", true),
	}
}

func (e OrderEntity) GetReadableName() string {
	return "Заказ"
}

func (e OrderEntity) GetId() int {
	return e.ID
}

func (e OrderEntity) GetName() string {
	return "order"
}

func (e *OrderEntity) Validate() bool {
	return len(e.Client_phone) == 11
}

func (e *OrderEntity) ValidateAndParseForm(form url.Values) bool {
	if !form.Has("name") || !form.Has("client_name") || !form.Has("client_phone") || !form.Has("date_created") {
		return false
	}
	e.Client_name = form.Get("client_name")
	e.Client_phone = form.Get("client_phone")
	e.Name = form.Get("name")
	e.Date_created = form.Get("date_created")
	var err error
	e.Creator_id, err = strconv.Atoi(form.Get("userid"))
	return err != nil
}
