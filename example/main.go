package main

import (
	"github.com/civet148/loader"
	"github.com/civet148/log"
	"github.com/urfave/cli/v2"
	"os"
	"time"
)

const (
	PROGRAM_NAME = "loader"
	CONFIG_NAME  = "run-params"
)

const (
	CMD_FLAG_NAME_DEBUG      = "debug"
	CMD_FLAG_NAME_STATIC     = "static"
	CMD_FLAG_NAME_IMAGE_PATH = "image-path"
	CMD_FLAG_NAME_DOMAIN     = "domain"
)

type Person struct {
	Age  int    `json:"age"`
	Name string `json:"name"`
}
type Config struct {
	Debug     bool     `cli:"debug" json:"debug" db:"debug" toml:"debug"`
	HttpAddr  string   `cli:"http-addr" json:"http_addr" db:"http_addr" toml:"http_addr"`
	Static    string   `cli:"static" json:"static" db:"static" toml:"static"`
	ImagePath string   `cli:"image-path" json:"image_path" toml:"image_path"`
	Domain    string   `cli:"domain" json:"domain" toml:"domain"`
	Timeout   int      `cli:"timeout" json:"timeout" toml:"timeout"`
	Keys      []string `cli:"keys" json:"keys" toml:"keys"`
	Person    Person   `json:"person"`
}

func main() {
	app := &cli.App{
		Name:    PROGRAM_NAME,
		Usage:   "[data source name]",
		Version: "v1.0.0",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  CMD_FLAG_NAME_DEBUG,
				Usage: "open debug mode",
			},
			&cli.StringFlag{
				Name:  CMD_FLAG_NAME_STATIC,
				Usage: "frontend static path",
			},
			&cli.StringFlag{
				Name:  CMD_FLAG_NAME_IMAGE_PATH,
				Usage: "image path",
			},
			&cli.StringFlag{
				Name:  CMD_FLAG_NAME_DOMAIN,
				Usage: "domain setting",
			},
		},
		Commands: nil,
		Action: func(cctx *cli.Context) error {
			var cfg = &Config{
				Debug:     true,
				HttpAddr:  ":80",
				Static:    "/opt/static",
				ImagePath: "/data/images",
				Domain:    "https://www.mydomain.com",
				Keys:      []string{"key1", "key2"},
				Person: Person{
					Age:  18,
					Name: "John",
				},
			}
			var strDSN = "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8"

			if cctx.Args().First() != "" {
				strDSN = cctx.Args().First()
			}
			loader.SetRefreshInterval(5)
			err := loader.Configure(
				cctx,
				strDSN,
				CONFIG_NAME,
				cfg,
			)
			if err != nil {
				log.Errorf("load run config from db error [%s]", err)
				return err
			}
			for i := 0; i < 10000; i++ {
				log.Infof("config from db [%+v]", cfg)
				time.Sleep(5*time.Second)
			}
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Errorf("exit in error %s", err)
		os.Exit(1)
		return
	}
}
