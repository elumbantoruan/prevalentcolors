package imageprocessor

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"io/ioutil"
	"net/http"

	"github.com/EdlinOrg/prominentcolor"
)

// ImageReader .
type ImageReader interface {
	Read(source string) (*ThreePrevalentColor, error)
}

// ThreePrevalentColor .
type ThreePrevalentColor struct {
	Source    string
	ColorRGB1 string
	ColorRGB2 string
	ColorRGB3 string
}

// String implements Stringer
func (pc *ThreePrevalentColor) String() string {
	return fmt.Sprintf("%s,%s,%s,%s", pc.Source, pc.ColorRGB1, pc.ColorRGB2, pc.ColorRGB3)
}

// ImageProcessor is a concreate type of ImageReader
type ImageProcessor struct{}

// NewImageProcessor creates an instance of ImageProcessor
func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}

// Read implements ImageReader.Read
func (ip *ImageProcessor) Read(source string) (*ThreePrevalentColor, error) {
	resp, err := http.Get(source)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("not found")
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(contents))
	if err != nil {
		return nil, err
	}
	colors, err := prominentcolor.Kmeans(img)
	if err != nil {
		return nil, err
	}

	return ip.convert(source, colors)
}

func (ip *ImageProcessor) convert(source string, colors []prominentcolor.ColorItem) (*ThreePrevalentColor, error) {

	if len(colors) != 3 {
		return nil, errors.New("source doesn't comply with 3 most prevalent colors")
	}

	prevalentColor := ThreePrevalentColor{
		Source:    source,
		ColorRGB1: fmt.Sprintf("#%s", colors[0].AsString()),
		ColorRGB2: fmt.Sprintf("#%s", colors[1].AsString()),
		ColorRGB3: fmt.Sprintf("#%s", colors[2].AsString()),
	}

	return &prevalentColor, nil
}
