package main

import (
	"encoding/json"
	"fmt"
	"github.com/perlin-network/safu-go/cmd/safu_cli/client"
	"github.com/perlin-network/safu-go/log"
	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	app := cli.NewApp()

	app.Name = "safu_cli"
	app.Author = "Perlin Network"
	app.Version = "0.0.1"
	app.Usage = "A cli that submits taint requests"

	commonFlags := []cli.Flag{
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "taint.host",
			Value: "localhost",
			Usage: "Taint server host `HOST`.",
		}),
		// note: use IntFlag for numbers, UintFlag don't seem to work with the toml files
		altsrc.NewIntFlag(cli.IntFlag{
			Name:  "taint.port",
			Value: 5050,
			Usage: "Taint server port `PORT`.",
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

	app.Commands = []cli.Command{
		{
			Name:  "register",
			Usage: "People who aren't registered to SAFU yet must call this first",
			Flags: commonFlags,
			Action: func(c *cli.Context) error {
				client, err := setup(c)
				if err != nil {
					return err
				}
				res, err := client.Register()
				if err != nil {
					return err
				}
				jsonOut, _ := json.Marshal(res)
				fmt.Printf("%s\n", jsonOut)
				return nil
			},
		},
		{
			Name:      "reset_rep",
			Usage:     "can only be called by admin; resets all +rep",
			Flags:     commonFlags,
			ArgsUsage: "<target_address>",
			Action: func(c *cli.Context) error {
				client, err := setup(c)
				if err != nil {
					return err
				}
				targetAddress := c.Args().Get(0)
				res, err := client.ResetRep(targetAddress)
				if err != nil {
					return err
				}
				jsonOut, _ := json.Marshal(res)
				fmt.Printf("%s\n", jsonOut)
				return nil
			},
		},
		{
			Name:      "plus_rep",
			Usage:     "Must be VIP member and can only do it once a day",
			Flags:     commonFlags,
			ArgsUsage: "<target_address> <scam_report_id>",
			Action: func(c *cli.Context) error {
				client, err := setup(c)
				if err != nil {
					return err
				}
				address := c.Args().Get(0)
				scamReportID := c.Args().Get(1)
				res, err := client.PlusRep(address, scamReportID)
				if err != nil {
					return err
				}
				jsonOut, _ := json.Marshal(res)
				fmt.Printf("%s\n", jsonOut)
				return nil
			},
		},
		{
			Name:      "neg_rep",
			Usage:     "Must be VIP member and can only do it once a day",
			Flags:     commonFlags,
			ArgsUsage: "<target_address> <scam_report_id>",
			Action: func(c *cli.Context) error {
				client, err := setup(c)
				if err != nil {
					return err
				}
				address := c.Args().Get(0)
				scamReportID := c.Args().Get(1)
				res, err := client.NegRep(address, scamReportID)
				if err != nil {
					return err
				}
				jsonOut, _ := json.Marshal(res)
				fmt.Printf("%s\n", jsonOut)
				return nil
			},
		},
		{
			Name:      "upgrade",
			Usage:     "Normal member upgrade to VIP member by paying 500 PERL and must have 20 rep",
			Flags:     commonFlags,
			ArgsUsage: "<target_address> <scam_report_id>",
			Action: func(c *cli.Context) error {
				client, err := setup(c)
				if err != nil {
					return err
				}
				res, err := client.Upgrade()
				if err != nil {
					return err
				}
				jsonOut, _ := json.Marshal(res)
				fmt.Printf("%s\n", jsonOut)
				return nil
			},
		},
		{
			Name:      "deposit",
			Usage:     "add to a balanced stored for the wallet",
			Flags:     commonFlags,
			ArgsUsage: "<deposit_amount>",
			Action: func(c *cli.Context) error {
				client, err := setup(c)
				if err != nil {
					return err
				}
				depositAmount, err := strconv.Atoi(c.Args().Get(0))
				if err != nil {
					return err
				}
				res, err := client.Deposit(depositAmount)
				if err != nil {
					return err
				}
				jsonOut, _ := json.Marshal(res)
				fmt.Printf("%s\n", jsonOut)
				return nil
			},
		},
		{
			Name:      "withdraw",
			Usage:     "withdraw a balance stored for the wallet",
			Flags:     commonFlags,
			ArgsUsage: "<withdraw_amount>",
			Action: func(c *cli.Context) error {
				client, err := setup(c)
				if err != nil {
					return err
				}
				withdrawAmount, err := strconv.Atoi(c.Args().Get(0))
				if err != nil {
					return err
				}
				res, err := client.Withdraw(withdrawAmount)
				if err != nil {
					return err
				}
				jsonOut, _ := json.Marshal(res)
				fmt.Printf("%s\n", jsonOut)
				return nil
			},
		},
		{
			Name:  "balance",
			Usage: "check the balance of the account",
			Flags: commonFlags,
			Action: func(c *cli.Context) error {
				client, err := setup(c)
				if err != nil {
					return err
				}
				res, err := client.Balance()
				if err != nil {
					return err
				}
				jsonOut, _ := json.Marshal(res)
				fmt.Printf("%s\n", jsonOut)
				return nil
			},
		},
		{
			Name:      "query",
			Usage:     "query an address for the taint value",
			Flags:     commonFlags,
			ArgsUsage: "<address>",
			Action: func(c *cli.Context) error {
				client, err := setup(c)
				if err != nil {
					return err
				}
				address := c.Args().Get(0)
				res, err := client.Query(address)
				if err != nil {
					return err
				}
				jsonOut, _ := json.Marshal(res)
				fmt.Printf("%s\n", jsonOut)
				return nil
			},
		},
		{
			Name:      "register_scam_report",
			Usage:     "make a scam report",
			Flags:     commonFlags,
			ArgsUsage: "<scam_report_id>",
			Action: func(c *cli.Context) error {
				client, err := setup(c)
				if err != nil {
					return err
				}
				scamReportID := c.Args().Get(0)
				res, err := client.RegisterScamReport(scamReportID)
				if err != nil {
					return err
				}
				jsonOut, _ := json.Marshal(res)
				fmt.Printf("%s\n", jsonOut)
				return nil
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse configuration/command-line arugments.")
	}
}

func setup(c *cli.Context) (*client.Client, error) {
	privateKeyFile := c.String("private_key_file")

	var privateKeyHex string
	if len(privateKeyFile) > 0 && privateKeyFile != "random" {
		bytes, err := ioutil.ReadFile(privateKeyFile)
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to open server private key file: %s", privateKeyFile)
		}
		privateKeyHex = strings.TrimSpace(string(bytes))
	} else {
		return nil, errors.Errorf("Invalid private key in file: %s", privateKeyFile)
	}

	client, err := client.NewClient(client.Config{
		PrivateKeyHex: privateKeyHex,
		Host:          c.String("taint.host"),
		Port:          c.Uint("taint.port"),
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}
