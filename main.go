package main

import (
	"fmt"
	"os"
)

func main() {
	if handle, err := os.Open("test.json"); err != nil {
		fmt.Println(err)
	} else if store := NewEventStoreReader(handle); store == nil {
		fmt.Println("Couldn't create eventstore")
	} else {
		writer := NewEventStoreWriter(0, os.Stdout)

		for {
			if record, err := store.Read(); err != nil {
				fmt.Println("Reading:", err)
				break
			} else if err := writer.Write(record); err != nil {
				fmt.Println("Writing:", err)
				break
			}
		}
	}
}

// various projections for different purposes:
// 1. validation of incoming input
// 2. generation of apt index files

// all events provided to all projections (simple)
// and they decide what to do
