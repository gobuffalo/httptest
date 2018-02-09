package willie

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

//EncodeToURLValues encodes structs into url.Values
func EncodeToURLValues(body interface{}) (url.Values, error) {
	result := url.Values{}

	rtype := reflect.TypeOf(body)
	rvalue := reflect.ValueOf(body)

	if rtype.Kind() == reflect.Ptr {
		rvalue = rvalue.Elem()
		rtype = reflect.TypeOf(rvalue.Interface())
	}

	if rtype.Kind() != reflect.Map && rtype.Kind() != reflect.Struct {
		return result, errors.New("cannot use passed type to build url.Values")
	}

	addValues(rtype, rvalue, result)
	return result, nil
}

func addValues(rtype reflect.Type, rvalue reflect.Value, values url.Values) error {
	if rtype.Kind() == reflect.Map {
		keys := rvalue.MapKeys()
		for _, key := range keys {
			addFromField(key.String(), key.Type(), rvalue.MapIndex(key), values)
		}

		return nil
	}

	for i := 0; i < rtype.NumField(); i++ {
		field := rtype.Field(i)
		value := rvalue.Field(i)
		addFromField(field.Name, field.Type, value, values)
	}

	return nil
}

func addFromField(namespace string, rtype reflect.Type, rvalue reflect.Value, values url.Values) {
	switch rtype.Kind() {
	case reflect.String:
		values.Add(namespace, fmt.Sprintf("%v", rvalue))
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		values.Add(namespace, strconv.FormatInt(rvalue.Int(), 10))
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		values.Add(namespace, strconv.FormatUint(rvalue.Uint(), 10))
	case reflect.Float32:
		values.Add(namespace, strconv.FormatFloat(rvalue.Float(), 'f', -1, 32))
	case reflect.Float64:
		values.Add(namespace, strconv.FormatFloat(rvalue.Float(), 'f', -1, 64))
	case reflect.Bool:
		values.Add(namespace, strconv.FormatBool(rvalue.Bool()))
	case reflect.Slice, reflect.Array:
		for index := 0; index < rvalue.Len(); index++ {
			k := fmt.Sprintf("%v[%v]", namespace, index)
			value := rvalue.Index(index)
			field := reflect.TypeOf(value.Interface())
			addFromField(k, field, value, values)
		}
	case reflect.Map:
		for _, key := range rvalue.MapKeys() {
			k := fmt.Sprintf("%v[%v]", namespace, key)
			value := rvalue.MapIndex(key)
			field := reflect.TypeOf(value.Interface())

			addFromField(k, field, value, values)
		}
	case reflect.Struct:
		if rvalue.Type() == reflect.TypeOf(time.Time{}) {
			val := rvalue.Interface().(time.Time).Format(time.RFC3339)
			values.Add(namespace, val)
			return
		}

		for i := 0; i < rtype.NumField(); i++ {
			value := rvalue.Field(i)
			field := rtype.Field(i)

			k := fmt.Sprintf("%v.%v", namespace, field.Name)
			addFromField(k, field.Type, value, values)
		}
	}
}
