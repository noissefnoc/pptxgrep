package main

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"html"
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

const (
	FNAME             = "pptxgrep"
	VERSION           = "0.0.2"
	SLIDE_PATH_PREFIX = "ppt/slides/slide"
)

func usage() {
	fmt.Printf("Usage:\n  %s [options] pattern pptx1 [pptx2 ... pptxN]\n\n", FNAME)
	fmt.Printf("Version:\n  %s\n\n", VERSION)
	fmt.Printf("Options:\n")
	flag.PrintDefaults()
}

func extractLocation(filePath, prefix string) string {
	return strings.TrimRight(strings.TrimLeft(filePath, prefix), ".xml")
}

func colorize(targetStr string, pattern *regexp.Regexp) string {
	var buf bytes.Buffer

	matchedIdx := pattern.FindAllStringSubmatchIndex(targetStr, -1)
	startIdx := 0

	for _, startEnd := range matchedIdx {
		matchedStart := startEnd[0]
		matchedEnd := startEnd[1]

		fmt.Fprintf(
			&buf,
			"%s\x1b[31m%s\x1b[0m",
			string([]byte(targetStr)[startIdx:matchedStart]),
			string([]byte(targetStr)[matchedStart:matchedEnd]))

		startIdx = matchedEnd
	}

	if len(targetStr) > startIdx {
		fmt.Fprint(&buf, string([]byte(targetStr)[startIdx:]))
	}

	return buf.String()
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

func pptxgrep(pattern *regexp.Regexp, arg string, color bool) error {
	r, err := zip.OpenReader(arg)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		var node Node
		var buf bytes.Buffer

		if strings.HasPrefix(f.Name, SLIDE_PATH_PREFIX) {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			b, err := ioutil.ReadAll(rc)
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

			unescapedString := html.UnescapeString(buf.String())

			if pattern.MatchString(unescapedString) {
				if color {
					unescapedString = colorize(unescapedString, pattern)
				}
				fmt.Printf("%s:%s:%s\n", arg, extractLocation(f.Name, SLIDE_PATH_PREFIX), unescapedString)
			}
		}
	}
	return nil
}

func main() {
	var version bool
	var color bool
	var ignoreCase bool
	flag.BoolVar(&version, "version", false, "print version")
	flag.BoolVar(&color, "color", false, "colorize matched pattern")
	flag.BoolVar(&ignoreCase, "i", false, "ignore case when matching pattern")
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() <= 1 {
		if version {
			fmt.Printf("Version: %s\n", VERSION)
		} else {
			flag.Usage()
		}
		os.Exit(0)
	}

	var pattern *regexp.Regexp

	for i, arg := range flag.Args() {
		if i == 0 {
			if ignoreCase {
				arg = `(?i)` + arg
			}
			pattern = regexp.MustCompile(arg)
		} else {
			if err := pptxgrep(pattern, arg, color); err != nil {
				log.Fatal(err)
			}
		}
	}
}
