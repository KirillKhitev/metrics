package codeanalizer

import (
	"golang.org/x/tools/go/analysis/analysistest"
	"log"
	"path/filepath"
	"testing"
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
