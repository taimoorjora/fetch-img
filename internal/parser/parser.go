package parser

import (
	"encoding/csv"
	"fmt"
	"os"
)

type Image struct {
	Name   string
	Link   string
	Path   string
	Status string
}

func ParseCSV(filepath, nameCol, linkCol string) ([]Image, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading header: %w", err)
	}

	nameIdx := -1
	linkIdx := -1
	for i, col := range header {
		if col == nameCol {
			nameIdx = i
		}
		if col == linkCol {
			linkIdx = i
		}
	}

	if nameIdx == -1 || linkIdx == -1 {
		return nil, fmt.Errorf("required columns not found: name=%d, link=%d", nameIdx, linkIdx)
	}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading records: %w", err)
	}

	images := make([]Image, 0, len(records))
	for _, record := range records {
		if len(record) <= nameIdx || len(record) <= linkIdx {
			continue
		}
		images = append(images, Image{
			Name:   record[nameIdx],
			Link:   record[linkIdx],
			Status: "pending",
		})
	}

	return images, nil
}

func WriteReport(images []Image, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating report file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"name", "link", "path", "status"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing header: %w", err)
	}

	for _, img := range images {
		record := []string{img.Name, img.Link, img.Path, img.Status}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing record: %w", err)
		}
	}

	return nil
}
