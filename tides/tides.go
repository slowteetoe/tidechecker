package tides

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	"sort"

	"golang.org/x/net/html/charset"
)

type tideLevel int

const (
	low tideLevel = iota
	high
)

// ObservationHolder maps a station id to the location
type ObservationHolder struct {
	Locations map[string]*Location
}

// Location holds station information and the tide predictions
type Location struct {
	XMLName     xml.Name      `xml:"datainfo"`
	Items       []*Prediction `xml:"data>item"`
	State       string        `xml:"state"`
	StationType string        `xml:"stationtype"`
	StationID   string        `xml:"stationid"`
}

func (loc Location) String() string {
	return fmt.Sprintf("Location: id=%s, state=%s, type=%s", loc.StationID, loc.State, loc.StationType)
}

// Prediction is provided by the NOAA
type Prediction struct {
	Level     string
	HighOrLow string  `xml:"highlow"`
	Date      string  `xml:"date"`
	Feet      float32 `xml:"predictions_in_ft"`
	Time      string  `xml:"time"`
	Direction string
}

func (p Prediction) String() string {
	return fmt.Sprintf("date=%v, time=%v, tideHeightFt=%0.1f, highOrLow=%s, dir=%s", p.Date, p.Time, p.Feet, p.HighOrLow, p.Direction)
}

// LoadDataStore loads all prediction files in the given directory
func (holder *ObservationHolder) LoadDataStore(dir string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		loc, err := loadPredictionFile(dir + "/" + file.Name())
		if err != nil {
			return err
		}
		holder.Locations[loc.StationID] = loc
	}
	return nil
}

func loadPredictionFile(file string) (*Location, error) {
	xmlFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer xmlFile.Close()

	// Since our XML is encoded in ISO-8859-1, we need to decode it
	decoder := xml.NewDecoder(xmlFile)
	decoder.CharsetReader = charset.NewReaderLabel

	var loc Location

	err = decoder.Decode(&loc)
	if err != nil {
		fmt.Printf("Unable to unmarshal correctly: %v\n", err)
		return nil, err
	}

	// Now loop over to decide whether the tide will be rising or falling until the next observation
	for i := range loc.Items {
		if loc.Items[i].HighOrLow == "H" {
			loc.Items[i].Direction = "FALLING"
		} else {
			loc.Items[i].Direction = "RISING"
		}
	}
	return &loc, nil
}

// FindNearestPrediction finds the closest prediction to the provided date, using binary search until we move to a smarter data structure
func (loc *Location) FindNearestPrediction(date string) *Prediction {

	i := sort.Search(len(loc.Items), func(i int) bool {
		return loc.Items[i].Date >= date
	})

	if i < len(loc.Items) && loc.Items[i].Date == date {
		return loc.Items[i]
	}

	return nil
}
