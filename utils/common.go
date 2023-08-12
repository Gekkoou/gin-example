package utils

import (
	"context"
	"fmt"
	"gin-example/global"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strconv"
)

func Translate(error error) string {
	errs, ok := error.(validator.ValidationErrors)
	if !ok {
		return error.Error()
	}
	for _, val := range errs.Translate(global.Trans) {
		return val
	}
	return error.Error()
}

func GetValidateErr(error error, obj any) error {
	errs, ok := error.(validator.ValidationErrors)
	if !ok {
		return error
	}
	for _, err := range errs {
		if field, ok := reflect.TypeOf(obj).FieldByName(err.Field()); ok {
			if e := field.Tag.Get("err"); e != "" {
				return fmt.Errorf("%s: %s", err.Field(), e)
			}
		}
	}
	return errs
}

func DelCache(key string) {
	global.Redis.Del(context.Background(), key).Result()
}

func GetCacheKeyById(key string, id int) string {
	return key + strconv.Itoa(id)
}
