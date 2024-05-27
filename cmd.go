package main

import (
	"github.com/globalmaksimum/vectbench/pkg/cli"
	_ "github.com/go-sql-driver/mysql"
)

var (
	tag       string //last git tag
	sha1ver   string // sha1 revision used to build the program
	buildTime string // when the executable was built
	hostname  string // build host
	goV       string //go compiler version used for build
)

func main() {

	cli.Initialize(tag, sha1ver, buildTime, hostname, goV)
	cli.Execute()

}
