package main

import (
	"flag"
	"fmt"
	"github.com/mattkimber/temporarily-late/internal/manifest"
	"github.com/mattkimber/temporarily-late/internal/template"
	"os"
)

var isRoad bool
var aniFlag bool

func init() {
	flag.BoolVar(&isRoad, "road", false, "Produce simpler templates for road vehicles")
	flag.BoolVar(&aniFlag, "anim", false, "Add the ANIM flag to output")
}

func main() {
	flag.Parse()
	for _, filename := range flag.Args() {
		processFile(filename)
	}
}

func processFile(filename string) {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		panic(fmt.Errorf("cannot open file %s: %v", filename, err))
	}

	m, err := manifest.FromJson(f)
	if err != nil {
		panic(fmt.Errorf("cannot parse file %s: %v", filename, err))
	}

	if isRoad {
		template.WriteRoadTemplates(m, aniFlag)
	} else {
		template.WriteTemplates(m, aniFlag)
	}
}
