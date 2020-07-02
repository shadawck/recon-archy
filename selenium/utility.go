package selenium

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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
