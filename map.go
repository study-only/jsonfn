package jsonfn

import (
	"errors"
	"strings"
	"reflect"
	"encoding/json"
)

func Marshal(entity interface{}, fields[]string) ([]byte, error) {
	node := parseFields(fields)
	m, err := loadRelation(entity, node)
	if err != nil {
		return nil, err
	}

	return json.Marshal(m)
}


func loadRelation(entity interface{}, node *node) (interface{}, error) {
	if isNil(entity) {
		return nil, nil
	}

	if reflect.TypeOf(entity).Kind() == reflect.Slice {
		s := reflect.ValueOf(entity)
		list := make([]interface{}, 0)
		for i := 0; i < s.Len(); i++ {
			itemData, err := loadRelation(s.Index(i).Interface(), node)
			if err != nil {
				return nil, err
			}

			list = append(list, itemData)
		}
		return list, nil
	}

	result, err := toMap(entity, node.GetFields())
	entityMap, ok := result.(map[string]interface{})
	if !ok {
		return nil, errors.New("unsupported entityMap")
	}

	if err != nil {
		return nil, nil
	}
	for _, child := range node.Children {
		if isNil(child) || child.IsLeaf() {
			continue
		}
		r := reflect.ValueOf(entity)
		methodValue := r.MethodByName(ucFirst(child.Name))
		methodType := methodValue.Type()
		if methodValue.IsValid() && methodType.NumIn() == 0 && methodType.NumOut() == 1 {
			ret := methodValue.Call([]reflect.Value{})
			if len(ret) == 1 {
				entityMap[child.Name], err = loadRelation(ret[0].Interface(), child)
			}
		}
	}

	return entityMap, nil
}

func toMap(entity interface{}, fields []string) (interface{}, error) {
	if reflect.TypeOf(entity).Kind() == reflect.Slice {
		s := reflect.ValueOf(entity)
		list := make([]interface{}, 0)
		for i := 0; i < s.Len(); i++ {
			itemData, err := toMap(s.Index(i).Interface(), fields)
			if err != nil {
				return nil, err
			}

			list = append(list, itemData)
		}
		return list, nil
	}

	bytes, err := json.Marshal(entity)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}

	for name := range result {
		if !inArray(name, fields) {
			delete(result, name)
		}
	}

	r := reflect.ValueOf(entity)
	if methodValue := r.MethodByName("ToMap"); methodValue.IsValid() {
		switch method := methodValue.Interface().(type) {
		case func(map[string]interface{}, []string) map[string]interface{}:
			return method(result, fields), nil
		}
	}

	return result, nil
}

func substr(str string, start, length int) string {

	rs := []rune(str)
	rl := len(rs)

	if rl == 0 {
		return ""
	}

	if start < 0 {
		start = rl + start
	}
	if start < 0 {
		start = 0
	}
	if start > rl-1 {
		return ""
	}

	end := rl

	if length < 0 {
		end = rl + length
	} else if length > 0 {
		end = start + length
	}

	if end < 0 || start >= end {
		return ""
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

func ucFirst(str string) string {
	return strings.ToUpper(substr(str, 0, 1)) + substr(str, 1, 0)
}

func isNil(a interface{}) bool {
	defer func() { recover() }()
	return a == nil || reflect.ValueOf(a).IsNil()
}

func inArray(needle string, array []string) bool {
	if len(array) == 0 {
		return true
	}
	for _, item := range array {
		if item == "*" {
			return true
		}
	}

	for _, item := range array {
		if needle == item {
			return true
		}
	}

	return false
}