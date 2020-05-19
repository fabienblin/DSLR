package describe

import (
	"fmt"
	"math"
	"sort"
	"strconv"

	"gonum.org/v1/gonum/mat"
)

var FeatureNames []string
var StatNames []string

func initPrintable(data [][]string) {
	FeatureNames = data[0][6:]
	for i, name := range FeatureNames { // shorten feature names
		if len(FeatureNames[i]) > 10 {
			FeatureNames[i] = name[:9] + "."
		}
	}
	StatNames = []string{"Count", "Mean", "Std", "Min", "25%", "50%", "75%", "Max"}
}

func PrintDescription(description [][]float64) {
	spacing := 15
	// print feature names
	fmt.Printf("%10s", " ")
	for _, feature := range FeatureNames {
		fmt.Printf("% -*s", spacing, feature)
	}
	fmt.Print("\n\n")

	// print values
	for i, stat := range StatNames {
		fmt.Printf("% -10s", stat)
		for j, _ := range FeatureNames {
			fmt.Printf("%-*.4f", spacing, description[i][j])
		}
		fmt.Print("\n")
	}
}

func count(values []float64) float64 {
	var count float64

	for _, num := range values {
		if !math.IsNaN(num) {
			count++
		}
	}
	return count
}

func sum(values []float64) float64 {
	var sum float64

	for _, num := range values {
		if !math.IsNaN(num) {
			sum += num
		}
	}
	return sum
}

func std(values []float64, mean float64, count float64) float64 {
	var std float64
	var variance float64

	for _, num := range values {
		if !math.IsNaN(num) {
			variance += math.Pow(num-mean, 2)
		}
	}
	variance /= count - 1
	std = math.Sqrt(variance)
	return std
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

type ByValue []float64

func (a ByValue) Len() int           { return len(a) }
func (a ByValue) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByValue) Less(i, j int) bool { return math.IsNaN(a[i]) || a[i] < a[j] }

func quartile(values []float64, count float64, centil float64) float64 {
	var quartile float64 = math.Inf(-1)

	sort.Sort(ByValue(values))

	var nbNan int
	for _, num := range values {
		if math.IsNaN(num) {
			nbNan++
		} else {
			break
		}
	}

	limit := nbNan + int(math.Round(((centil * count) / 100)))
	for i, num := range values {
		if i == limit {
			break
		}
		if !math.IsNaN(num) && quartile < num {
			quartile = num
		}
	}
	return quartile
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

func fillMatrix(data [][]string) *mat.Dense {
	nrows := len(data)
	ncols := len(data[0])
	matrix := mat.NewDense(nrows-1, ncols-6, nil)
	for i := 1; i < nrows; i++ {
		for j := 6; j < ncols; j++ {
			float, err := strconv.ParseFloat(data[i][j], 64)
			if err == nil {
				matrix.Set(i-1, j-6, float)
			} else {
				matrix.Set(i-1, j-6, math.NaN()) // empty fields are set to NaN
			}
		}
	}
	return matrix
}

// parse data to printable statistics
func Describe(data [][]string) [][]float64 {
	matrix := fillMatrix(data) // parse and cast data from string to float

	initPrintable(data) // init feature names and stat names

	// allocate stat array
	m := make([][]float64, len(StatNames))
	for i := range m {
		m[i] = make([]float64, len(FeatureNames))
	}

	// count
	for i, _ := range FeatureNames {
		m[0][i] = count(mat.Col(nil, i, matrix))
	}

	// mean
	for i, _ := range FeatureNames {
		m[1][i] = sum(mat.Col(nil, i, matrix)) / m[0][i]
	}

	// std
	for i, _ := range FeatureNames {
		m[2][i] = std(mat.Col(nil, i, matrix), m[1][i], m[0][i])
	}

	// min
	for i, _ := range FeatureNames {
		m[3][i] = min(mat.Col(nil, i, matrix))
	}

	// 25%
	for i, _ := range FeatureNames {
		m[4][i] = quartile(mat.Col(nil, i, matrix), m[0][i], 25)
	}

	// 50%
	for i, _ := range FeatureNames {
		m[5][i] = quartile(mat.Col(nil, i, matrix), m[0][i], 50)
	}

	// 75%
	for i, _ := range FeatureNames {
		m[6][i] = quartile(mat.Col(nil, i, matrix), m[0][i], 75)
	}

	// max
	for i, _ := range FeatureNames {
		m[7][i] = max(mat.Col(nil, i, matrix))
	}
	return m
}
