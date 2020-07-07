package selenium

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/matryer/try"
	"github.com/tebeka/selenium"
)

// ReadFile Read a file and return the line as a []string
func ReadFile(f *os.File) []string {
	defer f.Close()
	reader := bufio.NewReader(f)
	contents, _ := ioutil.ReadAll(reader)
	lines := strings.Split(string(contents), "\n")
	return lines
}

// WriteFile Write []string to file "fileName" and line by line
func WriteFile(fileName string, buff []string) {
	f, err := os.OpenFile(fileName+".archy", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	for _, v := range buff {
		fmt.Fprintln(f, v) // Right to file and append new line
		if err != nil {
			panic(err)
		}
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}

}

//SlicePrint Print String Slice
func SlicePrint(s []string) {
	fmt.Print("\n")
	for i, v := range s {
		fmt.Printf("%d, %s\n", i, v)
	}
}

// DecodeReconstruct Decode the searchURL and build the search URL for the main company
func DecodeReconstruct(searchURL string) string {
	decodedURL, err := url.QueryUnescape(searchURL)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	spl := strings.Split(decodedURL, "\"")
	toEncode := "\"" + spl[1] + "\"" + "]"
	facetCompanies := url.QueryEscape(toEncode)
	encodedBuild := spl[0] + facetCompanies
	return encodedBuild
}

// DecodeRetry Attempt getting currentURL
func DecodeRetry(wd selenium.WebDriver, currentURL string) string {

	var encodedURL string
	err := try.Do(

		func(attempt int) (bool, error) {
			var err error
			encodedURL = DecodeReconstruct(currentURL)

			// Wait 100 ms between each retry (to load the page)
			if err != nil {
				wd.SetImplicitWaitTimeout(time.Millisecond * 1000)
			}
			// attempt < 5 -> try 5 time
			return attempt < 5, err
		})

	if err != nil {
		log.Fatalln("error:", err)
		wd.Quit()
	}
	return encodedURL
}
