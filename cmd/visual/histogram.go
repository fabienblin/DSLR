package main

import (
	"encoding/csv"
	"log"
	"math"
	"os"
	"sort"
	"strconv"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

type Value []float64

func (a Value) Len() int           { return len(a) }
func (a Value) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Value) Less(i, j int) bool { return math.IsNaN(a[i]) || a[i] < a[j] }

func main() {
	// get csv data
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

	featureNames := data[0][6:]
	matrix := fillMatrix(data)
	nbPortions := 10

	houseFrequenceTab := houseDistribution(featureNames, matrix, nbPortions)

	for _, subject := range featureNames {
		plotHistogram(subject, houseFrequenceTab, nbPortions)
	}
}

func houseDistribution(featureNames []string, matrix map[string]*mat.Dense, nbPortions int) map[string]map[string][]float64 {
	houseFrequenceTab := make(map[string]map[string][]float64) // contains grade frequencies of each house in a float array (houseFrequenceTab[house][subject]=[]float64{...})

	for house, houseGrades := range matrix { // loop over houses
		_, ncols := houseGrades.Dims()
		frequenceTab := make(map[string][]float64)
		for i := 0; i < ncols; i++ { // loop over subject grades
			subject := featureNames[i]
			subjectGrades := mat.Col(nil, i, houseGrades)
			subjectGrades = cleanTab(subjectGrades)

			distrib := distribution(subjectGrades, nbPortions)

			frequenceTab[subject] = distrib
		}
		houseFrequenceTab[house] = frequenceTab
	}
	return houseFrequenceTab
}

func distribution(tab []float64, nbPortions int) []float64 {
	frequencesTab := make([]float64, nbPortions)
	max := max(tab)
	min := min(tab)
	diff := max - min
	portion := diff / float64(nbPortions)

	for p := 0; p < nbPortions; p++ {
		high := min + portion*(float64(p)+1)
		low := min + portion*float64(p)
		var freq float64
		for _, val := range tab {
			if val >= low && val <= high {
				freq++
			}
		}
		frequencesTab[p] = freq
	}
	return frequencesTab
}

// sort tab and eliminate NaN values
func cleanTab(tab []float64) []float64 {
	sort.Sort(Value(tab))
	var nbNan int
	for _, num := range tab {
		if math.IsNaN(num) {
			nbNan++
		} else {
			break
		}
	}
	cleantab := tab[nbNan:]
	return cleantab
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

func max(values []float64) float64 {
	var max float64 = math.Inf(-1)

	for _, num := range values {
		if !math.IsNaN(num) && max < num {
			max = num
		}
	}
	return max
}

func min(values []float64) float64 {
	var min float64 = math.Inf(1)

	for _, num := range values {
		if !math.IsNaN(num) && min > num {
			min = num
		}
	}
	return min
}

// create image of subject frequencies with all houses houseFrequenceTab[house][subject]=[]float64{...}
func plotHistogram(subject string, houseFrequenceTab map[string]map[string][]float64, nbPortions int) {
	subjectTab := make(map[string][]float64) // subjectTab[house]=[]float64{frequencies}
	for house, subjectFreq := range houseFrequenceTab {
		subjectTab[house] = subjectFreq[subject]
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = subject
	p.X.Label.Text = "Tranche de notes"
	p.Y.Label.Text = "Fréquence"

	w := vg.Points(10)
	i := 0
	for house, frequencies := range subjectTab {
		var plotValues plotter.Values = frequencies
		bar, err := plotter.NewBarChart(plotValues, w)
		if err != nil {
			panic(err)
		}
		bar.LineStyle.Width = vg.Length(0)
		bar.Color = plotutil.Color(i)
		bar.Offset = vg.Points(w.Points() * (float64(i) - 2))
		p.Add(bar)
		p.Legend.Add(house, bar)
		i++
	}

	p.Legend.Top = true

	folder := "./visual/histograms/"
	os.MkdirAll(folder, 0777)
	if err := p.Save(vg.Length(nbPortions)*2*vg.Centimeter+vg.Length(2), 10*vg.Centimeter, folder+subject+".png"); err != nil {
		panic(err)
	}

}
