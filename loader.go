package loader

import (
	"database/sql"
	"fmt"
	"github.com/civet148/log"
	"github.com/civet148/sqlca/v2"
	"github.com/urfave/cli/v2"
	"reflect"
	"strconv"
	"strings"
)

const (
	TagName_DB   = "db"
	TagName_JSON = "json"
	TagName_Bson = "bson"
	TagName_Toml = "toml"
	TagName_CLI  = "cli"
)

const (
	table_sql = "CREATE TABLE `run_config` (\n" +
		"`id` int NOT NULL AUTO_INCREMENT COMMENT 'incr id',\n  " +
		"`config_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'config name',\n " +
		"`config_key` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'config key',\n  " +
		"`config_value` text COLLATE utf8mb4_unicode_ci COMMENT 'config value',\n  " +
		"`remark` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'remark',\n  " +
		"`deleted` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'is deleted(0=false 1=true)',\n  " +
		"`created_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'created time',\n  " +
		"PRIMARY KEY (`id`),\n  UNIQUE KEY `UNIQ_NAME_KEY` (`config_name`,`config_key`)\n" +
		") ENGINE=InnoDB AUTO_INCREMENT=131 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='run config table';\n"
)

//Configure initialize or load run config from database
func Configure(ctx *cli.Context, strDSN, strConfigName string, model interface{}) error {

	if ctx == nil && strDSN == "" {
		return fmt.Errorf("CLI context and DSN is nil")
	}
	if strDSN != "" {
		db, err := sqlca.NewEngine(strDSN)
		if err != nil {
			log.Errorf(err.Error())
			return err
		}

		if err = daoInitTable(db, strConfigName); err != nil {
			return err
		}

		var params map[string]string
		params, err = daoGetConfigParams(db, strConfigName)
		values, _ := parseModelValues(model, TagName_DB)

		for k, v := range values {
			var strValue string
			switch v.(type) {
			case string:
				strValue = fmt.Sprintf("\"%v\"", v.(string))
			default:
				strValue = fmt.Sprintf("%v", v)
			}
			if _, ok := params[k]; ok {
				continue
			}
			var do = &RunConfigDO{
				ConfigName:  strConfigName,
				ConfigKey:   k,
				ConfigValue: strValue,
			}
			err = daoInsertConfig(db, do)
			if err != nil {
				err = log.Errorf(err.Error())
				return err
			}
		}

		//read run config params from database
		err = daoLoadConfig(db, strConfigName, model)
		if err != nil {
			err = log.Errorf("config json [%s] unmarshal error [%s]", err.Error())
			return err
		}
	}

	if ctx != nil {
		err := setCliValues(ctx, model)
		if err != nil {
			return err
		}
	}
	return nil
}

//overwrite model params by CLI flags
func setCliValues(ctx *cli.Context, model interface{}) error {
	values, err := parseModelValues(model, TagName_CLI)
	if err != nil {
		log.Errorf(err.Error())
		return err
	}
	if len(values) == 0 {
		return log.Errorf("no CLI tag found in model")
	}
	if err = setModelValues(ctx, model); err != nil {
		return err
	}
	return err
}

func parseModelValues(model interface{}, tag string) (map[string]interface{}, error) {
	var values map[string]interface{}
	typ := reflect.TypeOf(model)
	val := reflect.ValueOf(model)

	for {
		if typ.Kind() != reflect.Ptr { // pointer type
			break
		}
		typ = typ.Elem()
		val = val.Elem()
	}

	kind := typ.Kind()
	switch kind {
	case reflect.Struct:
		{
			values = parseStructField(typ, val, tag)
		}
	default:
		{
			return nil, fmt.Errorf("type of %v not support yet", typ.Kind())
		}
	}
	return values, nil
}

// parse struct fields
func parseStructField(typ reflect.Type, val reflect.Value, tag string) map[string]interface{} {
	var values = make(map[string]interface{})

	NumField := val.NumField()
	for i := 0; i < NumField; i++ {
		typField := typ.Field(i)
		valField := val.Field(i)

		if typField.Type.Kind() == reflect.Ptr {
			typField.Type = typField.Type.Elem()
			valField = valField.Elem()
		}
		if !valField.IsValid() || !valField.CanInterface() {
			continue
		}
		tagVal := getTag(typField, tag)
		if tagVal != "" {
			values[tagVal] = valField.Interface()
		}
	}

	return values
}

// get struct field's tag value
func getTag(sf reflect.StructField, tagName string) string {
	tagValue := sf.Tag.Get(tagName)
	return handleTagValue(tagName, tagValue)
}

func setModelValues(ctx *cli.Context, model interface{}) error {
	typ := reflect.TypeOf(model)
	val := reflect.ValueOf(model)

	for {
		if typ.Kind() != reflect.Ptr { // pointer type
			break
		}
		typ = typ.Elem()
		val = val.Elem()
	}

	kind := typ.Kind()
	switch kind {
	case reflect.Struct:
		{
			if err := setStructValue(ctx, typ, val); err != nil {
				return log.Errorf(err.Error())
			}
		}
	default:
		{
			return fmt.Errorf("type of %v not support yet", typ.Kind())
		}
	}
	return nil
}

// parse struct fields
func setStructValue(ctx *cli.Context, typ reflect.Type, val reflect.Value) (err error) {
	var flagValues = make(map[string]interface{})
	flags := ctx.FlagNames()
	for _, flag := range flags {
		flagValues[flag] = ctx.Value(flag)
	}
	NumField := val.NumField()
	for i := 0; i < NumField; i++ {
		typField := typ.Field(i)
		valField := val.Field(i)

		if typField.Type.Kind() == reflect.Ptr {
			typField.Type = typField.Type.Elem()
			valField = valField.Elem()
		}
		if !valField.IsValid() || !valField.CanInterface() {
			continue
		}
		tagVal := getTag(typField, TagName_CLI)
		if tagVal != "" {
			if v, ok := flagValues[tagVal]; ok {
				if err = setValue(typField.Type, valField, v); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func setValue(typ reflect.Type, val reflect.Value, v interface{}) (err error) {
	strValue := fmt.Sprintf("%v", v)
	return setValueString(typ, val, strValue)
}

//将string存储的值赋值到变量
func setValueString(typ reflect.Type, val reflect.Value, v string) (err error) {
	switch typ.Kind() {
	case reflect.Struct:
		s, ok := val.Addr().Interface().(sql.Scanner)
		if !ok {
			log.Warnf("struct type %s not implement sql.Scanner interface", typ.Name())
			return
		}
		if err := s.Scan(v); err != nil {
			panic(fmt.Sprintf("scan value %s to sql.Scanner implement object error [%s]", v, err))
		}
	case reflect.String:
		val.SetString(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, _ := strconv.ParseInt(v, 10, 64)
		val.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, _ := strconv.ParseUint(v, 10, 64)
		val.SetUint(i)
	case reflect.Float32, reflect.Float64:
		i, _ := strconv.ParseFloat(v, 64)
		val.SetFloat(i)
	case reflect.Bool:
		if v == "true" {
			val.SetBool(true)
		} else {
			val.SetBool(false)
		}
	case reflect.Ptr:
		typ = typ.Elem()
		err = setValueString(typ, val, v)
	default:
		err = fmt.Errorf("can't assign value [%v] to variant type [%v]\n", v, typ.Kind())
		return err
	}
	return nil
}

func handleTagValue(strTagName, strTagValue string) string {

	if strTagValue == "" {
		return ""
	}

	if strTagName == TagName_JSON {

		vs := strings.Split(strTagValue, ",")
		strTagValue = vs[0]
	} else {

	}
	return strTagValue
}
