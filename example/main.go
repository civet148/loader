package main

import (
	"github.com/civet148/loader"
	"github.com/civet148/log"
	"github.com/urfave/cli/v2"
	"os"
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

type Config struct {
	Debug     bool   `cli:"debug" json:"debug" db:"debug" toml:"debug"`
	HttpAddr  string `cli:"http-addr" json:"http_addr" db:"http_addr" toml:"http_addr"`
	Static    string `cli:"static" json:"static" db:"static" toml:"static"`
	ImagePath string `cli:"image-path" json:"image_path" db:"image_path" toml:"image_path"`
	Domain    string `cli:"domain" json:"domain" db:"domain" toml:"domain"`
	Timeout   int    `cli:"timeout" json:"timeout" db:"timeout" toml:"timeout"`
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
			}
			var strDSN = "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8"

			if cctx.Args().First() != "" {
				strDSN = cctx.Args().First()
			}

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
			log.Infof("config from db [%+v]", cfg)
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Errorf("exit in error %s", err)
		os.Exit(1)
		return
	}
}
