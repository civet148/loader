# loader
program config load from db and overwrite by CLI

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
+----+----------------+------------+-------------------------------------------+
| id | config_name    | config_key | config_value                              |
+----+----------------+------------+-------------------------------------------+
|  1 | user-service   | http_addr  | "0.0.0.0:80"                              |
|  2 | user-service   | static     | "/var/www/html"                           |
|  3 | user-service   | image_path | "/data/images"                            |
|  4 | user-service   | domain     | "http://user.mydomain.com/images"         |
|  5 | user-service   | debug      | true                                      |                                                                |
+----+----------------+------------+-------------------------------------------+

```