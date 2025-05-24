# Image Downloader

A simple CLI tool to download images from URLs listed in a CSV file.

## Usage
```csv
name,link
image1,https://example.com/image1.jpg
image2,https://example.com/image2.png
```

## Command Line Arguments

- `-input`: Path to the input CSV file (required)
- `-name-col`: Name of the column containing image names (default: "name")
- `-link-col`: Name of the column containing image URLs (default: "link")
- `-output`: Directory to save downloaded images (default: "downloads")

## Output

The tool will:
1. Download images to the specified output directory
2. Generate a report CSV file with download status
3. Show progress in the console