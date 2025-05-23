package entity

import (
	"errors"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"strconv"

	. "github.com/sergeykochiev/curs/backend/gui"
	"github.com/sergeykochiev/curs/backend/util"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	. "maragu.dev/gomponents"
	_ "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

type ResourceEntity struct {
	Id                            decimal.Decimal `gorm:"primaryKey"`
	Name                          string
	Date_last_updated             string
	Cost_by_one                   float64
	One_is_called                 string
	Quantity                      float64
	ResourceResupplyEntities      []ResourceResupplyEntity      `gorm:"foreignKey:Resource_id"`
	OrderResourceSpendingEntities []OrderResourceSpendingEntity `gorm:"foreignKey:Resource_id"`
	ItemResourceNeeds             []ItemResourceNeed            `gorm:"foreignKey:Resource_id"`
}

func (e ResourceEntity) GetEntityPageButtons() Group {
	return Group{}
}

func (e *ResourceEntity) GetFilters() Group {
	return Group{
		StringFilterComponent("Название включает", "name"),
		StringFilterComponent("Единицей является", "one_is_called"),
		DateFilterComponent("Дата последнего обновления в диапазоне", "date_last_updated"),
	}
}

func (e *ResourceEntity) GetPreloadedDb(db *gorm.DB) *gorm.DB {
	return db.Preload("ResourceResupplyEntities").Preload("OrderResourceSpendingEntities").Preload("ItemResourceNeeds.ItemEntity")
}

func (e *ResourceEntity) GetFilteredDb(filters url.Values, db *gorm.DB) *gorm.DB {
	if filters.Has("date_last_updated_lo") && filters.Get("date_last_updated_lo") != "" {
		db = db.Where("date_last_updated > ?", filters.Get("date_last_updated_lo"))
	}
	if filters.Has("date_last_updated_hi") && filters.Get("date_last_updated_hi") != "" {
		db = db.Where("date_last_updated < ?", filters.Get("date_last_updated_hi"))
	}
	if filters.Has("name") && filters.Get("name") != "" {
		db = db.Where("name LIKE ?", "%"+filters.Get("name")+"%")
	}
	if filters.Has("one_is_called") && filters.Get("one_is_called") != "" {
		db = db.Where("one_is_called = ?", filters.Get("one_is_called"))
	}
	return db
}

func (e *ResourceEntity) GetDataRow() Group {
	return Group{
		TableDataComponent(html.EscapeString(fmt.Sprintf("%d", e.GetId())), Td, fmt.Sprintf("/resource/%d", e.GetId())),
		TableDataComponent(e.Name, Td, ""),
		TableDataComponent(e.Date_last_updated, Td, ""),
		TableDataComponent(fmt.Sprintf("%f", e.Cost_by_one), Td, ""),
		TableDataComponent(e.One_is_called, Td, ""),
		TableDataComponent(fmt.Sprintf("%f", e.Quantity), Td, ""),
	}
}

func (e *ResourceEntity) GetTableHeader() Group {
	return Group{
		TableDataComponent("Id", Th, ""),
		TableDataComponent("Наименование", Th, ""),
		TableDataComponent("Дата обновления", Th, ""),
		TableDataComponent("Цена за единицу", Th, ""),
		TableDataComponent("Единица", Th, ""),
		TableDataComponent("Количество", Th, ""),
	}
}

func (e ResourceEntity) GetEntityPage(recursive bool) Group {
	return Group{
		LabeledFieldComponent("Наименование", e.Name),
		LabeledFieldComponent("Дата обновления", e.Date_last_updated),
		LabeledFieldComponent("Цена за единицу", fmt.Sprintf("%f", e.Cost_by_one)),
		LabeledFieldComponent("Единица", e.One_is_called),
		LabeledFieldComponent("Количество", fmt.Sprintf("%f", e.Quantity)),
		If(recursive, RelationCardArrComponent("Траты", e.OrderResourceSpendingEntities, func(ent OrderResourceSpendingEntity) Node {
			return RelationCardCoreComponent(util.GetOneReadableName(ent), util.GetOneHref(ent), Group{
				ent.GetEntityPage(false),
				ent.OrderEntity.GetEntityPage(false),
			})
		})),
		If(recursive, RelationCardArrComponent("Необходим товарам", e.ItemResourceNeeds, func(ent ItemResourceNeed) Node {
			return RelationCardCoreComponent(util.GetOneReadableName(ent), util.GetOneHref(ent), Group{
				ent.GetEntityPage(false),
				ent.ItemEntity.GetEntityPage(false),
			})
		})),
		If(recursive, RelationCardArrComponent("Поставки", e.ResourceResupplyEntities, func(ent ResourceResupplyEntity) Node {
			return RelationCardComponent(util.GetOneReadableName(ent), ent)
		})),
	}
}

func (e ResourceEntity) GetCreateForm(db *gorm.DB) Group {
	return Group{
		LabeledInputComponent("text", "", "name", "Название", "", true),
		LabeledInputComponent("text", "", "cost_by_one", "Стоимость за единицу", "", true),
		LabeledInputComponent("text", `По умолчанию - "Единица"`, "one_is_called", "Единица названа", "", false),
	}
}

func (e ResourceEntity) GetReadableName() string {
	return "Ресурс"
}

func (e ResourceEntity) GetId() int64 {
	return e.Id.IntPart()
}
func (e *ResourceEntity) Clear() {
	*e = ResourceEntity{}
}

func (e *ResourceEntity) SetId(id int64) {
	e.Id = decimal.NewFromInt(id)
}

func (e ResourceEntity) Validate() bool {
	return true
}

func (e ResourceEntity) TableName() string {
	return "resource"
}

func (e *ResourceEntity) ValidateAndParseForm(r *http.Request) error {
	form := r.Form
	if !form.Has("name") || !form.Has("cost_by_one") {
		return errors.New("Invalid fields")
	}
	if form.Has("one_is_called") {
		e.One_is_called = form.Get("one_is_called")
	}
	e.Name = form.Get("name")
	cost_by_one, err := strconv.ParseFloat(form.Get("cost_by_one"), 32)
	if err != nil {
		return err
	}
	e.Cost_by_one = float64(cost_by_one)
	quantity, err := strconv.ParseFloat(form.Get("quantity"), 32)
	if err != nil {
		return err
	}
	e.Quantity = float64(quantity)
	return nil
}
