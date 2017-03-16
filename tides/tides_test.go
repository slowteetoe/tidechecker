package tides_test

import (
	"testing"

	"github.com/slowteetoe/tidechecker/tides"
)

// TODO create a smaller, testing-only file

func TestLoadingNonexistentFile(t *testing.T) {
	holder := tides.ObservationHolder{}
	err := holder.LoadDataStore("data_does_not_exist")
	if err == nil {
		t.Error("Should have errored on dir that does not exist")
	}
}

func TestLoadingGoodFile(t *testing.T) {
	holder := tides.ObservationHolder{Locations: make(map[string]*tides.Location)}
	err := holder.LoadDataStore("../data")
	if err != nil {
		t.Errorf("Should not have errored: %v\n", err)
	}
	expected := 1410 // yes, it's magic
	location := holder.Locations["9410230"]
	if location == nil {
		t.Error("Failed to find location in holder")
	}
	actual := len(location.Items)
	if actual != expected {
		t.Errorf("Incorrect number of items, expected=%d, actual=%d", expected, actual)
	}
	expectedState := "CA"
	if state := location.State; state != expectedState {
		t.Errorf("Expected to be tides from %s, but was: %s", expectedState, state)
	}
	expectedStationType := "Harmonic"
	if location.StationType != expectedStationType {
		t.Errorf("Expected station type %s but was: %s", expectedStationType, location.StationType)
	}
}
