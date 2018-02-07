package willie

import (
	"fmt"
	"net/url"
	"reflect"
)

//EncodeToFormValues encodes structs into forms
func EncodeToFormValues(body interface{}) url.Values {
	result := url.Values{}

	rtype := reflect.TypeOf(body)
	rvalue := reflect.ValueOf(body)

	if rtype.Kind() == reflect.Ptr {
		rvalue = rvalue.Elem()
		rtype = reflect.TypeOf(rvalue.Interface())
	}

	for i := 0; i < rtype.NumField(); i++ {
		field := rtype.Field(i)
		value := rvalue.Field(i)

		switch field.Type.Kind() {
		case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
			result.Add(field.Name, fmt.Sprintf("%v", value))
		case reflect.Slice, reflect.Array:
			for index := 0; index < value.Len(); index++ {
				k := fmt.Sprintf("%v[%v]", field.Name, index)
				v := fmt.Sprintf("%v", value.Index(index))

				result.Add(k, v)
			}
		case reflect.Map:
			keys := value.MapKeys()
			for _, key := range keys {
				v := fmt.Sprintf("%v", value.MapIndex(key))
				k := fmt.Sprintf("%v.%v", field.Name, key)

				result.Add(k, v)
			}
		case reflect.Struct:
			merged := EncodeToFormValues(value.Interface())
			for mk := range merged {
				k := fmt.Sprintf("%v.%v", field.Name, mk)
				v := fmt.Sprintf("%v", merged.Get(mk))

				result.Add(k, v)
			}
		}
	}

	return result
}
