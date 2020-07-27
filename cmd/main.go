package main

import (
	"bufio"
	"bytes"
	"fmt"
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

		for scanner.Scan() {

			text := scanner.Text()
			if len(text) == 0 {
				continue
			}
			clrs, err := ip.Read(text)
			if err != nil {
				// write error
				de.Write([]byte(fmt.Sprintf("error getting image from:%s -- %s \n", text, err.Error())))
				continue
			}
			// write result
			err = dw.Write([]byte(fmt.Sprintf("%s\n", clrs)))
			if err != nil {
				log.Fatal(err)
			}
		}
		if scanner.Err() != nil {
			log.Fatal(err)
		}
	}

	if dr.Err() != nil {
		log.Fatal(dr.Err())
	}

}
