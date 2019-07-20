// Copyright 2017-2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gogf/gf/g/os/glog"
	"github.com/gogf/gf/g/os/gtime"
	"github.com/gogf/gf/g/text/gregex"
	"github.com/gogf/gf/g/text/gstr"
	"github.com/gogf/gf/g/util/gconv"
)

// Type assert api for String().
type apiString interface {
	String() string
}

// 格式化SQL语句。
// 1. 支持参数只传一个slice；
// 2. 支持占位符号数量自动扩展；
func formatQuery(query string, args []interface{}) (newQuery string, newArgs []interface{}) {
	newQuery = query
	// 查询条件参数处理，主要处理slice参数类型
	if len(args) > 0 {
		for index, arg := range args {
			rv := reflect.ValueOf(arg)
			kind := rv.Kind()
			if kind == reflect.Ptr {
				rv = rv.Elem()
				kind = rv.Kind()
			}
			switch kind {
			// '?'占位符支持slice类型, 这里会将slice参数拆散，并更新原有占位符'?'为多个'?'，使用','符号连接。
			case reflect.Slice, reflect.Array:
				for i := 0; i < rv.Len(); i++ {
					newArgs = append(newArgs, rv.Index(i).Interface())
				}
				// 如果参数直接传递slice，并且占位符数量与slice长度相等，
				// 那么不用替换扩展占位符数量，直接使用该slice作为查询参数
				if len(args) == 1 && gstr.Count(newQuery, "?") == rv.Len() {
					break
				}
				// counter用于匹配该参数的位置(与index对应)
				counter := 0
				newQuery, _ = gregex.ReplaceStringFunc(`\?`, newQuery, func(s string) string {
					counter++
					if counter == index+1 {
						return "?" + strings.Repeat(",?", rv.Len()-1)
					}
					return s
				})
			default:
				newArgs = append(newArgs, arg)
			}
		}
	}
	return
}

// 格式化Where查询条件。
func formatWhere(where interface{}, args []interface{}) (newWhere string, newArgs []interface{}) {
	// 条件字符串处理
	buffer := bytes.NewBuffer(nil)
	// 使用反射进行类型判断
	rv := reflect.ValueOf(where)
	kind := rv.Kind()
	if kind == reflect.Ptr {
		rv = rv.Elem()
		kind = rv.Kind()
	}
	switch kind {
	// map/struct类型
	case reflect.Map:
		fallthrough
	case reflect.Struct:
		for key, value := range structToMap(where) {
			if buffer.Len() > 0 {
				buffer.WriteString(" AND ")
			}
			// 支持slice键值/属性，如果只有一个?占位符号，那么作为IN查询，否则打散作为多个查询参数
			rv := reflect.ValueOf(value)
			switch rv.Kind() {
			case reflect.Slice:
				fallthrough
			case reflect.Array:
				count := gstr.Count(key, "?")
				if count == 0 {
					buffer.WriteString(key + " IN(?)")
					newArgs = append(newArgs, value)
				} else if count != rv.Len() {
					buffer.WriteString(key)
					newArgs = append(newArgs, value)
				} else {
					buffer.WriteString(key)
					// 如果键名/属性名称中带有多个?占位符号，那么将参数打散
					newArgs = append(newArgs, gconv.Interfaces(value)...)
				}
			default:
				if value == nil {
					buffer.WriteString(key)
				} else {
					// 支持key带操作符号
					if gstr.Pos(key, "?") == -1 {
						if gstr.Pos(key, "<") == -1 && gstr.Pos(key, ">") == -1 && gstr.Pos(key, "=") == -1 {
							buffer.WriteString(key + "=?")
						} else {
							buffer.WriteString(key + "?")
						}
					} else {
						buffer.WriteString(key)
					}
					newArgs = append(newArgs, value)
				}
			}
		}

	default:
		buffer.WriteString(gconv.String(where))
	}
	// 没有任何条件查询参数，直接返回
	if buffer.Len() == 0 {
		return "", args
	}
	newArgs = append(newArgs, args...)
	newWhere = buffer.String()
	// 查询条件参数处理，主要处理slice参数类型
	if len(newArgs) > 0 {
		// 支持例如 Where/And/Or("uid", 1) 这种格式
		if gstr.Pos(newWhere, "?") == -1 {
			if gstr.Pos(newWhere, "<") == -1 && gstr.Pos(newWhere, ">") == -1 && gstr.Pos(newWhere, "=") == -1 {
				newWhere += "=?"
			} else {
				newWhere += "?"
			}
		}
	}
	return
}

