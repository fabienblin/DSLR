package main

import (
	desc "pkg"
	"encoding/csv"
	"log"
	"os"
)

func main() {
	// get csv data
	var fileName string

	if len(os.Args) != 2 {
		log.Fatal("Usage :\n\tdescribe.exe <file.csv>")
		os.Exit(1)
	} else {
		fileName = os.Args[1]
	}

	dataFile, errData := os.Open(fileName)
	if errData != nil {
		log.Fatal("Failed opening file.")
		os.Exit(1)
	}
	defer dataFile.Close()

	r := csv.NewReader(dataFile)

	data, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	description := desc.Describe(data)

	desc.PrintDescription(description)
}
