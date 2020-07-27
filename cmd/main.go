package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"prevalentcolors/pkg/datareader"
	"prevalentcolors/pkg/datawriter"
	"prevalentcolors/pkg/imageprocessor"
)

func main() {

	dr, err := datareader.NewMMapReader("../data/input/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer dr.Close()

	ip := imageprocessor.NewImageProcessor()

	dw, err := datawriter.NewFileWriter("../data/output/result.csv")
	if err != nil {
		log.Fatal(err)
	}

	de, err := datawriter.NewFileWriter("../data/output/error.txt")
	if err != nil {
		log.Fatal(err)
	}

	for dr.Read() {
		data := dr.Data()

		scanner := bufio.NewScanner(bytes.NewReader(data))

		// var wg sync.WaitGroup

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
