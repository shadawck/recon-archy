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

	"github.com/tebeka/selenium"
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

/* Convert WebElement slice to String slice
 */
func WbToString(wb []selenium.WebElement) []string {
	slice := make([]string, len(wb))
	for i := range wb {
		slice[i], _ = wb[i].Text()
	}
	return slice
}

/* Convert a slice of WebElement Attribute to string slice
 */
func WbAttrToString(wb []selenium.WebElement, attr string) []string {
	slice := make([]string, len(wb))
	for i := range wb {
		slice[i], _ = wb[i].GetAttribute(attr)
	}
	return slice
}

/* Scrolling simulation
   TODO -> check if scrolling with loadScript in JS is faster or not
*/
func Scroll(wd selenium.WebDriver) {
	wd.SetImplicitWaitTimeout(time.Second * 3)
	time.Sleep(time.Second * 2)
	wd.KeyDown(selenium.PageDownKey)
	wd.KeyUp(selenium.PageDownKey)
	time.Sleep(time.Second * 2)
}

/* Decode the search URL and build the search URL for the main company
 */
func DecodeReconstruct(searchURL string) string {
	decodedURL, err := url.QueryUnescape(searchURL)
	if err != nil {
		log.Fatal(err)
	}

	spl := strings.Split(decodedURL, "\"")
	for i, v := range spl {
		fmt.Printf("%d, %s\n", i, v)
	}

	toEncode := "\"" + spl[1] + "\"" + "]"
	facetCompanies := url.QueryEscape(toEncode)

	encodedBuild := spl[0] + facetCompanies
	return encodedBuild
}
