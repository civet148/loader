package loader

import (
	"github.com/civet148/log"
	"github.com/urfave/cli/v2"
	"github.com/civet148/sqlca/v2"
)

const(
	table_sql= "CREATE TABLE `run_config` (\n  `id` int NOT NULL AUTO_INCREMENT COMMENT 'incr id',\n  `config_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'config name',\n  `config_key` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'config key',\n  `config_value` text COLLATE utf8mb4_unicode_ci COMMENT 'config value',\n  `remark` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'remark',\n  `deleted` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'is deleted(0=false 1=true)',\n  `created_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'created time',\n  PRIMARY KEY (`id`),\n  UNIQUE KEY `UNIQ_NAME_KEY` (`config_name`,`config_key`)\n) ENGINE=InnoDB AUTO_INCREMENT=131 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='run config table';\n"
)

func GetConfig(strDSN, strName string, cfg interface{}, cctx *cli.Context, flags...string) (err error) {
	db := sqlca.NewEngine(strDSN)
	db.Debug(true)

	var id int
	var count int64
	count, err = db.Model(&id).
		Table(TableNameRunConfig).
		Select(RUN_CONFIG_COLUMN_ID).
		Equal(RUN_CONFIG_COLUMN_DELETED, 0).
		Limit(1).
		Query()
	if err != nil {

		//create table if not exist
		_, _, err = db.ExecRaw(table_sql)
		if err != nil {
			return log.Errorf(err.Error())
		}
	}

	if count == 0 {
		//TODO initial run config params
		log.Warnf("TODO initial run config params")
	} else {
		//TODO read run config params from database
	}
	return nil
}

