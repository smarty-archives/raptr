package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	if store, err := NewEventStore("test.json"); err != nil {
		fmt.Println(err)
	} else if record, err := store.Next(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Print(record.MetaRecord())
		text, _ := json.Marshal(record.Message)
		fmt.Println(string(text))
		//fmt.Printf("%#v\n", record)

		if record, err = store.Next(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Print(record.MetaRecord())
			text, _ := json.Marshal(record.Message)
			fmt.Println(string(text))
			// fmt.Printf("%#v\n", record)
		}
	}
}

// various projections for different purposes:
// 1. validation of incoming input
// 2. generation of apt index files

// all events provided to all projections (simple)
// and they decide what to do
