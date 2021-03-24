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
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	. "github.com/dave/jennifer/jen"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(dir)

		//forderSlice := strings.Split(dir, "/")
		forderName := strings.SplitAfter(dir, "src/")[1]
		//forderName := forderSlice[len(forderSlice)-1]

		err = os.Mkdir("routes", os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		err = os.Chdir(dir + "/routes")
		if err != nil {
			log.Fatal(err)
		}

		routesFile := NewFilePathName(dir+"/routes", "routes")
		routesFile.Func().Id("RegisterApi").Params().Block()
		routesFile.Save("routes.go")

		err = os.Chdir(dir)
		if err != nil {
			log.Fatal(err)
		}

		f := NewFilePathName(dir, "main")
		// app := fiber.New()
		// app.Use(cors.New(), logger.New(), recover.New())
		// routes.RegisterApi(app)
		// app.Listen(":12270")
		fmt.Println("creating..")
		f.ImportAlias("github.com/gofiber/fiber/v2", "fiber")
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
			Qual(forderName+"/routes", "RegisterApi").Call(),
			Id("app").Dot("Listen").Call(Lit(":8000")),
		)

		mainName := "main.go"
		fmt.Println("saving..")
		err = f.Save(mainName)
		fmt.Println("saved..")
		if err != nil {
			fmt.Println(err)
		}
		modCommand := exec.Command("go", "mod", "init")
		err = modCommand.Run()

		if err != nil {
			fmt.Println(err)
		}

		buildCommand := exec.Command("go", "build", mainName)

		err = buildCommand.Run()
		// if err != nil {
		// 	log.Println("err: ", err)
		// }

		fmt.Println("Enjoy Bogus-CLI 🤪")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
