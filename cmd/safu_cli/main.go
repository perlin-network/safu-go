package main

import (
	"encoding/json"
	"fmt"
	"github.com/perlin-network/safu-go/cmd/safu_cli/client"
	"github.com/perlin-network/safu-go/log"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"
	"os"
	"sort"
	"strconv"
	"time"
)

func main() {
	app := cli.NewApp()

	app.Name = "safu_cli"
	app.Author = "Perlin Network"
	app.Version = "0.0.1"
	app.Usage = "A cli that submits taint requests"

	taintFlags := []cli.Flag{
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
			Name:  "account_id",
			Value: "account_id",
			Usage: "The public key of the wallet that originally registered.",
		}),
	}

	ledgerFlags := []cli.Flag{
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "wavelet.host",
			Value: "localhost",
			Usage: "Wavelet api server host `HOST`.",
		}),
		// note: use IntFlag for numbers, UintFlag don't seem to work with the toml files
		altsrc.NewIntFlag(cli.IntFlag{
			Name:  "wavelet.port",
			Value: 3000,
			Usage: "Wavelet api server port `PORT`.",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "wctl.path",
			Value: "./wctl",
			Usage: "Path to wctl binary",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "private_key_file",
			Value: "wallet.txt",
			Usage: "TXT file that contain's the wallet's private key `PRIVATE_KEY_FILE`.",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "contract.address",
			Value: "TODO",
			Usage: "Address of the smart contract",
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
			Flags: ledgerFlags,
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
			Flags:     ledgerFlags,
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
			Flags:     ledgerFlags,
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
			Flags:     ledgerFlags,
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
			Flags:     ledgerFlags,
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
			Usage:     "Add to the balance of an account",
			Flags:     ledgerFlags,
			ArgsUsage: "<deposit_amount>",
			Action: func(c *cli.Context) error {
				client, err := setup(c)
				if err != nil {
					return err
				}
				depositAmount, err := strconv.ParseInt(c.Args().Get(0), 10, 64)
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
			Name:      "query",
			Usage:     "query an address for the taint value",
			Flags:     taintFlags,
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
			Flags:     append(ledgerFlags, taintFlags...),
			ArgsUsage: "<scammer_address> <victim_address> <title> <content>",
			Action: func(c *cli.Context) error {
				client, err := setup(c)
				if err != nil {
					return err
				}
				scammerAddress := c.Args().Get(0)
				victimAddress := c.Args().Get(1)
				title := c.Args().Get(2)
				content := c.Args().Get(3)
				res, err := client.RegisterScamReport(scammerAddress, victimAddress, title, content)
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
	client, err := client.NewClient(&client.Config{
		PrivateKeyFile:       c.String("private_key_file"),
		TaintHost:            c.String("taint.host"),
		TaintPort:            c.Uint("taint.port"),
		WCTLPath:             c.String("wctl.path"),
		WaveletHost:          c.String("wavelet.host"),
		WaveletPort:          c.Uint("wavelet.port"),
		AccountID:            c.String("account_id"),
		SmartContractAddress: c.String("contract.address"),
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}
