package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

type flagParams struct {
	sourceFile string
	lineCount  int
	destFile   string
	maxFiles   int
}

func main() {
	sourceFile := flag.String("i", "", "File to be split.")
	lineCount := flag.Int("l", 0, "Maximum lines file to be split into")
	destFile := flag.String("o", "", "Destination file name.")
	maxFiles := flag.Int("m", 0, "Maximum number of files to be output. (0 for all)")
	flag.Parse()

	params := flagParams{*sourceFile, *lineCount, *destFile, *maxFiles}

	// Check for errors with the cli parameters
	err := params.checkFlagErrors()
	if err != nil {
		flag.PrintDefaults()
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Source file: %v \n", params.sourceFile)
	fmt.Printf("Destination file: %v \n", params.destFile)

	fileCount, err := params.splitFile()
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Split complete. %v files.\n", fileCount)

}

func (param flagParams) checkFlagErrors() error {
	var err error

	if param.sourceFile == "" || param.destFile == "" {
		return fmt.Errorf("the source and destination file names cannot be blank")
	}

	if param.sourceFile == param.destFile && err == nil {
		return fmt.Errorf("the source and destination files must be different")
	}

	if param.lineCount < 1 {
		return fmt.Errorf("the file cannot be split to less than 1 line per file")
	}

	if param.maxFiles < 0 {
		return fmt.Errorf("maximum file count must be zero (maximum files) or greater")
	}

	return err

}

func (param flagParams) incFilename(counter int) string {

	extn := filepath.Ext(param.destFile)
	return param.destFile[0:len(param.destFile)-len(extn)] + strconv.Itoa(counter) + extn

}

func (param flagParams) splitFile() (int, error) {

	var err error
	var fileCount int

	// Open the source file for reading
	fileReader, err := os.Open(param.sourceFile)
	if err != nil {
		return fileCount, err
	}
	defer fileReader.Close()
	bufioReader := bufio.NewReader(fileReader)

	// Open the Destination file for writing
	for {
		fileCount++
		outputWriter, err := os.Create(param.incFilename(fileCount))
		if err != nil {
			err := fmt.Errorf("creating output file: %v", err)
			return fileCount, err
		}
		defer outputWriter.Close()

		owriter := bufio.NewWriter(outputWriter)
		defer owriter.Flush()

		for i := 0; i < param.lineCount; i++ {
			fileLines, readerr := bufioReader.ReadString('\n')
			if readerr == io.EOF {
				return fileCount, err
			}
			if _, err := owriter.WriteString(fileLines); err != nil {
				return fileCount, fmt.Errorf("writing to output file: %v", err)
			}
		}
		if param.maxFiles > 0 && fileCount >= param.maxFiles {
			return fileCount, err
		}
	}

}
