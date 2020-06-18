package main

import (
	"encoding/csv"
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {

	var fileName string

	if len(os.Args) != 2 {
		log.Fatal("Usage :\n\thistogram.exe <file.csv>")
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

	// featureNames := data[0][6:]
	matrix := fillMatrix(data)

	for hous, m := range matrix {
		nrows, ncols := m.Dims()
		points := make(plotter.XYs, nrows)
		for i := 0; i < ncols-1; i++ {
			firstFeature := mat.Col(nil, i, m)
			for j := 1; j < ncols; j++ {
				secondFeature := mat.Col(nil, j, m)
				for k := 0; k < nrows; k++ {
					points[k].X = firstFeature[k]
					points[k].Y = secondFeature[k]
				}
			}
		}
	}
	/////////////////////////////////////////////////////
	// Get some random points

	scatterData := points

	// Create a new plot, set its title and
	// axis labels.
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Points Example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	// Draw a grid behind the data
	p.Add(plotter.NewGrid())

	// Make a scatter plotter and set its style.
	s, err := plotter.NewScatter(scatterData)
	if err != nil {
		panic(err)
	}
	s.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 255}

	p.Add(s)
	p.Legend.Add("scatter", s)

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "points.png"); err != nil {
		panic(err)
	}
}

func randomPoints(n int) plotter.XYs {
	pts := make(plotter.XYs, n)
	for i := range pts {
		if i == 0 {
			pts[i].X = rand.Float64()
		} else {
			pts[i].X = pts[i-1].X + rand.Float64()
		}
		pts[i].Y = pts[i].X + 10*rand.Float64()
	}
	return pts
}

func fillMatrix(data [][]string) map[string]*mat.Dense {
	firstCol := 6
	nrows := len(data)
	ncols := len(data[0])
	matrix := make(map[string]*mat.Dense)

	for i := 1; i < nrows; i++ {
		house := data[i][1]

		if len(house) > 0 {
			toAdd := mat.NewDense(1, ncols-firstCol, nil)
			for j := firstCol; j < ncols; j++ { // fill the array
				grade, err := strconv.ParseFloat(data[i][j], 64)

				if err == nil {
					toAdd.Set(0, j-firstCol, grade)
				} else {
					toAdd.Set(0, j-firstCol, math.NaN())
				}
			}
			if matrix[house] == nil { // array init for the 4 houses
				matrix[house] = mat.DenseCopyOf(toAdd)
			} else {
				var x mat.Dense
				x.Stack(matrix[house], toAdd)
				matrix[house] = &x
			}
		}
	}
	return matrix
}
