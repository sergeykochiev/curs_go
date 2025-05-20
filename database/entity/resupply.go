package entity

import (
	"errors"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"strconv"

	. "github.com/sergeykochiev/curs/backend/gui"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	. "maragu.dev/gomponents"
	_ "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

type ResourceResupplyEntity struct {
	Id             decimal.Decimal `gorm:"primaryKey"`
	Resource_id    decimal.Decimal
	Quantity_added float64
	Date           string
	ResourceEntity ResourceEntity `gorm:"foreignKey:Resource_id"`
}

func (e ResourceResupplyEntity) GetEntityPageButtons() Group {
	return Group{}
}

func (e ResourceResupplyEntity) GetFilters() Group {
	return Group{
		DateFilterComponent("Дата в диапазоне", "date"),
		StringFilterComponent("Название ресурса включает", "resource_name"),
	}
}

func (e *ResourceResupplyEntity) GetPreloadedDb(db *gorm.DB) *gorm.DB {
	return db.Joins("ResourceEntity")
}

func (e *ResourceResupplyEntity) GetFilteredDb(filters url.Values, db *gorm.DB) *gorm.DB {
	if filters.Has("date_lo") && filters.Get("date_lo") != "" {
		db = db.Where("date > ?", filters.Get("date_lo"))
	}
	if filters.Has("date_hi") && filters.Get("date_hi") != "" {
		db = db.Where("date < ?", filters.Get("date_hi"))
	}
	if filters.Has("resource_name") && filters.Get("resource_name") != "" {
		db = db.Where("ResourceEntity__name LIKE ?", "%"+filters.Get("resource_name")+"%")
	}
	return db
}

func (e ResourceResupplyEntity) GetDataRow() Group {
	return Group{
		TableDataComponent(html.EscapeString(fmt.Sprintf("%d", e.GetId())), Td, fmt.Sprintf("/resource_resupply/%d", e.GetId())),
		TableDataComponent(e.ResourceEntity.Name, Td, ""),
		TableDataComponent(fmt.Sprintf("%f", e.ResourceEntity.Cost_by_one), Td, ""),
		TableDataComponent(e.Date, Td, ""),
		TableDataComponent(fmt.Sprintf("%f", e.Quantity_added), Td, ""),
	}
}

func (e ResourceResupplyEntity) GetTableHeader() Group {
	return Group{
		TableDataComponent("Id", Th, ""),
		TableDataComponent("Название", Th, ""),
		TableDataComponent("Цена за единицу", Th, ""),
		TableDataComponent("Единица", Th, ""),
		TableDataComponent("Количество", Th, ""),
	}
}

func (e ResourceResupplyEntity) GetEntityPage(recursive bool) Group {
	return Group{
		LabeledFieldComponent("Количество добавлено (единиц)", fmt.Sprintf("%f", e.Quantity_added)),
		LabeledFieldComponent("Дата поставки", e.Date),
		If(recursive, Group{
			RelationCardComponent(fmt.Sprintf("Поставлен ресурс #%d", e.ResourceEntity.GetId()), &e.ResourceEntity),
		}),
	}
}

func (e ResourceResupplyEntity) GetCreateForm(db *gorm.DB) Group {
	var res []*ResourceEntity
	db.Find(&res)
	return Group{
		SelectComponent(res, "", func(r *ResourceEntity) string { return r.Name }, "Выберите ресурс", "resource_id", true, -1),
		LabeledInputComponent("number", "", "quantity_added", "Кол-во добавлено", "", true),
		LabeledInputComponent("date", "", "date", "Дата поставки", "", true),
	}
}

func (e ResourceResupplyEntity) GetReadableName() string {
	return "Поставка ресурса"
}

func (e *ResourceResupplyEntity) Validate() bool {
	return true
}

func (e ResourceResupplyEntity) GetId() int64 {
	return e.Id.IntPart()
}

func (e *ResourceResupplyEntity) Clear() {
	*e = ResourceResupplyEntity{}
}

func (e *ResourceResupplyEntity) SetId(id int64) {
	e.Id = decimal.NewFromInt(id)
}

func (e ResourceResupplyEntity) TableName() string {
	return "resource_resupply"
}

func (e *ResourceResupplyEntity) ValidateAndParseForm(r *http.Request) error {
	form := r.Form
	if !form.Has("resource_id") || !form.Has("quantity_added") || !form.Has("date") {
		return errors.New("Invalid fields")
	}
	var err error
	e.Resource_id, err = decimal.NewFromString(form.Get("resource_id"))
	if err != nil {
		return err
	}
	quantity_added, err := strconv.ParseFloat(form.Get("quantity_added"), 32)
	if err != nil {
		return err
	}
	e.Quantity_added = float64(quantity_added)
	e.Date = form.Get("date")
	return nil
}

func (e *ResourceResupplyEntity) AfterCreate(tx *gorm.DB) (err error) {
	e.ResourceEntity.Id = e.Resource_id
	res := tx.First(&e.ResourceEntity)
	if res.Error != nil {
		return res.Error
	}
	e.ResourceEntity.Quantity += e.Quantity_added
	res = tx.Updates(&e.ResourceEntity)
	if res.Error != nil {
		return res.Error
	}
	return
}
