package main

func main() {
}

// two lines per message: first line is metadata, e.g. name,seq,date
// second line is payload, e.g.:
//// PackageRemoved,1,ISO time
//{}

// various projections for different purposes:
// 1. validation of incoming input
// 2. generation of apt index files

// all events provided to all projections (simple)
// and they decide what to do
