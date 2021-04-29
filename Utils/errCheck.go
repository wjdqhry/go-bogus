package Utils

import (
	"github.com/fatih/color"
	"os"
)

func ErrCheck(err error) {
	if err != nil {
		c := color.New(color.BgRed).Add(color.Bold)
		c.Println(err.Error())
		os.Exit(0)
	}
}
