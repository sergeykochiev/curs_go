package report

import (
	"fmt"

	billgen_types "github.com/sergeykochiev/billgen/types"
)

const ResourceSpendinQuery string = `select resource.name as name, max("order".date_ended) as last_date, sum(order_resource_spending.quantity_spent) * resource.cost_by_one as money_spent from resource left join order_resource_spending on resource.id = order_resource_spending.resource_id join "order" on order_resource_spending.order_id = "order".id where "order".ended = 1%s%s group by resource.id order by one_is_called desc, money_spent desc`

type ResourceSpending struct {
	Name        string
	Last_date   string
	Money_spent float64
}

func (rs ResourceSpending) ToTRow() []billgen_types.TDData {
	return []billgen_types.TDData{
		{Value: rs.Name, Align: 1},
		{Value: rs.Last_date, Align: 1},
		{Value: fmt.Sprintf("%f", rs.Money_spent), Align: 1},
	}
}

func (rs ResourceSpending) ToTHead() []billgen_types.THData {
	return []billgen_types.THData{
		{Value: "Название ресурса", Width: 0},
		{Value: "Дата последней траты", Width: 0},
		{Value: "Рублей потрачено", Width: 0},
	}
}

func (rs ResourceSpending) ToTFoot(rsa []ResourceSpending) []billgen_types.TDData {
	var sum float64
	for _, rsi := range rsa {
		sum += rsi.Money_spent
	}
	return []billgen_types.TDData{
		{Value: "", Align: 1},
		{Value: "Итого", Align: 1},
		{Value: fmt.Sprintf("%f", sum), Align: 1},
	}
}

func (rs ResourceSpending) GetName() string {
	return "Траты ресурсов"
}

func (ip ResourceSpending) GetQuery(is_dl bool, is_dh bool) string {
	return fmt.Sprintf(ResourceSpendinQuery, func() string {
		if is_dl {
			return ` and "order".date_ended > ?`
		}
		return ""
	}(), func() string {
		if is_dh {
			return ` and "order".date_ended < ?`
		}
		return ""
	}())
}

func (rs *ResourceSpending) Clear() {
	*rs = ResourceSpending{}
}
