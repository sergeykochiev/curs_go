package report

import (
	"fmt"

	billgen_types "github.com/sergeykochiev/billgen/types"
	"github.com/sergeykochiev/curs/backend/util"
)

type ResourceSpending struct {
	Name          string
	Last_date     string
	One_is_called string
	Count_spent   float32
}

func (rs ResourceSpending) ToTRow() []billgen_types.TDData {
	return []billgen_types.TDData{
		{Value: rs.Name, Align: 1},
		{Value: rs.Last_date, Align: 1},
		{Value: rs.One_is_called, Align: 1},
		{Value: fmt.Sprintf("%f", rs.Count_spent), Align: 1},
	}
}

func (rs ResourceSpending) ToTHead() []billgen_types.THData {
	return []billgen_types.THData{
		{Value: "Название ресурса", Width: 0},
		{Value: "Дата последней траты", Width: 0},
		{Value: "Единица названа", Width: 0},
		{Value: "Кол-во потрачено", Width: 0},
	}
}

func (rs ResourceSpending) GetName() string {
	return "Траты ресурсов"
}

func (ip ResourceSpending) GetQuery(is_dl bool, is_dh bool) string {
	return fmt.Sprintf(`select resource.name as name, resource.one_is_called as one_is_called, max("order".date_ended) as last_date, count(order_resource_spending.id) as count_spent from resource left join order_resource_spending on resource.id = order_resource_spending.resource_id join "order" on order_resource_spending.order_id = "order".id where "order".ended = 1%s%s group by resource.id order by one_is_called desc, count_spent desc`, util.ConditionalArg(is_dl, ` and "order".date_ended > ?`, ""), util.ConditionalArg(is_dh, ` and "order".date_ended < ?`, ""))
}

func (rs *ResourceSpending) Clear() {
	*rs = ResourceSpending{}
}
