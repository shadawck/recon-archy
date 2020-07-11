/*
Package selenium implements linkedin data scraping with selenium automatisation.
*/
package selenium

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/matryer/try"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/firefox"
)

const (
	// These paths will be different on your system.
	seleniumPath    = "~/go/pkg/mod/github.com/tebeka/selenium@v0.9.9/vendor/selenium-server-standalone.jar"
	geckoDriverPath = "~/go/pkg/mod/github.com/tebeka/selenium@v0.9.9/vendor/geckodriver"
	htmlunitpath    = "~/go/pkg/mod/github.com/tebeka/selenium@v0.9.9/vendor/htmlunit-driver.jar"
	LITTLE_WAIT     = time.Millisecond * 100
	MEDIUM_WAIT     = time.Second
	BIG_WAIT        = time.Second * 3
	port            = 4444
)

// Get number of search page found
func lenPage(wd selenium.WebDriver) int {
	// Scroll the page to ensure page is entirely loaded
	scroll(wd, 2)

	pageNumber := findsRetry(wd, ".artdeco-pagination__indicator--number.ember-view", "Page Number")

	//fmt.Printf("\nPageNumber findRetry(): %v", pageNumber)

	lenPage, err := pageNumber[len(pageNumber)-1].Text()
	if err != nil {
		panic(err)
	}

	conv, err := strconv.Atoi(lenPage)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nText Page Number is : %s\n", lenPage)

	return conv
}

// Go to next search Page
func nextPage(wd selenium.WebDriver, page int, searchURL string) {
	baseURL := searchURL + "&origin=FACETED_SEARCH&page=" + strconv.Itoa(page)
	wd.Get(baseURL)
}

// Convert WebElement slice to String slice
func wbToStringSlice(wb []selenium.WebElement) []string {
	slice := make([]string, len(wb))
	for i := range wb {
		slice[i], _ = wb[i].Text()
	}

	return slice
}

// Convert a slice of WebElement Attribute to string slice
func wbAttrToStringSlice(wb []selenium.WebElement, attr string) []string {
	slice := make([]string, len(wb))
	for i := range wb {
		slice[i], _ = wb[i].GetAttribute(attr)
	}
	return slice
}

// Scrolling simulation
//   TODO -> check if scrolling with loadScript in JS is faster or not
func scroll(wd selenium.WebDriver, x int) {
	wd.SetImplicitWaitTimeout(time.Second)
	// scroll x time (for headless mode for example)
	for i := 0; i < x; i++ {
		wd.KeyDown(selenium.PageDownKey)
		wd.KeyUp(selenium.PageDownKey)
	}
}

func initService(port int) selenium.WebDriver {

	opts := []selenium.ServiceOption{
		//selenium.StartFrameBuffer(),           // Start an X frame buffer for the browser to run in.
		selenium.GeckoDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		//selenium.Output(os.Stderr),            // Output debug information to STDERR.
	}

	selenium.SetDebug(false)
	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer service.Stop()

	caps := selenium.Capabilities{"browserName": "firefox"}
	caps.AddFirefox(firefox.Capabilities{Args: []string{"--headless", "--safe-mode"}})

	fmt.Printf("\nStart New Remote on port: %s ", strconv.Itoa(port))
	remoteAdress := "http://localhost:" + strconv.Itoa(port) + "/wd/hub"
	wd, err := selenium.NewRemote(caps, remoteAdress)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	fmt.Printf("\nRemote management at : %s", remoteAdress)

	// Wait for Ctr-C or Killing signal
	CloseHandler(wd)

	return wd
}

// CloseHandler handle program Interruption and perform a clean exit
func CloseHandler(wd selenium.WebDriver) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("\n\r- Ctrl+C pressed in Terminal : Cleaning Exit")
		wd.Quit()
		os.Exit(0)
	}()
}

// signIn perform Linkedin SignIn
func signIn(wd selenium.WebDriver) {
	// Navigate to Linkedin
	fmt.Printf("\nNavigating to Signup")
	if err := wd.Get("https://www.linkedin.com/login?fromSignIn=true&trk=guest_homepage-basic_nav-header-signin"); err != nil {
		panic(err)
	}

	// Load Credential for Linkedin SignIn
	fmt.Printf("\nPerform Signup")
	file, err := os.Open(".creds")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read credential file
	creds := ReadFile(file)

	// Find linkedin username field
	username, err := wd.FindElement(selenium.ByID, "username")
	if err != nil {
		panic(err)
	}
	// And fill it with the username
	err = username.SendKeys(creds[0])
	if err != nil {
		panic(err)
	}

	// Find linkedin password field
	pass, err := wd.FindElement(selenium.ByID, "password")
	if err != nil {
		panic(err)
	}

	// And fill it with the password
	err = pass.SendKeys(creds[1])
	if err != nil {
		panic(err)
	}

	// Click the SignIn button
	btn := findRetry(wd, ".btn__primary--large", "SignIn Button")

	if err := btn.Click(); err != nil {
		panic(err)
	}

	fmt.Printf("\nLogged In !")
}

