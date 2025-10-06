package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/wardviaene/meetparser/pkg/parser"
)

func main() {
	var filename string
	flag.StringVar(&filename, "filename", "", "parse filename")

	flag.Parse()

	if filename == "" {
		flag.PrintDefaults()
		os.Exit(0)
	}

	filenameWithoutSuffix := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))

	if !fileExists("pdf-column-extractor-1.0-SNAPSHOT.jar") {
		fmt.Printf("pdf-column-extractor-1.0-SNAPSHOT.jar doesn't exist. Build the jar file first and place it in the current directory")
		os.Exit(1)
	}

	cmd := exec.Command(
		"java",
		"-jar", "pdf-column-extractor-1.0-SNAPSHOT.jar",
		"./"+filename,
		"./"+filenameWithoutSuffix+".txt",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error processing %s: %s\n%s.\n", filename, output, err)
	}

	result, err := parser.ParsePDFText(filenameWithoutSuffix + ".txt")
	if err != nil {
		log.Fatalf("Error extracting text: %v", err)
	}

	// write times
	if len(result.Times) > 0 {
		csvBytes, err := parser.MarshalCSV(result.Times)
		if err != nil {
			log.Fatalf("Error creating csv: %s", err)
		}

		err = os.WriteFile(filenameWithoutSuffix+"-times.csv", csvBytes, 0644)
		if err != nil {
			log.Fatalf("Error creating csv file: %s", err)

		}
	}

	// write relays
	if len(result.RelayTimes) > 0 {
		csvBytes, err := parser.MarshalCSV(result.RelayTimes)
		if err != nil {
			log.Fatalf("Error creating csv (relays): %s", err)
		}

		err = os.WriteFile(filenameWithoutSuffix+"-relays.csv", csvBytes, 0644)
		if err != nil {
			log.Fatalf("Error creating csv file (relays): %s", err)

		}
	}

	// write events
	if len(result.Events) > 0 {
		csvBytes, err := parser.MarshalCSV(result.Events)
		if err != nil {
			log.Fatalf("Error creating csv (events): %s", err)
		}

		err = os.WriteFile(filenameWithoutSuffix+"-events.csv", csvBytes, 0644)
		if err != nil {
			log.Fatalf("Error creating csv file (events): %s", err)

		}
	}

	fmt.Println("CSV written.")
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}
