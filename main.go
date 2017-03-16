package main

import (
	"fmt"

	"github.com/slowteetoe/tidechecker/tides"
)

func main() {

	holder := tides.ObservationHolder{Locations: make(map[string]*tides.Location)}

	err := holder.LoadDataStore("data")
	if err != nil {
		fmt.Printf("Failed to load data: %v\n", err)
		return
	}

	for index := 0; index < 5; index++ {
		fmt.Printf("%v\n", holder.Locations["9410230"].Items[index])
	}

}
