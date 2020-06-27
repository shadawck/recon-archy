package selenium

import (
	"bufio"
	"io/ioutil"
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
