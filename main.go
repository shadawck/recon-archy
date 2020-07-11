package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/remiflavien1/recon-archy/selenium"
	"github.com/urfave/cli/v2"
)

func main() {

	// FLAG TARGET
	var threads string
	var company string

	app := &cli.App{
		Name:  "ReconArchy",
		Usage: "Crawl 1000 employees of a choosen company and build their organizational chart",

		Commands: []*cli.Command{
			{
				Name:  "crawl",
				Usage: "crawl employees specific to a company",

				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "threads",
						Value:       "1",
						Aliases:     []string{"t"},
						Usage:       "Adjust number of crawling worker",
						Destination: &threads,
					},
					&cli.StringFlag{
						Name:        "company",
						Aliases:     []string{"c"},
						Usage:       "Name of the target company",
						Destination: &company,
					},
				}, // Crawl Command Flags

				Action: func(c *cli.Context) error {
					// ARGUMENT HANDLING (in one Action - for now there is only one commands and no subcommands)
					if company == "" {
						fmt.Println("Please choose a company to target")
						os.Exit(0)
					}
					threadsInt, errConv := strconv.Atoi(threads)
					if errConv != nil {
						panic(errConv)
					}
					if threadsInt > 4 {
						fmt.Printf("\nBy setting THREADS to %d, you will spawn %d WebDriver on your machine.", threadsInt, threadsInt)
						fmt.Printf("\nI think she can not handle it.\nSetup THREAD to 4 at max the next time.")
						threadsInt = 4
						fmt.Printf("\nAuto - Setting THREADS to 4\n")
					}

					// MAIN - START PROGRAM
					selenium.Start(company, threadsInt)

					return nil
				}, // Action
			}, // Command command
		}, // Command main
	} // App

	// Run apps cli
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
