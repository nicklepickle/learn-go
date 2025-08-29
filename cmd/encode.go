package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"errors"
)

type Color struct {
	Name string `xml:",chardata"` // inner text
	Hash string `xml:"hash,attr"` // hash attribute
}

func (c *Color) ToArray() []string {
	return []string{c.Name, c.Hash}
}

// need a root element for xml encoding
type Root struct {
	Colors []Color
}

var Colors = []Color{
	{"Red", "#FF0000"},
	{"White", "#FFFFFF"},
	{"Cyan", "#00FFFF"},
	{"Silver", "#C0C0C0"},
	{"Blue", "#0000FF"},
	{"Grey", "#808080"},
	{"DarkBlue", "#00008B"},
	{"Black", "#000000"},
	{"LightBlue", "#ADD8E6"},
	{"Orange", "#FFA500"},
	{"Purple", "#800080"},
	{"Brown", "#A52A2A"},
	{"Yellow", "#FFFF00"},
	{"Maroon", "#800000"},
	{"Lime", "#00FF00"},
	{"Green", "#008000"},
	{"Magenta", "#FF00FF"},
	{"Olive", "#808000"},
	{"Pink", "#FFC0CB"},
	{"Aquamarine", "#7FFFD4"},
}

func GetColors(encoding string) (string, error) {
	switch encoding {
	case "json":
		bytes, err := json.Marshal(Colors)
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	case "xml":
		root := Root{}
		root.Colors = Colors
		bytes, err := xml.Marshal(root)
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	case "csv":
		buf := new(bytes.Buffer)
		writer := csv.NewWriter(buf)
		// header
		writer.Write([]string{"Name", "Hash"})

		for _, c := range Colors {
			if err := writer.Write(c.ToArray()); err != nil {
				return "", err
			}
		}

		writer.Flush()

		return buf.String(), nil
	}
	return "", errors.New("encoding not recognized")
}
