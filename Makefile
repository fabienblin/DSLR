
all:	dependencies
	mkdir -p bin
	go build -o bin/describe.exe cmd/describe/describe.go cmd/describe/stat.go
	go build -o bin/scatter_plot.exe cmd/visual/scatter_plot.go
	go build -o bin/histogram.exe cmd/visual/histogram.go
	go build -o bin/pair_plot.exe cmd/visual/pair_plot.go

clean:
	rm -rf ./bin

re : clean all

dependencies:
	go get -u gonum.org/v1/gonum
	go get -u gonum.org/v1/plot

.PHONY: describe