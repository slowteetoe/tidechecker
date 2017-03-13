package main

import (
	"encoding/xml"
	"fmt"
	"os"

	"golang.org/x/net/html/charset"
)

type TideLevel int

const (
	LOW TideLevel = iota
	HIGH
)

type ObservationHolder struct {
	XMLName     xml.Name       `xml:"datainfo"`
	Items       []*Observation `xml:"data>item"`
	State       string         `xml:"state"`
	StationType string         `xml:"stationtype"`
}

type Observation struct {
	Level     string
	HighOrLow string  `xml:"highlow"`
	Date      string  `xml:"date"`
	Feet      float32 `xml:"predictions_in_ft"`
	Time      string  `xml:"time"`
	Direction string
}

func main() {
	xmlFile, err := os.Open("data/9410230_annual.xml")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()

	var holder ObservationHolder

	// Since our XML is encoded in ISO-8859-1, we need to decode it
	decoder := xml.NewDecoder(xmlFile)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&holder)

	if err != nil {
		fmt.Printf("Unable to unmarshal correctly: %v\n", err)
		return
	}

	// lag over observations to see if rising or falling
	for i := range holder.Items {
		if i == 0 {
			holder.Items[i].HighOrLow = "???"
		} else if holder.Items[i-1].HighOrLow == "H" {
			holder.Items[i].Direction = "FALLING"
		} else {
			holder.Items[i].Direction = "RISING"
		}
	}

	for index := 0; index < 5; index++ {
		fmt.Printf("%v\n", holder.Items[index])
	}

}
