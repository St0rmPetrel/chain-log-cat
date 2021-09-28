package main

import (
	"time"

	"github.com/St0rmPetrel/chain-log-cat/processors/filetracker"
)

func main() {
	ft := filetracker.New("..", "", 5*time.Second)
	ft.TrackChanges()
}
