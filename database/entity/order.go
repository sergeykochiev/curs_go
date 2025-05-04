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
	Company_name sql.NullString
	Date_created string
	Date_ended   sql.NullString
	Ended        int
	Creator_id   int
	UserEntity   UserEntity `gorm:"foreignKey:Creator_id"`
}

func (e *OrderEntity) GetFilters() Group {
	return Group{
		DateFilterComponent("Дата создания в диапазоне", "date_created"),
		BoolFilterComponent("Закончен?", "ended"),
		DateFilterComponent("Дата завершения в диапазоне", "date_ended"),
		StringFilterComponent("Название включает", "name"),
		StringFilterComponent("Название компании включает", "company_name"),
		StringFilterComponent("Имя клиента включает", "client_name"),
		StringFilterComponent("Телефон клиента", "client_phone"),
	}
}

func (e *OrderEntity) GetFilteredDb(filters url.Values, db *gorm.DB) *gorm.DB {
	if filters.Has("date_created_lo") && filters.Get("date_created_lo") != "" {
		db = db.Where("date_created > ?", filters.Get("date_created_lo"))
	}
	if filters.Has("date_created_hi") && filters.Get("date_created_hi") != "" {
		db = db.Where("date_created < ?", filters.Get("date_created_hi"))
	}
	if filters.Has("date_ended_lo") && filters.Get("date_ended_lo") != "" {
		db = db.Where("date_ended > ?", filters.Get("date_ended_lo"))
	}
	if filters.Has("date_ended_hi") && filters.Get("date_ended_hi") != "" {
		db = db.Where("date_ended < ?", filters.Get("date_ended_hi"))
	}
	if filters.Has("client_name") && filters.Get("client_name") != "" {
		db = db.Where("client_name LIKE ?", "%"+filters.Get("client_name")+"%")
	}
	if filters.Has("name") && filters.Get("name") != "" {
		db = db.Where("name LIKE ?", "%"+filters.Get("name")+"%")
	}
	if filters.Has("company_name") && filters.Get("company_name") != "" {
		db = db.Where("company_name LIKE ?", "%"+filters.Get("company_name")+"%")
	}
	if filters.Has("client_phone") && filters.Get("client_phone") != "" {
		db = db.Where("client_phone = ?", filters.Get("client_phone"))
	}
	if filters.Has("ended") && filters.Get("ended") != "" {
		db = db.Where("ended = ?", filters.Get("ended"))
	}
	return db.Joins("UserEntity")
}

func (e OrderEntity) GetDataRow() Group {
	return Group{
		Div(Class("px-[2px] grid place-items-center"), Text(html.EscapeString(fmt.Sprintf("%d", e.ID)))),
		TableCellComponent(e.Name),
		TableCellComponent(e.Client_name),
		TableCellComponent(e.Client_phone),
		TableCellComponent(e.Company_name.String),
		TableCellComponent(e.Date_created),
		TableCellComponent(ConditionalArg(e.Ended == 1, "Да", "Нет")),
		TableCellComponent(e.Date_ended.String),
		TableCellComponent(e.UserEntity.Name),
	}
}

func (e OrderEntity) GetTableHeader() Group {
	return Group{
		Div(Class("px-[2px] grid place-items-center"), Text("ID")),
		TableCellComponent("Название"),
		TableCellComponent("Имя клиента"),
		TableCellComponent("Телефон клиента"),
		TableCellComponent("Компания клиента"),
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
		LabeledFieldComponent("Компания клиента", ConditionalArg(e.Company_name.Valid, e.Company_name.String, "-")),
		LabeledFieldComponent("Дата создания", e.Date_created),
		LabeledFieldComponent("Завершен", ConditionalArg(e.Ended == 1, "ДА", "НЕТ")),
		LabeledFieldComponent("Дата завершения", ConditionalArg(e.Date_ended.Valid, e.Date_ended.String, "-")),
		If(recursive, Div(
			Class("bg-gray-100 flex flex-col gap-[8px] p-[8px]"),
			H2(Text(fmt.Sprintf("Создал пользователь #%d", e.Creator_id))),
			LabeledFieldComponent("Имя", e.UserEntity.Name),
		)),
	}
}

func (e OrderEntity) GetCreateForm(db *gorm.DB) Group {
	return Group{
		LabeledInputComponent("text", "", "name", "Название заказа", "", true),
		LabeledInputComponent("text", "", "client_name", "Имя клиента", "", true),
		LabeledInputComponent("number", "", "client_phone", "Телефон клиента", "", true),
		LabeledInputComponent("number", "", "company_name", "Компания клиента", "", false),
		LabeledInputComponent("date", "", "date_created", "Дата создания", "", true),
	}
}

func (e OrderEntity) GetReadableName() string {
	return "Заказ"
}

func (e OrderEntity) GetId() int {
	return e.ID
}

func (e OrderEntity) TableName() string {
	return "order"
}

func (e *OrderEntity) Validate() bool {
	return len(e.Client_phone) == 11
}

func (e *OrderEntity) ValidateAndParseForm(form url.Values) bool {
	if !form.Has("name") || !form.Has("client_name") || !form.Has("client_phone") || !form.Has("date_created") {
		return false
	}
	if form.Has("company_name") {
		e.Company_name.String = form.Get("company_name")
		e.Company_name.Valid = true
	}
	e.Client_name = form.Get("client_name")
	e.Client_phone = form.Get("client_phone")
	e.Name = form.Get("name")
	e.Date_created = form.Get("date_created")
	var err error
	e.Creator_id, err = strconv.Atoi(form.Get("userid"))
	return err != nil
}
