package entity

import (
	"database/sql"
	"errors"
	"fmt"
	"html"
	"net/http"
	"net/url"

	"github.com/LeKovr/num2word"
	"github.com/shopspring/decimal"

	billgen_types "github.com/sergeykochiev/billgen/types"
	. "github.com/sergeykochiev/curs/backend/gui"
	. "github.com/sergeykochiev/curs/backend/util"
	"gorm.io/gorm"
	. "maragu.dev/gomponents"
	_ "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

type OrderEntity struct {
	Id                            decimal.Decimal `gorm:"primaryKey"`
	Name                          string
	Client_name                   string
	Client_phone                  string
	Company_name                  sql.NullString
	Date_created                  string
	Date_ended                    sql.NullString
	Ended                         bool
	Creator_id                    decimal.Decimal
	UserEntity                    UserEntity                    `gorm:"foreignKey:Creator_id"`
	OrderItemFulfillmentEntities  []OrderItemFulfillmentEntity  `gorm:"foreignKey:Order_id"`
	OrderResourceSpendingEntities []OrderResourceSpendingEntity `gorm:"foreignKey:Order_id"`
}

func (e OrderEntity) GetEntityPageButtons() Group {
	return Group{
		If(e.Ended, Group{
			ButtonComponent("Создать счет", A, Href(fmt.Sprintf("/order/%d/bill", e.GetId()))),
			ButtonComponent("Создать накладную", A, Href(fmt.Sprintf("/order/%d/invoice", e.GetId()))),
		}),
		If(!e.Ended, ButtonComponent("Завершить сейчас", A, Href(fmt.Sprintf("/order/%d/end", e.GetId())))),
	}
}

func (e OrderEntity) GetBIL(db *gorm.DB) billgen_types.BillItemList {
	len := len(e.OrderItemFulfillmentEntities)
	var bia = make([]billgen_types.BillItem, len)
	var summ float64
	for i, item := range e.OrderItemFulfillmentEntities {
		bia[i].Name = item.ItemEntity.Name
		bia[i].Cost = item.ItemEntity.Cost_by_one
		bia[i].Count = int(item.Quantity_fulfilled)
		bia[i].One_is_called = item.ItemEntity.One_is_called
		bia[i].Summ = float64(bia[i].Count) * bia[i].Cost
		summ += bia[i].Summ
	}
	return billgen_types.BillItemList{
		Bia:        bia,
		Len:        len,
		Summ:       summ,
		SummString: num2word.RuMoney(float64(summ), true),
	}
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

func (e *OrderEntity) GetPreloadedDb(db *gorm.DB) *gorm.DB {
	return db.Joins("UserEntity").Preload("OrderItemFulfillmentEntities.ItemEntity.ItemResourceNeeds.ResourceEntity")
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
	return db
}

func (e OrderEntity) GetDataRow() Group {
	return Group{
		TableDataComponent(html.EscapeString(fmt.Sprintf("%d", e.GetId())), Td, fmt.Sprintf("/order/%d", e.GetId())),
		TableDataComponent(e.Name, Td, ""),
		TableDataComponent(e.Client_name, Td, ""),
		TableDataComponent(e.Client_phone, Td, ""),
		TableDataComponent(e.Company_name.String, Td, ""),
		TableDataComponent(e.Date_created, Td, ""),
		TableDataComponent(ConditionalArg(e.Ended, "Да", "Нет"), Td, ""),
		TableDataComponent(ConditionalArg(e.Date_ended.Valid, e.Date_ended.String, "-"), Td, ""),
		TableDataComponent(e.UserEntity.Name, Td, ""),
	}
}

func (e OrderEntity) GetTableHeader() Group {
	return Group{
		TableDataComponent("ID", Th, ""),
		TableDataComponent("Название", Th, ""),
		TableDataComponent("Имя клиента", Th, ""),
		TableDataComponent("Телефон клиента", Th, ""),
		TableDataComponent("Компания клиента", Th, ""),
		TableDataComponent("Дата создания", Th, ""),
		TableDataComponent("Завершен", Th, ""),
		TableDataComponent("Дата завершения", Th, ""),
		TableDataComponent("Создатель", Th, ""),
	}
}

func (e OrderEntity) GetEntityPage(recursive bool) Group {
	return Group{
		LabeledFieldComponent("Название", e.Name),
		LabeledFieldComponent("Имя клиента", e.Client_name),
		LabeledFieldComponent("Телефон клиента", e.Client_phone),
		LabeledFieldComponent("Компания клиента", ConditionalArg(e.Company_name.Valid, e.Company_name.String, "-")),
		LabeledFieldComponent("Дата создания", e.Date_created),
		LabeledFieldComponent("Завершен", ConditionalArg(e.Ended, "ДА", "НЕТ")),
		LabeledFieldComponent("Дата завершения", ConditionalArg(e.Date_ended.Valid, e.Date_ended.String, "-")),
		If(recursive, Div(
			Class("bg-gray-100 flex flex-col gap-[8px] p-[8px]"),
			H2(Text(fmt.Sprintf("Создал пользователь #%d", e.UserEntity.GetId()))),
			LabeledFieldComponent("Имя", e.UserEntity.Name),
		)),
		If(recursive, RelationCardArrComponent("Предоставленные товары", e.OrderItemFulfillmentEntities, func(ent OrderItemFulfillmentEntity) Node {
			return RelationCardCoreComponent(GetOneReadableName(ent), GetOneHref(ent), Group{
				ent.GetEntityPage(false),
				ent.ItemEntity.GetEntityPage(false),
			})
		})),
		If(recursive, RelationCardArrComponent("Потраченные ресурсы", e.OrderResourceSpendingEntities, func(ent OrderResourceSpendingEntity) Node {
			return RelationCardCoreComponent(GetOneReadableName(ent), GetOneHref(ent), Group{
				ent.GetEntityPage(false),
				ent.ResourceEntity.GetEntityPage(false),
			})
		})),
	}
}

func (e OrderEntity) GetCreateForm(db *gorm.DB) Group {
	return Group{
		LabeledInputComponent("text", "", "name", "Название заказа", "", true),
		LabeledInputComponent("text", "", "client_name", "Имя клиента", "", true),
		LabeledInputComponent("number", "", "client_phone", "Телефон клиента", "", true),
		LabeledInputComponent("text", "", "company_name", "Компания клиента", "", false),
		LabeledInputComponent("date", "", "date_created", "Дата создания", "", true),
	}
}

func (e OrderEntity) GetReadableName() string {
	return "Заказ"
}

func (e OrderEntity) GetId() int64 {
	return e.Id.IntPart()
}

func (e *OrderEntity) Clear() {
	*e = OrderEntity{}
}

func (e *OrderEntity) SetId(id int64) {
	e.Id = decimal.NewFromInt(id)
}

func (e *OrderEntity) TableName() string {
	return "order"
}

func (e *OrderEntity) Validate() bool {
	return len(e.Client_phone) == 11
}

func (e *OrderEntity) ValidateAndParseForm(r *http.Request) error {
	form := r.Form
	if !form.Has("name") || !form.Has("client_name") || !form.Has("client_phone") || !form.Has("date_created") {
		return errors.New("Invalid fields")
	}
	if form.Has("company_name") {
		e.Company_name.String = form.Get("company_name")
		e.Company_name.Valid = true
	}
	e.Client_name = form.Get("client_name")
	e.Client_phone = form.Get("client_phone")
	e.Name = form.Get("name")
	e.Date_created = form.Get("date_created")
	e.Creator_id = r.Context().Value("user").(UserEntity).Id
	return nil
}
