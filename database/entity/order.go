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
		TableCellComponent(e.Name),
		TableCellComponent(e.Client_name),
		TableCellComponent(e.Client_phone),
		TableCellComponent(e.Date_created),
		TableCellComponent(ConditionalArg(e.Ended == 1, "Да", "Нет")),
		TableCellComponent(e.Date_ended.String),
		TableCellComponent(e._Creator.Name),
	}
}

func (e OrderEntity) GetTableHeader() Group {
	return Group{
		Div(Class("px-[2px] grid place-items-center"), Text("ID")),
		TableCellComponent("Название"),
		TableCellComponent("Имя клиента"),
		TableCellComponent("Телефон клиента"),
		TableCellComponent("Дата создания"),
		TableCellComponent("Завершен"),
		TableCellComponent("Дата завершения"),
		TableCellComponent("Создатель"),
	}
}

func (e OrderEntity) GetEntityPage(recursive bool) Group {
	return Group{
		LabeledFieldComponent("Название", e.Name),
		LabeledFieldComponent("Имя клиента", e.Client_name),
		LabeledFieldComponent("Телефон клиента", e.Client_phone),
		LabeledFieldComponent("Дата создания", e.Date_created),
		LabeledFieldComponent("Завершен", ConditionalArg(e.Ended == 1, "ДА", "НЕТ")),
		LabeledFieldComponent("Дата завершения", ConditionalArg(e.Date_ended.Valid, e.Date_ended.String, "-")),
		If(recursive, Div(
			Class("bg-gray-100 flex flex-col gap-[8px] p-[8px]"),
			H2(Text(fmt.Sprintf("Создал пользователь #%d", e.Creator_id))),
			LabeledFieldComponent("Имя", e._Creator.Name),
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
