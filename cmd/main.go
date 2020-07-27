package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"prevalentcolors/pkg/datareader"
	"prevalentcolors/pkg/datawriter"
	"prevalentcolors/pkg/imageprocessor"
)

func main() {

	// takes an argument and provides default value
	input := flag.String("input", "../data/input/input.txt", "input file")
	output := flag.String("output", "../data/output/result.csv", "output file")
	errorLog := flag.String("errorlog", "../data/output/error.txt", "error log")
	chunkSize := flag.Int("chunksize", 4096, "chunk size for data reader")

	flag.Parse()

	dr, err := datareader.NewMMapReader(*input, *chunkSize)
	if err != nil {
		log.Fatal(err)
	}
	defer dr.Close()

	ip := imageprocessor.NewImageProcessor()

	dw, err := datawriter.NewFileWriter(*output)
	if err != nil {
		log.Fatal(err)
	}

	de, err := datawriter.NewFileWriter(*errorLog)
	if err != nil {
		log.Fatal(err)
	}

	for dr.Read() {
		data := dr.Data()

		scanner := bufio.NewScanner(bytes.NewReader(data))

		var urls []string
		for scanner.Scan() {
			urls = append(urls, scanner.Text())
		}
		if scanner.Err() != nil {
			log.Fatal(err)
		}

		ch := make(chan struct{})

		for _, url := range urls {

			go func(url string) {

				clrs, err := ip.Read(url)
				if err != nil {
					// write error
					de.Write([]byte(fmt.Sprintf("error getting image from:%s -- %s \n", url, err.Error())))
				} else {
					// write result
					err = dw.Write([]byte(fmt.Sprintf("%s\n", clrs)))
					if err != nil {
						log.Fatal(err)
					}
				}

				ch <- struct{}{}
			}(url)
		}

		//closer
		for range urls {
			<-ch
		}
	}

	if dr.Err() != nil && dr.Err() != io.EOF {
		log.Fatal(dr.Err())
	}

}
