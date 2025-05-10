package util

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	billgen_types "github.com/sergeykochiev/billgen/types"
	"github.com/sergeykochiev/curs/backend/types"
)

func MakeArrayOf[T interface{}](i T) []T {
	return reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(i)), 0, 0).Interface().([]T)
}

func GetOneReadableName[T interface {
	types.Identifier
	types.HtmlTemplater
}](ent T) string {
	return fmt.Sprintf("%s #%d", ent.GetReadableName(), ent.GetId())
}

func GetOneHref[T interface {
	types.Identifier
}](ent T) string {
	return fmt.Sprintf("/%s/%d", ent.TableName(), ent.GetId())
}

func RunOnQ(q *chan func(), f func() error) error {
	err := make(chan error)
	*q <- func() {
		err <- f()
	}
	return <-err
}

func GetCompanyInfoFromEnv() billgen_types.CompanyInfo {
	return billgen_types.CompanyInfo{
		Inn:     os.Getenv("COMP_INN"),
		Name:    os.Getenv("COMP_NAME"),
		Address: os.Getenv("COMP_ADDRESS"),
		Number:  os.Getenv("COMP_NUMBER"),
		Details: billgen_types.CompanyDetails{
			Bank: os.Getenv("COMP_DET_BANK"),
			Rs:   os.Getenv("COMP_DET_RS"),
			Ks:   os.Getenv("COMP_DET_KS"),
			Bik:  os.Getenv("COMP_DET_BIK"),
		},
		PersonResp: os.Getenv("COMP_PERSON_RESP"),
	}
}

func ConditionalArg[T any](condition bool, arg T, notarg T) T {
	if condition {
		return arg
	}
	return notarg
}

func GetCurrentTime() string {
	return time.Now().Format(time.DateTime)
}

func GetCurrentDate() string {
	return time.Now().Format(time.DateOnly)
}

func GetRussianMonthGenitive(month int) string {
	switch month {
	case 1:
		return "января"
	case 2:
		return "февраля"
	case 3:
		return "марта"
	case 4:
		return "апреля"
	case 5:
		return "мая"
	case 6:
		return "июня"
	case 7:
		return "июля"
	case 8:
		return "августа"
	case 9:
		return "сентября"
	case 10:
		return "октября"
	case 11:
		return "ноября"
	case 12:
		return "декабря"
	}
	log.Fatal("They discovered a new month")
	return ""
}
