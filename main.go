package main

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

type Node struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:"-"`
	Content []byte     `xml:",innerxml"`
	Nodes   []Node     `xml:",any"`
}

func walk(node *Node, w io.Writer) error {
	switch node.XMLName.Local {
	case "t":
		fmt.Fprintf(w, string(node.Content))
	default:
		for _, n := range node.Nodes {
			if err := walk(&n, w); err != nil {
				return err
			}
		}
	}

	return nil
}

func pptxgrep(pattern *regexp.Regexp, arg string) error {
	r, err := zip.OpenReader(arg)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		var node Node
		var buf bytes.Buffer

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

			err = xml.Unmarshal(b, &node)
			if err != nil {
				return err
			}

			err = walk(&node, &buf)
			if err != nil {
				return err
			}

			if pattern.MatchString(buf.String()) {
				fmt.Printf("%s: %s\n\n", f.Name, buf.String())
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
