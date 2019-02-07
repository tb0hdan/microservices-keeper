package ccd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/tb0hdan/microservices-keeper/repository/logs"
)

const (
	MaxCCDWords     = 8
	CCDSeparator    = "-"
	PreviousCCDLink = "[<-previous]"
	PreviousCCD     = PreviousCCDLink + "(%s)"
	NextCCDLink     = "[next->]"
	NextCCD         = "|\n" + NextCCDLink + "(%s)"
	NoNextCCD       = "|\n" + NextCCDLink + "(None)"
)

type CCD struct {
	// TODO: Consider adding working directory here
}

func (c *CCD) ReadCCDs(directory string) (ccds []string) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		logs.Logger.Fatal(err)
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".md") {
			ccds = append(ccds, f.Name())
		}
	}
	return
}

func (c *CCD) ReadCCD(directory, ccd string) (contents string, err error) {
	tmp, err := ioutil.ReadFile(path.Join(directory, ccd))
	if err != nil {
		return "", err
	}
	contents = string(tmp)
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
			logs.Logger.Fatal(err)
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

func (c *CCD) HasHeader(directory, ccd string) bool {
	contents, err := c.ReadCCD(directory, ccd)
	if err != nil {
		return false
	}
	for _, line := range strings.Split(contents, "\n") {
		if strings.Contains(line, PreviousCCDLink) || strings.Contains(line, NextCCDLink) {
			return true
		}
	}
	return false
}

func (c *CCD) RebuildIndex(directory string) error { // nolint
	var (
		CCDHeader string
	)
	ccds := c.ReadCCDs(directory)
	for idx, ccd := range ccds {
		// First CCD
		if idx == 0 && len(ccds) > 1 {
			CCDHeader = fmt.Sprintf(NextCCD, ccds[idx+1])
		} else if idx == 0 && len(ccds) == 1 {
			CCDHeader = NoNextCCD
		}
		// Other CCDs
		if idx > 0 && len(ccds) > 2 && idx+1 < len(ccds) {
			CCDHeader = fmt.Sprintf(PreviousCCD, ccds[idx-1]) + "\n" + fmt.Sprintf(NextCCD, ccds[idx+1])
		} else if idx > 0 && len(ccds) > 2 {
			CCDHeader = fmt.Sprintf(PreviousCCD, ccds[idx-1])
		}
		// Check header
		if c.HasHeader(directory, ccd) {
			// TODO: Update index properly
			continue
		}
		// Got header, proceed with rebuilding index
		tmpfile, err := ioutil.TempFile("", "microservices-keeper")
		if err != nil {
			return err
		}
		ccdContents, err := c.ReadCCD(directory, ccd)
		if err != nil {
			return err
		}

		content := CCDHeader + ccdContents
		if _, err := tmpfile.Write([]byte(content)); err != nil {
			return err
		}
		if err := tmpfile.Close(); err != nil {
			return err
		}
		if err := os.Rename(tmpfile.Name(), path.Join(directory, ccd)); err != nil {
			return err
		}
	}
	return nil
}

func NewCCD() (ccd *CCD) {
	ccd = &CCD{}
	return
}
