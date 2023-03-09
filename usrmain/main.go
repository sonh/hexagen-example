package main

import (
	"fmt"
	"golang.org/x/tools/go/packages"
	"log"
	"strings"
)

func main() {

	var err error
	fmt.Println(fmt.Sprintf("%v: abc", err))

	parsePackage([]string{"/home/sonhuynh/projects/x/generator/usr"}, nil)
}

type String string

func parsePackage(patterns []string, tags []string) {
	cfg := &packages.Config{
		Mode:       packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax,
		Tests:      false,
		BuildFlags: []string{fmt.Sprintf("-tags=%s", strings.Join(tags, " "))},
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("error: %d packages found", len(pkgs))
	}

	for _, pkg := range pkgs {
		fmt.Printf("%+v \n", pkg.Name)
	}
}
