package repository

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus" // nolint
)

const (
	MaxCCDWords  = 8
	CCDSeparator = "-"
)

type CCD struct {
}

func (c *CCD) ReadCCDs(directory string) (ccds []string) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".md") {
			ccds = append(ccds, f.Name())
		}
	}
	return
}

func (c *CCD) CreateCCD(text, dpath string) (name string, err error) {
	var (
		words []string
	)
	for _, line := range strings.Split(text, "\n") {
		for _, word := range strings.Split(line, " ") {
			words = append(words, strings.ToLower(word))
		}
	}
	lastCCDID := int64(0)
	ccds := c.ReadCCDs(dpath)
	if len(ccds) > 0 {
		lastCCD := ccds[len(ccds)-1]

		lastCCDID, err = strconv.ParseInt(strings.Split(lastCCD, "-")[0], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
	}
	maxWords := len(words)

	if len(words) >= MaxCCDWords {
		maxWords = MaxCCDWords
	}

	// 6 zeroes
	name = strings.TrimRight(fmt.Sprintf("%06d%s", lastCCDID+1, CCDSeparator)+strings.Join(words[:maxWords],
		CCDSeparator), CCDSeparator) + ".md"
	return
}

func (c *CCD) WriteToCCD(text, filename string) error {
	currentTime := time.Now()
	date := currentTime.Format("2006-01-02")
	template := fmt.Sprintf("Date: %s\n", date)
	template += "## Decision\n" + text
	return ioutil.WriteFile(filename, []byte(template), 0644)
}

func NewCCD() (ccd *CCD) {
	ccd = &CCD{}
	return
}
