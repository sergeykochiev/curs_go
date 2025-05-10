package database

import (
	"context"
	"reflect"

	"github.com/shopspring/decimal"
	"gorm.io/gorm/schema"
)

// JSONSerializer json serializer
type DecimalIdSerializer struct {
}

// Scan implements serializer interface
func (DecimalIdSerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) (err error) {
	field.ReflectValueOf(ctx, dst).Set(reflect.ValueOf(decimal.NewFromInt(dbValue.(int64))))
	return
}

// Value implements serializer interface
func (DecimalIdSerializer) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	return fieldValue.(decimal.Decimal).IntPart(), nil
}
