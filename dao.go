package loader

import (
	"encoding/json"
	"github.com/civet148/log"
	"github.com/civet148/sqlca/v2"
)

func daoInsertConfig(db *sqlca.Engine, do *RunConfigDO) (err error) {
	_, err = db.Model(&do).Table(TableNameRunConfig).Insert()
	if err != nil {
		err = log.Errorf(err.Error())
		return err
	}
	return nil
}

func daoInitTable(db *sqlca.Engine, strConfigName string) (err error) {
	var id int
	_, err = db.Model(&id).
		NoVerbose().
		Table(TableNameRunConfig).
		Select(RUN_CONFIG_COLUMN_ID).
		Equal(RUN_CONFIG_COLUMN_CONFIG_NAME, strConfigName).
		Query()
	if err != nil {

		//create table if not exist
		_, _, err = db.ExecRaw(table_sql)
		if err != nil {
			return log.Errorf(err.Error())
		}
		log.Infof("table %s not exist, auto create it [OK]", TableNameRunConfig)
	}
	return nil
}

func daoGetConfigParams(db *sqlca.Engine, strConfigName string) (params map[string]string, err error) {

	var dos []*RunConfigDO
	params = make(map[string]string)

	_, err = db.Model(&dos).
		NoVerbose().
		Table(TableNameRunConfig).
		Select(RUN_CONFIG_COLUMN_CONFIG_KEY, RUN_CONFIG_COLUMN_CONFIG_VALUE).
		Equal(RUN_CONFIG_COLUMN_CONFIG_NAME, strConfigName).
		Query()
	if err != nil {

		//create table if not exist
		_, _, err = db.ExecRaw(table_sql)
		if err != nil {
			return nil, log.Errorf(err.Error())
		}
		log.Infof("table %s not exist, auto create it [OK]", TableNameRunConfig)
	}
	for _, do := range dos {
		params[do.ConfigKey] = do.ConfigValue
	}
	return params, nil
}

func daoLoadConfig(db *sqlca.Engine, strConfigName string, model interface{}) (err error) {
	/*
	 SELECT  CONCAT('{', GROUP_CONCAT('"', config_key, '":', config_value, '"'), '}') AS config FROM run_config  WHERE 1=1 AND config_name='user-backend';
	*/
	var strConfigJson string
	if _, err = db.Model(&strConfigJson).
		Table(TableNameRunConfig).
		Select("CONCAT('{', GROUP_CONCAT('\"', config_key, '\":', config_value), '}') AS config").
		Equal(RUN_CONFIG_COLUMN_CONFIG_NAME, strConfigName).
		Query(); err != nil {
		err = log.Errorf("load config from database error [%s]", err.Error())
		return err
	}
	if err = json.Unmarshal([]byte(strConfigJson), model); err != nil {
		err = log.Errorf("config json %s unmarshal error [%s]", strConfigJson, err.Error())
		return err
	}
	return nil
}
