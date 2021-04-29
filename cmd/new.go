/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/wjdqhry/go-bogus/Utils"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	. "github.com/dave/jennifer/jen"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "start a new project",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) <= 0 {
			color.HiRed("Please Write Your Project Name... 😅")
			os.Exit(0)
		}

		projectName := args[0]
		if projectName != "" {
			dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(dir)

			//forderSlice := strings.Split(dir, "/")
			if !strings.Contains(dir, "src") {
				log.Fatalf("Not in GOPATH")
			}

			// create project
			Utils.ErrCheck(os.Mkdir(projectName, os.ModePerm))

			// directory to project
			dir = dir + "/" + projectName

			Utils.ErrCheck(os.Chdir(dir))

			Utils.ErrCheck(os.Mkdir("routes", os.ModePerm))
			Utils.ErrCheck(os.Chdir(dir + "/routes"))

			routesFile := NewFilePathName(dir+"/routes", "routes")
			routesFile.ImportAlias("github.com/gofiber/fiber/v2", "fiber")
			routesFile.Func().Id("RegisterApi").Params(Id("api").Qual("github.com/gofiber/fiber/v2", "Router")).Block(
				Id("api").Dot("Get").Call(Lit("/ping"), Func().Params(Id("c").Id("*fiber.Ctx")).Error().Block(
					Return(Id("c").Dot("SendString").Call(Lit("pong"))),
				)),
				Line(),
			)
			Utils.ErrCheck(routesFile.Save("routes.go"))

			Utils.ErrCheck(os.Chdir(dir))

			//for import
			forderName := strings.SplitAfter(dir, "src/")[1]
			f := NewFilePathName(dir, "main")
			// app := fiber.New()
			// app.Use(cors.New(), logger.New(), recover.New())
			// routes.RegisterApi(app)
			// app.Listen(":12270")
			color.Green("creating..")
			f.ImportAlias("github.com/gofiber/fiber/v2", "fiber")
			f.ImportAlias("github.com/gofiber/fiber/v2/middleware/recover", "recover")

			//f.ImportName("github.com/gofiber/fiber/v2", "fiber")
			f.Func().Id("main").Params().Block(
				Qual("fmt", "Println").Call(Lit("Hello, world")),
				//f.PackagePrefix = "pkg",
				Id("app").Op(":=").Qual("github.com/gofiber/fiber/v2", "New").Call(),
				//f.Line(),

				Id("app").Dot("Use").Call(
					Qual("github.com/gofiber/fiber/v2/middleware/cors", "New").Call(),
					Qual("github.com/gofiber/fiber/v2/middleware/logger", "New").Call(),
					Qual("github.com/gofiber/fiber/v2/middleware/recover", "New").Call(),
				),
				Qual(forderName+"/routes", "RegisterApi").Call(Id("app")),
				Id("app").Dot("Listen").Call(Lit(":8000")),
			)

			mainName := "main.go"
			color.Green("saving......")
			Utils.ErrCheck(f.Save(mainName))

			color.Blue("saved......")

			modCommand := exec.Command("go", "mod", "init")
			Utils.ErrCheck(modCommand.Run())

			color.Green("mod init~")

			buildCommand := exec.Command("go", "build", mainName)

			Utils.ErrCheck(buildCommand.Run())

			color.Cyan("Project Created\nEnjoy Bogus-CLI 🤪")
		} else {
			color.HiRed("Please write your project name 😅")
		}

	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
