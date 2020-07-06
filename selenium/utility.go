package selenium

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"
)

/*
Read Files

*/
func ReadFile(f *os.File) []string {
	defer f.Close()
	reader := bufio.NewReader(f)
	contents, _ := ioutil.ReadAll(reader)
	lines := strings.Split(string(contents), "\n")
	return lines
}

/*Print String Slice
 */
func SlicePrint(s []string) {
	fmt.Print("\n")
	for i, v := range s {
		fmt.Printf("%d, %s\n", i, v)
	}
}

/* Decode the search URL and build the search URL for the main company
 */
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
