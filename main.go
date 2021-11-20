package main

import (
	_ "github.com/lib/pq"
	"github.com/themisir/databoard/cli"
)

func main() {
	app := cli.NewApp()
	app.Run()
}
