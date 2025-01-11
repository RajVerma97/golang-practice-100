package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"github.com/nfnt/resize"
)

type Config struct {
	InputPaths       []string `json:"inputPaths"`
	OutputDirector   string   `json:"outputDirectory"`
	ThumbnailWidth   int      `json:"thumbnailWidth"`
	ThumbnailHeight  int      `json:"thumbnailHeight"`
	Concurrency      int      `json:"concurrency"`
	SupportedFormats []string `json:"supportedFormats"`
}

func loadConfig(configFilePath string) (*Config, error) {
	file, err := os.Open(configFilePath)

	if err != nil {
		return nil, fmt.Errorf("not able to open config.json %s", err)
	}
	defer file.Close()

	var config *Config

	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("not able to decode json config %s", err)
	}
	return config, nil

}

func createFile(inputPath string, config *Config) error {

	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("not able to open input file %s: %v", inputPath, err)
	}
	defer file.Close()

	img, format, err := image.Decode(file)

	if err != nil {
		return fmt.Errorf("not able to decode image %s: %v", inputPath, err)
	}
	isSupported := false

	for _, supportedFormat := range config.SupportedFormats {
		if strings.EqualFold(supportedFormat, format) {
			isSupported = true
			break
		}
	}

	if !isSupported {
		return fmt.Errorf("unsupported format: %s", format)
	}

	resizedImg := resize.Resize(
		uint(config.ThumbnailWidth),
		uint(config.ThumbnailHeight),
		img,
		resize.Lanczos3,
	)

	dirPath := config.OutputDirector
	err = os.MkdirAll(dirPath, 0744)

	if err != nil {
		return fmt.Errorf("not able to create directory %s with err %s", dirPath, err)
	}

	inputPathSplice := strings.Split(inputPath, "/")

	fileName := inputPathSplice[len(inputPathSplice)-1]

	outputFilePath := dirPath + fileName
	outFile, err := os.Create(outputFilePath)

	if err != nil {
		return fmt.Errorf("not able to create file %s", err)
	}
	defer outFile.Close()

	switch format {
	case "jpeg", "jpg":
		err = jpeg.Encode(outFile, resizedImg, nil)
	case "png":
		err = png.Encode(outFile, resizedImg)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
	if err != nil {
		return fmt.Errorf("error saving resized image %s: %v", outputFilePath, err)
	}
	fmt.Printf("Thumbnail created: %s\n", outputFilePath)

	return nil

}
func generateThumnails(config *Config) error {
	inputPaths := config.InputPaths
	var failedPaths []string

	for _, inputPath := range inputPaths {
		if err := createFile(inputPath, config); err != nil {
			fmt.Printf("Failed to process %s: %v\n", inputPath, err)
			failedPaths = append(failedPaths, inputPath)
		}
	}

	if len(failedPaths) > 0 {
		return fmt.Errorf("failed to process %d files", len(failedPaths))
	}

	return nil
}

func main() {

	configFilePath := "./thumbnail-generator/config.json"

	config, err := loadConfig(configFilePath)

	if err != nil {
		fmt.Println(err)
	}

	err = generateThumnails(config)
	if err != nil {
		fmt.Printf("err %s", err)
	}

}
