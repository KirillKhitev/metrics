package codeanalizer

import (
	"log"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

var TestData = func() string {
	testdata, err := filepath.Abs("../../cmd")
	if err != nil {
		log.Fatal(err)
	}
	return testdata
}

func TestExitOnMainAnalyzer(t *testing.T) {
	analysistest.Run(t, TestData(), ExitOnMainAnalyzer, "./...")
}
