package main

import (
	"encoding/csv"
	"log"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
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

	description := describe(data)

	printDescription(description)

	groupA := plotter.Values{20, 35, 30, 35, 27}
	groupB := plotter.Values{25, 32, 34, 20, 25}
	groupC := plotter.Values{12, 28, 15, 21, 8}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Bar chart"
	p.Y.Label.Text = "Heights"

	w := vg.Points(20)

	barsA, err := plotter.NewBarChart(groupA, w)
	if err != nil {
		panic(err)
	}
	barsA.LineStyle.Width = vg.Length(0)
	barsA.Color = plotutil.Color(0)
	barsA.Offset = -w

	barsB, err := plotter.NewBarChart(groupB, w)
	if err != nil {
		panic(err)
	}
	barsB.LineStyle.Width = vg.Length(0)
	barsB.Color = plotutil.Color(1)

	barsC, err := plotter.NewBarChart(groupC, w)
	if err != nil {
		panic(err)
	}
	barsC.LineStyle.Width = vg.Length(0)
	barsC.Color = plotutil.Color(2)
	barsC.Offset = w

	p.Add(barsA, barsB, barsC)
	p.Legend.Add("Group A", barsA)
	p.Legend.Add("Group B", barsB)
	p.Legend.Add("Group C", barsC)
	p.Legend.Top = true
	p.NominalX("One", "Two", "Three", "Four", "Five")

	if err := p.Save(5*vg.Inch, 3*vg.Inch, "barchart.png"); err != nil {
		panic(err)
	}
}