func captchaCheck(wd selenium.WebDriver) bool {
	captcha, err := wd.FindElement(selenium.ByID, "captchaInternalPath")
	if err != nil {
		panic(err)
	}
	if captcha != nil {
		fmt.Printf("\nThere is a CAPTCHA to solve on your account")
		fmt.Printf("\nLogin in GUI mode to solve it and relauch ReconArchy")
		captchaBool := true
		return captchaBool
	}
	captchaBool := false
	return captchaBool

}

// findRetry Retry for findElement function
func findRetry(wd selenium.WebDriver, selector string, id string) selenium.WebElement {

	var found selenium.WebElement
	err := try.Do(

		func(attempt int) (bool, error) {
			var err error
			found, err = wd.FindElement(selenium.ByCSSSelector, selector)

			// Wait 100 ms between each retry
			if err != nil {
				time.Sleep(time.Millisecond * 500)
			}
			// attempt < 5 -> try 5 time
			fmt.Printf("\n (attempts %d) for %v for grabing: %s", attempt, wd.SessionID(), id)
			return attempt < 5, err
		})

	if err != nil {
		wd.Quit()
		panic(err)
	}
	return found
}

// findRetry Retry for findElements function
func findsRetry(wd selenium.WebDriver, selector string, id string) []selenium.WebElement {

	var found []selenium.WebElement
	err := try.Do(

		func(attempt int) (bool, error) {
			var err error
			found, err = wd.FindElements(selenium.ByCSSSelector, selector)

			// Wait 100 ms between each retry (to load the page)
			if err != nil {
				//wd.SetImplicitWaitTimeout(time.Millisecond * 100)
				time.Sleep(time.Millisecond * 500)
			}
			// attempt < 5 -> try 5 time
			fmt.Printf("\n (attempts %d) for %v for grabing: %s", attempt, wd.SessionID(), id)
			return attempt < 5, err
		})

	if err != nil {
		wd.Quit()
		panic(err)
	}
	return found
}

// searchPage process company name and go to filtered page result for crawling
// Implement Filtering Options
func searchFilteredPage(wd selenium.WebDriver, comp string) string {
	/* SOLUTION 1 - INPUT 1 : Navigate to the companies page with just the companies name */
	// Navigate to Linkedin companies search result (for givven company name "comp")

	var searchURL string = "https://www.linkedin.com/search/results/companies/?keywords=" + comp + "&origin=SWITCH_SEARCH_VERTICAL"
	if err := wd.Get(searchURL); err != nil {
		panic(err)
	}

	fmt.Printf("\nSearching Company Page")
	time.Sleep(LITTLE_WAIT)
	firstCompanyLink := findsRetry(wd, ".app-aware-link.ember-view", "company-link ") // click on the first company found in the search result
	fmt.Printf("\nCompany found")

	// click on the second link on the slice. In fact image link can't be clicked so just use the second link
	// TODO Need To be OPTIMIZED -> I get all the "app-aware-link" of the page but i just need the first one
	if err := firstCompanyLink[1].Click(); err != nil {
		panic(err)
	}

	/* SOLUTION 2 - INPUT 2: Navigate directly to the companies page; The url need to be input by the user
	if err := wd.Get(compURL); err != nil {
		panic(err)
	}*/

	// Click on "See all X Employees on Linkedin"
	fmt.Printf("\nGetting %s employees", comp)
	wd.SetImplicitWaitTimeout(LITTLE_WAIT)

	employees := findRetry(wd, ".ember-view.link-without-visited-state.inline-block", "employee")
	if err := employees.Click(); err != nil {
		panic(err)
	}

	// TODO -> Add filter selection to select only wanted companies or subsidiary companies
	// get and process actual search url
	// URL structc: https://www.linkedin.com/search/results/people/?facetCurrentCompany=["1259"%2C"2274"%2C"208298"%2C"1260"%2C"53472064"]
	scroll(wd, 2) // SUPER IMPORRTANT SCROLL
	wd.SetImplicitWaitTimeout(LITTLE_WAIT)

	//fmt.Printf("\nThe current url to be decoded is %s", currentURL)

	fmt.Printf("\nDecoding URL...")
	encodedURL := DecodeRetry(wd)
	fmt.Printf("\nURL Re-encoded...")

	return encodedURL

}

// StartProcess Start an OS Process : Used to start selenium standalone server on n port for multiple webdriver worker
func StartProcess(port int) *os.Process {
	cmd := exec.Command("bash", "-c", "java -jar "+seleniumPath+" -port "+strconv.Itoa(port))
	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start cmd: %v", err)
	}
	fmt.Printf(" \nStartProcess : %v", cmd.Process)
	return cmd.Process
}

