package pkg

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

func GeneratePlaceholderMap(v interface{}, keys map[string]interface{}, prefix string) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(v).Elem()
	typ := val.Type()

	if val.Kind() == reflect.String {
		name := fmt.Sprintf("%s%s", prefix, strings.ToUpper(ToSnakeCase(typ.Name())))
		envPlaceholder := fmt.Sprintf("${%s}", name)
		keys[name] = val.Interface()
		result[typ.Name()] = envPlaceholder

		return result
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		fieldName := fieldType.Tag.Get("yaml")
		if fieldName == "" {
			fieldName = strings.ToLower(fieldType.Name)
		}
		switch field.Kind() {
		case reflect.Struct:
			nestedMap := GeneratePlaceholderMap(field.Addr().Interface(), keys, prefix+strings.ToUpper(fieldType.Name)+"_")
			result[fieldName] = nestedMap
		// case reflect.Slice:
		// 	for j := 0; j < field.Len(); j++ {
		// 		nestedMap := GeneratePlaceholderMap(field.Index(j).Addr().Interface(), keys, prefix+strings.ToUpper(fieldType.Name)+"_"+strconv.Itoa(j)+"_")
		// 		result[fieldName] = nestedMap
		// 	}
		default:
			name := fmt.Sprintf("%s%s", prefix, strings.ToUpper(ToSnakeCase(fieldType.Name)))
			envPlaceholder := fmt.Sprintf("${%s}", name)
			keys[name] = field.Interface()
			result[fieldName] = envPlaceholder
		}
	}

	return result
}

// https://stackoverflow.com/questions/56616196/how-to-convert-camel-case-string-to-snake-case
var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
