package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/orcastor/iwork-converter/iwork2html"
	"github.com/orcastor/iwork-converter/iwork2text"
)

func main() {
	verbose := flag.Bool("v", false, "Enable verbose debug output")
	flag.Parse()

	if len(flag.Args()) < 2 {
		fmt.Printf(`Converts iWork files to html/json/txt

Usage:
    %s [-v] infile.key outfile.html
    %s [-v] infile.pages outfile.html
    %s [-v] infile.numbers outfile.html
    %s [-v] infile.key outfile.txt

Options:
    -v    Enable verbose debug output

`, os.Args[0], os.Args[0], os.Args[0], os.Args[0])
		return
	}

	infile := flag.Args()[0]
	outfile := flag.Args()[1]

	// Set debug mode globally
	iwork2html.SetDebugMode(*verbose)
	iwork2text.SetDebugMode(*verbose)

	switch {
	case strings.HasSuffix(outfile, ".txt"):
		if err := iwork2text.Convert(infile, outfile); err != nil {
			panic(err)
		}
	default:
		if err := iwork2html.Convert(infile, outfile); err != nil {
			panic(err)
		}
	}
}
