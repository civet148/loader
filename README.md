# loader
program config load from db and overwirte by CLI

# example database

- run_config table create

```sql
CREATE TABLE `run_config` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'incr id',
  `config_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'config name',
  `config_key` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'config key',
  `config_value` text COLLATE utf8mb4_unicode_ci COMMENT 'config value',
  `remark` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'remark',
  `deleted` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'is deleted(0=false 1=true)',
  `created_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'created time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `UNIQ_NAME_KEY` (`config_name`,`config_key`)
) ENGINE=InnoDB AUTO_INCREMENT=131 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='run config table';

```

- run_config table records

```shell
+----+----------------+------------+----------------------------------------------------------------------+
| id | config_name    | config_key | config_value                                                         |
+----+----------------+------------+----------------------------------------------------------------------+
|  1 | user-backend   | http_addr  | "0.0.0.0:80"                                                         |
|  2 | user-backend   | static     | "/opt/static"                                                        |
|  3 | user-backend   | image_path | "/data/nft-printer/images"                                           |
|  4 | user-backend   | domain     | "http://dev-nftprinter-bcos.storeros.com/nft-printer/images"         |
|  5 | user-backend   | debug      | true                                                                 |
|  6 | user-backend   | ak         | "jIgMtVTYIfSA2bUteP"                                                 |
|  7 | user-backend   | sk         | "ncmTbkUMNz7nyZFL3DZSUOxvrDQcmIGt0PZI"                               |
|  8 | user-backend   | nft_url    | "https://dev-dcs-system.storeros.com/api/v1/chain/detail/collection" |
|  9 | user-backend   | bcos_url   | "http://192.168.20.108:8545"                                         |
| 10 | user-backend   | chain_id   | 1                                                                    |
| 11 | user-backend   | group_id   | 1                                                                    |
| 12 | system-backend | http_addr  | "0.0.0.0:80"                                                         |
| 13 | system-backend | static     | "/opt/static"                                                        |
| 14 | system-backend | image_path | "/data/nft-printer/images"                                           |
| 15 | system-backend | domain     | "http://dev-nftprinter-bcos.storeros.com/nft-printer/images"         |
| 16 | system-backend | debug      | true                                                                 |
+----+----------------+------------+----------------------------------------------------------------------+

```