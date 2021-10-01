package main

import (
	"fmt"
	"time"

	"github.com/St0rmPetrel/chain-log-cat/processors/filetracker"
)

func main() {
	ft := filetracker.New("..", "", 100*time.Second)
	b, _ := ft.TrackChanges()
	fmt.Println("\n" + string(b))
}
