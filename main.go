package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

func pptxgrep(pattern *regexp.Regexp, arg string) error {
	r, err := zip.OpenReader(arg)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if strings.HasPrefix(f.Name, "ppt/slides/slide") {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			b, _ := ioutil.ReadAll(rc)
			if err != nil {
				return err
			}

			if pattern.Match(b) {
				fmt.Printf("%s\n", f.Name)
			}
		}
	}

	return nil
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	var pattern *regexp.Regexp

	for i, arg := range flag.Args() {
		if i == 0 {
			pattern = regexp.MustCompile(arg)
		} else {
			if err := pptxgrep(pattern, arg); err != nil {
				log.Fatal(err)
			}
		}
	}
}
