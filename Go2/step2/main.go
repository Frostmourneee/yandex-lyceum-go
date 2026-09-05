package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"time"
)

func main() {
	f, _ := os.Open("literature.txt")
	offset := int64(1024)
	f.Seek(offset, 0)
}

func ReadContent(filename string) string {
	data, err := os.ReadFile(filename)
	if err != nil {
		return ""
	}
	return string(data)
}

func LineByNum(inputFilename string, lineNum int) string {
	f, _ := os.Open(inputFilename)
	fileScanner := bufio.NewScanner(f)

	isSuccessfulRead := false
	for i := 0; i < lineNum+1; i++ {
		isSuccessfulRead = fileScanner.Scan()
		// if !isSuccessfulRead {
		// 	return ""
		// }
	}
	if !isSuccessfulRead {
		return ""
	}
	return fileScanner.Text()
}

func CopyFilePart(inputFilename, outFileName string, startpos int) error {
	f, err := os.OpenFile(inputFilename, os.O_RDONLY, 0600)
	if err != nil {
		return err
	}
	f.Seek(int64(startpos), 0)
	defer f.Close()

	outFile, err := os.Create(outFileName)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, f)
	if err != nil {
		return err
	}
	return nil
}

func ModifyFile(filename string, pos int, val string) {
	f, _ := os.OpenFile(filename, os.O_WRONLY, 0600)
	f.Seek(int64(pos), 0)
	f.WriteString(val)
}

func ExtractLog(inputFileName string, start, end time.Time) ([]string, error) {
	f, err := os.Open(inputFileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data := make([]string, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		log := scanner.Text()

		datestr := log[:10]
		date, err := time.Parse("02.01.2006", datestr)
		if err != nil {
			return nil, err
		}
		if !date.Before(start) && !date.After(end) {
			data = append(data, log)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, errors.New("no logs found")
	}

	return data, nil
}
