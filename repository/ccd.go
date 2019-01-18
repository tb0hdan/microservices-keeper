package repository

import (
	"fmt"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

type CCD struct {
}

func (c *CCD) readACCDs(directory string) (ccds []string) {
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

func (c *CCD) createCCD(text, dpath string) (name string, err error) {
	var (
		words []string
	)
	for _, line := range strings.Split(text, "\n") {
		for _, word := range strings.Split(line, " ") {
			words = append(words, strings.ToLower(word))
		}
	}
	last_ccd_id := int64(0)
	ccds := c.readACCDs(dpath)
	if len(ccds) > 0 {
		last_ccd := ccds[len(ccds)-1]

		last_ccd_id, err = strconv.ParseInt(strings.Split(last_ccd, "-")[0], 10, 64)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 6 zeroes
	name = fmt.Sprintf("%06d-", last_ccd_id+1) + strings.Join(words[:MAX_CCD_WORDS], "-") + ".md"
	return
}

func (c *CCD) writeToCCD(text, filename string) error {
	currentTime := time.Now()
	date := currentTime.Format("2006-01-02")
	template := fmt.Sprintf("Date: %s\n", date)
	template += "## Decision\n" + text
	return ioutil.WriteFile(filename, []byte(template), 0644)
}
