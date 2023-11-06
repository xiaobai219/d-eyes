package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/urfave/cli/v2"
	"github.com/xuri/excelize/v2"

	"d-eyes/basicinfo/info"
	"d-eyes/configcheck/check"
	"d-eyes/filedetection"
	"d-eyes/logo"
	"d-eyes/output"
	"d-eyes/process/controller"
	"d-eyes/yaraobj"
)

var path string
var rule string
var thread int
var pid int

func main() {
	logo.ShowLogo()
	app := &cli.App{
		Name:  "D-Eyes",
		Usage: "The Eyes of Darkness from Nsfocus spy on everything.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "path",
				Aliases:     []string{"P"},
				Value:       "/",
				Usage:       "--path / or -P / (Only For Filescan)",
				Destination: &path,
			},
			&cli.IntFlag{
				Name:        "pid",
				Aliases:     []string{"p"},
				Value:       -1,
				Usage:       "--pid 666 or -p 666 (Only For processcan.'-1' means all processes.)",
				Destination: &pid,
			},
			&cli.StringFlag{
				Name:    "rule",
				Aliases: []string{"r"},
				//Value:   5,
				Usage:       "--rule Ransom.Wannacrypt or -r Ransom.Wannacrypt",
				Destination: &rule,
			},
			&cli.IntFlag{
				Name:        "thread",
				Aliases:     []string{"t"},
				Value:       5,
				Usage:       "--thread 1 or -t 1 (Only For Filescan)",
				Destination: &thread,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "filescan",
				Aliases: []string{"fs"},
				Usage:   "Command for scanning filesystem",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "path",
						Aliases:     []string{"P"},
						Value:       "/",
						Usage:       "--path / or -P / (Only For Filescan)",
						Destination: &path,
					},
					&cli.StringFlag{
						Name:    "rule",
						Aliases: []string{"r"},
						//Value:   5,
						Usage:       "--rule Ransom.Wannacrypt or -r Ransom.Wannacrypt",
						Destination: &rule,
					},
					&cli.IntFlag{
						Name:        "thread",
						Aliases:     []string{"t"},
						Value:       5,
						Usage:       "--thread 1 or -t 1",
						Destination: &thread,
					},
				},
				Action: func(c *cli.Context) error {
					// fmt.Println("added task: ", c.Args().First())
					//
					//
					var paths []string
					r := []output.Result{}
					paths = strings.Split(path, ",")
					var start = time.Now()
					var sum = 0

					if rule == "" {
						yaraRule := "./yaraRules"
						rules, err := yaraobj.LoadAllYaraRules(yaraRule)
						if err != nil {
							color.Redln("LoadCompiledRules goes error !!!")
							color.Redln("GetRules err: ", err)
							os.Exit(1)
						}
						for _, path := range paths {
							files := filedetection.StartFileScan(path, rules, thread, &r)
							sum += files
						}
					} else {
						yaraRule := "./yaraRules/" + rule + ".yar"
						_, err := os.Lstat(yaraRule)
						if err != nil {
							color.Redln("There is no such rule yet !!!")
							os.Exit(1)
						}
						rules, err := yaraobj.LoadSingleYaraRule(yaraRule)
						if err != nil {
							color.Redln("GetRules err: ", err)
							os.Exit(1)
						}
						for _, path := range paths {
							files := filedetection.StartFileScan(path, rules, thread, &r)
							sum += files
						}
					}

					if len(r) > 0 {
						length := len(r)
						categories := map[string]string{
							"A1": "Risk Description", "B1": "Risk File Path",
						}
						var values = make(map[string]string)
						vulsumTmp := 0
						for i := 0; i < length; i++ {
							vulsumTmp++
							color.Error.Println("[ Risk ", vulsumTmp, " ]")
							fmt.Print("Risk Description: ")
							color.Warn.Println(r[i].Risk)
							fmt.Println("Risk File Path: ")
							color.Warn.Println(r[i].RiskPath)
							//set excel values
							excelValuetmpA := "A" + strconv.Itoa(vulsumTmp+1)
							excelValuetmpB := "B" + strconv.Itoa(vulsumTmp+1)
							values[excelValuetmpA] = r[i].Risk
							values[excelValuetmpB] = r[i].RiskPath
						}
						//output to a excel
						f := excelize.NewFile()
						f.SetColWidth("Sheet1", "A", "B", 50)
						for k, v := range categories {
							f.SetCellValue("Sheet1", k, v)
						}
						for k, v := range values {
							f.SetCellValue("Sheet1", k, v)
						}
						style, err := f.NewStyle(
							&excelize.Style{
								Font: &excelize.Font{
									Bold:  true,
									Size:  11,
									Color: "e83723",
								},
							},
						)
						if err != nil {
							fmt.Println(err)
						}
						f.SetCellStyle("Sheet1", "A1", "A1", style)
						f.SetCellStyle("Sheet1", "B1", "B1", style)
						// save the result to Deyes.xlsx
						if err := f.SaveAs("D-Eyes.xlsx"); err != nil {
							fmt.Println(err)
						}
					} else {
						fmt.Println("\nNo suspicious files found. Your computer is safe with the rules you choose.")
					}
					var end = time.Now().Sub(start)
					fmt.Println("Consuming Time: ", end, "  Number of scanned files: ", sum)

					return nil
				},
			},
			{
				Name:    "processcan",
				Aliases: []string{"ps"},
				Usage:   "Command for scanning processes",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:        "pid",
						Aliases:     []string{"p"},
						Value:       -1,
						Usage:       "--pid 666 or -p 666 ('-1' means all processes.)",
						Destination: &pid,
					},
					&cli.StringFlag{
						Name:        "rule",
						Aliases:     []string{"r"},
						Usage:       "--rule Ransom.Wannacrypt or -r Ransom.Wannacrypt",
						Destination: &rule,
					},
				},
				Action: func(c *cli.Context) error {
					var start = time.Now()
					controller.ScanProcess(pid, rule)
					var end = time.Now().Sub(start)
					fmt.Println("Consuming Time: ", end)
					return nil
				},
			},
			{
				Name:    "selfcheck",
				Aliases: []string{"sc"},
				Usage:   "Command for checking some files which may have backdoors",
				Action: func(c *cli.Context) error {
					check.Trigger()
					return nil
				},
			},
			{
				Name:  "host",
				Usage: "Command for displaying basic host information",
				Action: func(c *cli.Context) error {
					color.Infoln("Host Info:")
					info.DisplayBaseInfo()
					return nil
				},
			},
			{
				Name:  "users",
				Usage: "Command for displaying all the users on the host",
				Action: func(c *cli.Context) error {
					color.Infoln("AllUsers:")
					info.DisplayAllUsers()
					return nil
				},
			},
			{
				Name:  "top",
				Usage: "Command for displaying the top 15 processes in CPU usage",
				Action: func(c *cli.Context) error {
					info.Top()
					return nil
				},
			},
			{
				Name:  "netstat",
				Usage: "Command for displaying host network information",
				Action: func(c *cli.Context) error {
					color.Infoln("Network Info:")
					info.DisplayNetStat()
					return nil
				},
			},
			{
				Name:  "task",
				Usage: "Command for displaying all the tasks on the host",
				Action: func(c *cli.Context) error {
					color.Infoln("Task:")
					info.DisplayPlanTask()
					return nil
				},
			},
			{
				Name:  "autoruns",
				Usage: "Command for displaying all the autoruns on the host",
				Action: func(c *cli.Context) error {
					color.Infoln("Autoruns:")
					info.CallDisplayAutoruns()
					return nil
				},
			},
			{
				Name:  "export",
				Usage: "Command for exporting basic host information",
				Action: func(c *cli.Context) error {
					info.SaveSummaryBaseInfo()
					return nil
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