// KillProcess Kill process started by StartProcess()
func KillProcess(proc []*os.Process) {
	fmt.Printf("\nKillProcess : %v", proc)
	for _, p := range proc {
		p.Kill()
	}
}

// CreateWorker Spawn "t" webdriver simulate multithreading and reduce runtime
func CreateWorker(currentWd selenium.WebDriver, threadNumber int, initialPort int) []selenium.WebDriver {
	// create a map to store Process
	//var proc []*os.Process

	// Can be useful for later. For now, server are started manualy or from external script
	// each webdriver need a different port
	//for i := 0; i < threadNumber; i++ {
	//	currentProc := StartProcess(initialPort + i)
	//	proc = append(proc, currentProc)
	//}

	// Create WebDriver map to give instruction to each WebDriver
	fmt.Printf("\nInitialising worker")
	var workers []selenium.WebDriver
	workers = append(workers, currentWd)
	for i := 1; i < threadNumber; i++ {
		wd := initService(initialPort + i)
		workers = append(workers, wd)
	}

	return workers
}

func calcSplitting(lenPage int, threadNumber int) [][]int {
	r := lenPage / threadNumber
	//rest := lenPage % threadNumber
	//fmt.Printf("\nrest %d", rest)
	//fmt.Printf("\nStep : %d", r)

	var step [][]int
	var tmp []int

	for i := 0; i < threadNumber; i++ {
		tmp = nil // terrible way to do this. Think of a nicer implementation for callSplitting
		tmp = append(tmp, (i*r)+1)
		tmp = append(tmp, (i+1)*r)
		step = append(step, tmp)
	}
	return step
}

// PerformSelenium action like signIn, searching and crawling but for map of webdriver. Need to
// divide page crawling between webdriver "worker"
// Useless to perform repetitive action against all webDriver. The initial WebDriver perform the basic action
// and transmit results to workers.
func populateWorker(wg *sync.WaitGroup, wd selenium.WebDriver, id string, lenPage int, url string, startPage int, stopPage int, comp string) {
	defer wg.Done()

	// worker 1 (initial wd) is aleardy signin
	if id != "worker_0" {
		fmt.Printf("\n%s signing", id)
		signIn(wd)
	}

	for i := startPage; i <= stopPage; i++ {
		nextPage(wd, i, url)
		scroll(wd, 2)
		fmt.Printf("\n%s try to extract page number %d", id, i)
		//wd.SetImplicitWaitTimeout(LITTLE_WAIT)
		//users := findsRetry(wd, ".actor-name")
		//
		//usersText := wbToStringSlice(users)
		//SlicePrint(usersText)
		//
		//profileURL := findsRetry(wd, ".search-result__result-link")
		//
		//// filter profile url
		//var selection []selenium.WebElement
		//for i := 0; i < len(profileURL); i += 2 {
		//	selection = append(selection, profileURL[i])
		//}
		//
		//profileURLText := wbAttrToStringSlice(selection, "href")
		//SlicePrint(profileURLText)

		description := findsRetry(wd, ".subline-level-1", "Description")
		descText := wbToStringSlice(description)
		WriteFile("./data/data_"+comp+id, descText)

		//location := findsRetry(wd, ".subline-level-2")
		//locText := wbToStringSlice(location)
		//SlicePrint(locText)
	}
	wd.Close()
	wd.Quit()

}

// Start setup and start the main process
func Start(comp string) {

	threadNumber := 4

	// Initial Webdriver
	wd := initService(port)
	// captcha checking -> if captcha
	signIn(wd)

	captcha := captchaCheck(wd)
	if captcha == true {
		wd.Close()
		wd.Quit()
		os.Exit(1)
	}

	encodedURL := searchFilteredPage(wd, comp)

	// IMPLEMENT SESSIONS SPLITTING
	// Go to filtered search page
	if err := wd.Get(encodedURL); err != nil {
		panic(err)
	}

	//// MAIN CRAWLING
	scroll(wd, 2)
	var lenPage int = lenPage(wd)
	fmt.Printf("\nThere is %d page to crawl !", lenPage)

	page := calcSplitting(lenPage, threadNumber)
	fmt.Printf("\nCalulating page range")
	for i, v := range page {
		fmt.Printf("\nworker_%d extract range : %v", i, v)
	}

	// Spawn ThreadNumber of worker
	var workers []selenium.WebDriver
	workers = CreateWorker(wd, threadNumber, port)
	for i, wd := range workers {
		fmt.Printf("\nSession ID for worker %d : %v", i, wd.SessionID())
	}
	var wg sync.WaitGroup

	for i := 0; i < threadNumber; i++ {
		wg.Add(1)
		go populateWorker(&wg, workers[i], "worker_"+strconv.Itoa(i), lenPage, encodedURL, page[i][0], page[i][1], comp)
	}

	fmt.Println("\nMain: Waiting for workers to finish")
	wg.Wait()
	fmt.Println("\nMain: Completed")

}
