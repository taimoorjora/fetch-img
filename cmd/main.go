package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"fetchimg/internal/downloader"
	"fetchimg/internal/parser"
)

func main() {
	inputFile := flag.String("input", "", "Input CSV file containing image links")
	nameCol := flag.String("name-col", "name", "Column name for image names")
	linkCol := flag.String("link-col", "link", "Column name for image links")
	outputDir := flag.String("output", "downloads", "Output directory for downloaded images")
	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Error: input file is required")
		flag.Usage()
		os.Exit(1)
	}

	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	images, err := parser.ParseCSV(*inputFile, *nameCol, *linkCol)
	if err != nil {
		fmt.Printf("Error parsing input file: %v\n", err)
		os.Exit(1)
	}

	downloader := downloader.NewDownloader(*outputDir)
	for i := range images {
		path, err := downloader.DownloadImage(images[i].Name, images[i].Link)
		if err != nil {
			images[i].Status = fmt.Sprintf("error: %v", err)
			fmt.Printf("Error downloading %s: %v\n", images[i].Name, err)
			continue
		}
		images[i].Path = path
		images[i].Status = "success"
		fmt.Printf("Successfully downloaded: %s\n", images[i].Name)
	}

	reportPath := filepath.Join(*outputDir, "download_report.csv")
	if err := parser.WriteReport(images, reportPath); err != nil {
		fmt.Printf("Error writing report: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Completed: %s\n", reportPath)
}
