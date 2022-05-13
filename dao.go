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

func daoGetConfigCount(db *sqlca.Engine, strConfigName string) (count int64, err error) {
	var id int

	count, err = db.Model(&id).
		NoVerbose().
		Table(TableNameRunConfig).
		Select(RUN_CONFIG_COLUMN_ID).
		Equal(RUN_CONFIG_COLUMN_CONFIG_NAME, strConfigName).
		Limit(1).
		Query()
	if err != nil {

		//create table if not exist
		_, _, err = db.ExecRaw(table_sql)
		if err != nil {
			return 0, log.Errorf(err.Error())
		}
		log.Infof("table %s not exist, auto create it [OK]", TableNameRunConfig)
	}
	return count, nil
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
		err = log.Errorf("config json [%s] unmarshal error [%s]", err.Error())
		return err
	}
	return nil
}