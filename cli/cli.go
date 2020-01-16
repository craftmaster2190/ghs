package cli

import (
	"github.com/urfave/cli"
	"os"
)

func Start() {
	app := cli.NewApp()
	app.Name = "Github Review Statistics"
	app.Usage = "Displays statistics for reviews KPIs"
	app.Version = "0.0.1"
	app.Authors = []cli.Author{{Name: "Ondrej Burkert"}}

	flags := []cli.Flag{
		cli.StringFlag{
			Name:  "org",
			Value: "",
			Usage: "Organization",
		},
		cli.StringFlag{
			Name:  "team",
			Value: "",
			Usage: "Team name",
		},
		cli.StringSliceFlag{
			Name:  "repo",
			Value: nil,
			Usage: "Repository names (May be repeated)",
		},
		cli.IntFlag{
			Name:  "start, s",
			Value: -14,
			Usage: "Start of range - days from now. E. g. -14",
		},
		cli.IntFlag{
			Name:  "end, e",
			Value: -1,
			Usage: "End of range - days from now. E. g. -7",
		},
		cli.StringFlag{
			Name:  "api-token, t",
			Value: "",
			Usage: "Github API token with repo and org scope",
		},
	}

	singleFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "org",
			Value: "",
			Usage: "Organization",
		},
		cli.StringFlag{
			Name:  "repo",
			Value: "",
			Usage: "Repository",
		},
		cli.StringFlag{
			Name:  "api-token, t",
			Value: "",
			Usage: "Github API token with repo and org scope",
		},
		cli.IntFlag{
			Name:  "pr-number, p",
			Value: 0,
			Usage: "PR Number",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "stats",
			Usage:  "Get Statistics",
			Action: stats,
			Flags:  flags,
		},
		{
			Name:   "single",
			Usage:  "Get Single PR Details",
			Action: singlePr,
			Flags:  singleFlags,
		},
	}

	app.Run(os.Args)
}
