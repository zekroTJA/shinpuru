package main

import (
	"flag"
	"fmt"

	"github.com/zekroTJA/shinpuru/internal/services/config"
)

var c = flag.String("c", "", "")

func main() {
	flag.Parse()

	p := config.NewPaerser(flag.Args(), *c)

	fmt.Println(p.Parse())
}
