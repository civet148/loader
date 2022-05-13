package loader

import (
	"github.com/civet148/log"
	"github.com/civet148/sqlca/v2"
)

func InsertConfig(db *sqlca.Engine, do *RunConfigDO) (err error) {
	_, err = db.Model(&do).Table(TableNameRunConfig).Insert()
	if err != nil {
		err = log.Errorf(err.Error())
		return err
	}
	return nil
}

func GetConfigCount(db *sqlca.Engine, strConfigName string) (count int64, err error) {
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