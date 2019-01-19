package main

import (
	"encoding/json"
	"fmt"
	"github.com/perlin-network/safu-go/api"
	"github.com/perlin-network/safu-go/log"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"
	"os"
	"sort"
	"time"
)

// Config describes how to start the node
type Config struct {
	PrivateKeyFile string
	Host           string
	Port           uint
	DatabasePath   string
	ResetDatabase  bool
}

func main() {
	app := cli.NewApp()

	app.Name = "taint_server"
	app.Author = "Perlin Network"
	app.Version = "0.0.1"
	app.Usage = "A server that processes taint requests"

	app.Flags = []cli.Flag{
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "host",
			Value: "localhost",
			Usage: "Listen for peers on host address `HOST`.",
		}),
		// note: use IntFlag for numbers, UintFlag don't seem to work with the toml files
		altsrc.NewIntFlag(cli.IntFlag{
			Name:  "port",
			Value: 3000,
			Usage: "Listen for peers on port `PORT`.",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "db.path",
			Value: "testdb",
			Usage: "Load/initialize LevelDB store from `DB_PATH`.",
		}),
		altsrc.NewBoolFlag(cli.BoolFlag{
			Name:  "db.reset",
			Usage: "Clear out the existing data in the datastore before initializing the DB.",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "private_key_file",
			Value: "wallet.txt",
			Usage: "TXT file that contain's the node's private key `PRIVATE_KEY_FILE`. Leave `PRIVATE_KEY_FILE` = 'random' if you want to randomly generate a wallet.",
		}),
	}

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("Version:    %s\n", c.App.Version)
		fmt.Printf("Built:      %s\n", c.App.Compiled.Format(time.ANSIC))
	}

	app.Action = func(c *cli.Context) error {
		config := &Config{
			PrivateKeyFile: c.String("private_key_file"),
			Host:           c.String("host"),
			Port:           c.Uint("port"),
			DatabasePath:   c.String("db.path"),
			ResetDatabase:  c.Bool("db.reset"),
		}

		// start the plugin
		if err := runServer(config); err != nil {
			return err
		}

		return nil
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse configuration/command-line arugments.")
	}
}

func runServer(c *Config) error {
	jsonConfig, _ := json.MarshalIndent(c, "", "  ")
	log.Debug().Msgf("Config: %s", string(jsonConfig))

	api.Run(fmt.Sprintf("%s:%d", c.Host, c.Port))

	return nil
}
