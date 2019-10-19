package main

import (
	"archive/zip"
	"encoding/xml"
	"flag"
	"io/ioutil"
	"log"
	"os"
)

type Relationship struct {
	Text       string `xml:",chardata"`
	ID         string `xml:"Id,attr"`
	Type       string `xml:"Type,attr"`
	Target     string `xml:"Target,attr"`
	TargetMode string `xml:"TargetMode,attr"`
}

type Relationships struct {
	XMLName      xml.Name       `xml:"Relationships"`
	Text         string         `xml:",chardata"`
	Xmlns        string         `xml:"xmlns,attr"`
	Relationship []Relationship `xml:"Relationship"`
}

type Presentation struct {
	XMLName xml.Name `xml:"presentation"`
	Text    string   `xml:",chardata"`
}

type Node struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:"-"`
	Content []byte     `xml",innerxml"`
	Nodes   []Node     `xml:",any"`
}

func (n *Node) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	n.Attrs = start.Attr
	type node Node

	return d.DecodeElement((*node)(n), &start)
}

func pptxgrep(arg string) error {
	r, err := zip.OpenReader(arg)
	if err != nil {
		return err
	}
	defer r.Close()

	var rels Relationships

	for _, f := range r.File {
		switch f.Name {
		case "ppt/_rels/presentation.xml.rels":
			rc, err := f.Open()
			defer rc.Close()

			b, _ := ioutil.ReadAll(rc)
			if err != nil {
				return err
			}

			err = xml.Unmarshal(b, &rels)
			if err != nil {
				return err
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
	for _, arg := range flag.Args() {
		if err := pptxgrep(arg); err != nil {
			log.Fatal(err)
		}
	}
}
