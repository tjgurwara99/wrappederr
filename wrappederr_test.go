package wrappederr_test

import (
	"testing"

	"github.com/tjgurwara99/wrappederr"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, wrappederr.Analyzer, "withoutwrap")
}
