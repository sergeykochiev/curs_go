package util

import (
	"os"
	"time"

	billgen_types "github.com/sergeykochiev/billgen/types"
)

func GetCompanyInfoFromEnv() billgen_types.CompanyInfo {
	return billgen_types.CompanyInfo{
		Inn:     os.Getenv("COMP_NAME"),
		Name:    os.Getenv("COMP_INN"),
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
