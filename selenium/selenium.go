/*
Package selenium implements linkedin data scraping with selenium automatisation.
*/
package selenium

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tebeka/selenium"
)

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

	fmt.Print(firstCompanyLink)

	/* SOLUTION 2 - INPUT 2: Navigate directly to the companies page; The url need to be input by the user */
	//if err := wd.Get(compURL); err != nil {
	//	panic(err)
	//}

	// Click on "See all X Employees on Linkedin"
	employees, err := wd.FindElement(selenium.ByCSSSelector, ".ember-view.link-without-visited-state.inline-block")
	if err != nil {
		panic(err)
	}

	if err := employees.Click(); err != nil {
		panic(err)
	}

	//// Search for company name
	//searchBox, err := wd.FindElement(selenium.ByCSSSelector, ".search-global-typeahead__input")
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = searchBox.SendKeys("atos")
	//if err != nil {
	//	panic(err)
	//}
	//searchButton, err := wd.FindElement(selenium.ByCSSSelector, ".search-global-typeahead__button")
	//if err != nil {
	//	panic(err)
	//}
	//if err := searchButton.Click(); err != nil {
	//	panic(err)
	//}
	//
	//// Click on more button to display filter options
	//moreButton, err := wd.FindElement(selenium.ByCSSSelector, ".artdeco-dropdown__trigger")
	//if err != nil {
	//	panic(err)
	//}
	//if err := moreButton.Click(); err != nil {
	//	panic(err)
	//}
	//
	//// Click on "companies" to filter result by companies

	// return empty for testing
	return make([]string, 1)
}

/* LinkedinUsers
Return list of users related to a company
Step to perform :
1 - Go to Linkedin.com
2 - Search for the company name
3 - Tap on the company page
4 - Tap on ~"See Peeople working here"
*/
func LinkedinUsers(comp string) {
	start("atos")
}
