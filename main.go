package main

import (
	"fmt"
	"github.com/ChrisTheBaron/strava-ical/structs"
	"github.com/ChrisTheBaron/strava-ical/utils"
	"github.com/ChrisTheBaron/strava-ical/web"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {

	var configfilepath string
	var config structs.Config

	app := cli.NewApp()
	app.Name = "strava-ical"
	app.HideVersion = true
	app.Usage = "Allows subscription to your Strava history using ICAL."
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Chris Taylor",
			Email: "christhebaron@gmail.com",
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config-file, config, c",
			Usage:       "Load configuration from `FILE` (required)",
			Destination: &configfilepath,
		},
	}

	// Before the application runs, let's just do some validation
	app.Before = func(c *cli.Context) error {
		if "" == configfilepath {
			return cli.NewExitError("Config file is required", 1)
		}
		var err error
		//Get the config from the config.yaml file
		config, err = utils.GetConfigFromFile(configfilepath)

		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		return nil
	}

	// Now we have passed validation we can get on with it
	app.Action = func(c *cli.Context) error {

		log.SetFlags(log.Llongfile)

		s, err := web.NewServer(&config)

		if err != nil {
			return cli.NewExitError(err.Error(), 2)
		}

		l := fmt.Sprintf("%s:%d", config.Server.Address, config.Server.Port)

		log.Printf("Listening on: %s", l)

		s.Run(l)

		return nil

	}

	app.Run(os.Args)

}
