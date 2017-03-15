package tides_test

import (
	"testing"

	"github.com/slowteetoe/tidechecker/tides"
)

// TODO create a smaller, testing-only file

func TestLoadingNonexistentFile(t *testing.T) {
	holder := tides.ObservationHolder{}
	err := holder.LoadDataStore("data/9410230_annual.xml")
	if err == nil {
		t.Error("Should have errored on file that does not exist")
	}
}

func TestLoadingGoodFile(t *testing.T) {
	holder := tides.ObservationHolder{}
	err := holder.LoadDataStore("../data/9410230_annual.xml")
	if err != nil {
		t.Errorf("Should not have errored: %v\n", err)
	}
	expected := 1410 // yes, it's magic
	actual := len(holder.Items)
	if actual != expected {
		t.Errorf("Incorrect number of items, expected=%d, actual=%d", expected, actual)
	}
	expectedState := "CA"
	if state := holder.State; state != expectedState {
		t.Errorf("Expected to be tides from %s, but was: %s", expectedState, state)
	}
	expectedStationType := "Harmonic"
	if holder.StationType != expectedStationType {
		t.Errorf("Expected station type %s but was: %s", expectedStationType, holder.StationType)
	}
}
