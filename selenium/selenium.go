/*
Package selenium implements linkedin data scraping with selenium automatisation.
*/
package selenium

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/tebeka/selenium"
)

func lenPage(wd selenium.WebDriver) int {
	pageNumber, err := wd.FindElements(selenium.ByCSSSelector, ".artdeco-pagination__indicator.artdeco-pagination__indicator--number")
	if err != nil {
		panic(err)
	}

	lenPage, _ := pageNumber[len(pageNumber)-1].Text()

	conv, err := strconv.Atoi(lenPage)
	if err != nil {
		panic(err)
	}

	return conv

}

/*
Utity function to init Selenium
*/
func start(comp string) []string {

	const (
		// These paths will be different on your system.
		seleniumPath    = "/home/wr3ck3r/go/pkg/mod/github.com/tebeka/selenium@v0.9.9/vendor/selenium-server.jar"
		geckoDriverPath = "/home/wr3ck3r/go/pkg/mod/github.com/tebeka/selenium@v0.9.9/vendor/geckodriver"
		port            = 8080
	)
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
	wd, err := selenium.NewRemote(caps, "http://localhost:8080/wd/hub")

	if err != nil {
		panic(err)
	}

	// Navigate to Linkedin
	if err := wd.Get("https://www.linkedin.com/login?fromSignIn=true&trk=guest_homepage-basic_nav-header-signin"); err != nil {
		panic(err)
	}

	// SIGN-IN

	// Load Credential for Linkedin SignIn
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
	btn, err := wd.FindElement(selenium.ByCSSSelector, ".btn__primary--large")
	if err != nil {
		panic(err)
	}
	if err := btn.Click(); err != nil {
		panic(err)
	}

	// Just display the username (page loading)
	wd.SetImplicitWaitTimeout(time.Second * 2)
	connectedUser, err := wd.FindElement(selenium.ByCSSSelector, ".profile-rail-card__actor-link")
	if err != nil {
		panic(err)
	}

	var output string

	time.Sleep(time.Millisecond * 100)
	wd.SetImplicitWaitTimeout(time.Second * 2)

	output, err = connectedUser.Text()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n %s is connected\n", output)

	/* SOLUTION 1 - INPUT 1 : Navigate to the companies page with just the companies name */
	// Navigate to Linkedin companies search result (for givven company name "comp")

	var searchURL string = "https://www.linkedin.com/search/results/companies/?keywords=" + comp + "&origin=SWITCH_SEARCH_VERTICAL"

	if err := wd.Get(searchURL); err != nil {
		panic(err)
	}

	// click on the first company found in the search result
	time.Sleep(time.Millisecond * 100)
	firstCompanyLink, err := wd.FindElements(selenium.ByCSSSelector, ".app-aware-link.ember-view")
	if err != nil {
		panic(err)
	}

	// click on the second link on the slice. In fact image link can't be clicked so just use the second link
	// TODO Need To be OPTIMIZED -> I get all the "app-aware-link" of the page but i just need the first one
	if err := firstCompanyLink[1].Click(); err != nil {
		panic(err)
	}

	/* SOLUTION 2 - INPUT 2: Navigate directly to the companies page; The url need to be input by the user */
	//if err := wd.Get(compURL); err != nil {
	//	panic(err)
	//}

	// Click on "See all X Employees on Linkedin"
	wd.SetImplicitWaitTimeout(time.Second * 3)
	employees, err := wd.FindElement(selenium.ByCSSSelector, ".ember-view.link-without-visited-state.inline-block")
	if err != nil {
		panic(err)
	}

	if err := employees.Click(); err != nil {
		panic(err)
	}

	// Scroll the page to load the page entirely
	Scroll(wd)

	// TODO -> Add filter selection to select only wanted companies or subsidiary companies
	// get and process actual search url
	// URL structc : https://www.linkedin.com/search/results/people/?facetCurrentCompany=["1259"%2C"2274"%2C"208298"%2C"1260"%2C"53472064"]
	wd.SetImplicitWaitTimeout(time.Second * 3)
	currentURL, err := wd.CurrentURL()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", currentURL)

	encodedURL := DecodeReconstruct(currentURL)
	// reconstruct URL

	if err := wd.Get(encodedURL); err != nil {
		panic(err)
	}

	/* Get and Format data */
	// Loop Through pages (1..n)

	var lenPage int = 1

	for i := 0; i < lenPage; i++ {

		wd.SetImplicitWaitTimeout(time.Second * 3)
		time.Sleep(time.Second * 2)
		users, err := wd.FindElements(selenium.ByCSSSelector, ".actor-name")
		if err != nil {
			panic(err)
		}

		usersText := WbToString(users)
		SlicePrint(usersText)

		// ProfileUrl
		profileURL, err := wd.FindElements(selenium.ByCSSSelector, ".search-result__result-link")
		if err != nil {
			panic(err)
		}

		profileURLText := WbAttrToString(profileURL, "href")
		SlicePrint(profileURLText)

		// Description
		description, err := wd.FindElements(selenium.ByCSSSelector, ".subline-level-1")
		if err != nil {
			panic(err)
		}
		descText := WbToString(description)
		SlicePrint(descText)

		// Location
		location, err := wd.FindElements(selenium.ByCSSSelector, ".subline-level-2")
		if err != nil {
			panic(err)
		}

		locText := WbToString(location)
		SlicePrint(locText)

		// function to click on the next Button
		// nextPage()

	}

	return make([]string, 1)
}

/**/
func LinkedinUsers(comp string) {
	start("sncf")
}
