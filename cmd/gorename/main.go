package main

import (
	"flag"

	"github.com/Orogenesis/gorename/pkg/replacer"
)

var (
	rootDir     string
	prevPath    string
	newPath     string
	useModules  bool
	printResult bool
)

func init() {
	flag.StringVar(&rootDir, "root-dir", ".", "root directory")
	flag.StringVar(&prevPath, "path", "", "current path")
	flag.StringVar(&newPath, "new-path", "", "new path")
	flag.BoolVar(&useModules, "use-modules", true, "replace module name in go.mod")
	flag.BoolVar(&printResult, "print-result", true, "print result")
	flag.Parse()
}

func main() {
	r := replacer.New(
		prevPath,
		newPath,
		replacer.UseModules(useModules),
		replacer.PrintResult(printResult),
		replacer.RootDir(rootDir),
	)
	err := r.Run()
	if err != nil {
		panic(err)
	}
}
