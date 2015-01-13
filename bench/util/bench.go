package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ParseFunc func([]byte) bool

var Counter = 0
var Time = 0.0

func BuildParsingVisitor(parseFile ParseFunc) filepath.WalkFunc {
	Counter = 0
	visit := func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			if strings.Contains(path, ".DS_Store") {
				return nil
			}

			if strings.Contains(path, ".txt") {
				return nil
			}

			data, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Println(err)
				os.Exit(3)
			}

			start := time.Now().UnixNano()
			if parseFile(data) {
				end := time.Now().UnixNano()
				seconds := float64(end-start) / 1000000000.0

				Time += seconds
				Counter++

				fmt.Printf("%s,%f\n", path, seconds)
			} else {
				log.Printf("Failed to parse: %s\n", path)
			}
		}

		return nil
	}
	return visit
}
