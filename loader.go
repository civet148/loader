package loader

import (
	"fmt"
	"github.com/civet148/log"
	"github.com/civet148/sqlca/v2"
	"github.com/urfave/cli/v2"
	"reflect"
)

const (
	json_tag = "json"
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

func GetConfig(strDSN, strConfigName string, model interface{}, cctx *cli.Context, flags ...string) error {
	db, err := sqlca.NewEngine(strDSN)
	if err != nil {
		log.Errorf(err.Error())
		return err
	}

	db.Debug(true)

	var id int
	var count int64
	count, err = db.Model(&id).
		Table(TableNameRunConfig).
		Select(RUN_CONFIG_COLUMN_ID).
		Equal(RUN_CONFIG_COLUMN_CONFIG_NAME, strConfigName).
		Limit(1).
		Query()
	if err != nil {

		//create table if not exist
		_, _, err = db.ExecRaw(table_sql)
		if err != nil {
			return log.Errorf(err.Error())
		}
		log.Infof("table %s not exist, auto create it [OK]", TableNameRunConfig)
	}

	values, _ := parseModelValues(model)
	log.Infof("config values [%+v]", values)
	if count == 0 {
		//TODO initial run config params
		log.Warnf("TODO initial run config params")
	} else {
		//TODO read run config params from database
	}
	return nil
}

func parseModelValues(model interface{}) (map[string]interface{}, error) {
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
			values = parseStructField(typ, val)
		}
	default:
		{
			return nil, fmt.Errorf("type of %v not support yet", typ.Kind())
		}
	}
	return values, nil
}

// parse struct fields
func parseStructField(typ reflect.Type, val reflect.Value) map[string]interface{} {
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
		tagVal := getTag(typField, json_tag)
		values[tagVal] = valField.Interface()
	}

	return values
}

// get struct field's tag value
func getTag(sf reflect.StructField, tagName string) string {
	return sf.Tag.Get(tagName)
}
