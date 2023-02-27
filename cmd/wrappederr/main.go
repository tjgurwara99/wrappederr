package main

import (
	"github.com/tjgurwara99/wrappederr"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(wrappederr.Analyzer)
}