// 将预处理参数转换为底层数据库引擎支持的格式。
// 主要是判断参数是否为复杂数据类型，如果是，那么转换为基础类型。
func convertParam(value interface{}) interface{} {
	rv := reflect.ValueOf(value)
	kind := rv.Kind()
	if kind == reflect.Ptr {
		rv = rv.Elem()
		kind = rv.Kind()
	}
	switch kind {
	case reflect.Struct:
		// 底层数据库引擎支持 time.Time/*time.Time 类型
		if v, ok := value.(time.Time); ok {
			if v.IsZero() {
				return "null"
			}
			return value
		}
		if v, ok := value.(*time.Time); ok {
			if v.IsZero() {
				return ""
			}
			return value
		}
		return gconv.String(value)
	}
	return value
}

// 打印SQL对象(仅在debug=true时有效)
func printSql(v *Sql) {
	s := fmt.Sprintf("%s, %v, %s, %s, %d ms, %s", v.Sql, v.Args,
		gtime.NewFromTimeStamp(v.Start).Format("Y-m-d H:i:s.u"),
		gtime.NewFromTimeStamp(v.End).Format("Y-m-d H:i:s.u"),
		v.End-v.Start,
		v.Func,
	)
	if v.Error != nil {
		s += "\nError: " + v.Error.Error()
		glog.Stack(true, 2).Error(s)
	} else {
		glog.Debug(s)
	}
}

// 格式化错误信息
func formatError(err error, query string, args ...interface{}) error {
	if err != nil && err != sql.ErrNoRows {
		errStr := fmt.Sprintf("DB ERROR: %s\n", err.Error())
		errStr += fmt.Sprintf("DB QUERY: %s\n", query)
		if len(args) > 0 {
			errStr += fmt.Sprintf("DB PARAM: %v\n", args)
		}
		err = errors.New(errStr)
	}
	return err
}

// 根据insert选项获得操作名称
func getInsertOperationByOption(option int) string {
	operator := "INSERT"
	switch option {
	case OPTION_REPLACE:
		operator = "REPLACE"
	case OPTION_SAVE:
	case OPTION_IGNORE:
		operator = "INSERT IGNORE"
	}
	return operator
}

// 将对象转换为map，如果对象带有继承对象，那么执行递归转换。
// 该方法用于将变量传递给数据库执行之前。
func structToMap(obj interface{}) map[string]interface{} {
	data := gconv.Map(obj)
	for key, value := range data {
		rv := reflect.ValueOf(value)
		kind := rv.Kind()
		if kind == reflect.Ptr {
			rv = rv.Elem()
			kind = rv.Kind()
		}
		switch kind {
		case reflect.Struct:
			// 底层数据库引擎支持 time.Time/*time.Time 类型
			if _, ok := value.(time.Time); ok {
				continue
			}
			if _, ok := value.(*time.Time); ok {
				continue
			}
			// 如果执行String方法，那么执行字符串转换
			if s, ok := value.(apiString); ok {
				data[key] = s.String()
				continue
			}
			delete(data, key)
			for k, v := range structToMap(value) {
				data[k] = v
			}
		}
	}
	return data
}

// 使用递归的方式将map键值对映射到struct对象上，注意参数<pointer>是一个指向struct的指针。
func mapToStruct(data map[string]interface{}, pointer interface{}) error {
	return gconv.StructDeep(data, pointer)
}
