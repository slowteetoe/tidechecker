package main

import (
	"fmt"

	"github.com/slowteetoe/tidechecker/tides"
)

func main() {

	var holder tides.ObservationHolder

	err := holder.LoadDataStore()
	if err != nil {
		fmt.Printf("Failed to load data: %v\n", err)
		return
	}

	for index := 0; index < 5; index++ {
		fmt.Printf("%v\n", holder.Items[index])
	}

}
