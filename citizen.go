package starcitizen

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Citizen represents an RSI citizen profile
type Citizen struct {
	UEENumber int64
	Handle    string
	Moniker   string

	Organizations []*OrgMembership
}

var rsiBaseURL = &url.URL{Scheme: "https", Host: "robertsspaceindustries.com"}

// RetrieveCitizen gets an RSI citizen's profile from the robertsspaceindustries.com
// website.
func RetrieveCitizen(client *http.Client, profileHandle string) (*Citizen, error) {
	const profileURLFormat = "https://robertsspaceindustries.com/citizens/%s"

	profileURL := fmt.Sprintf(profileURLFormat, profileHandle)

	resp, err := client.Get(profileURL)
	if err != nil {
		return nil, err
	}

	document, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}

	citizen := &Citizen{}

	citizen.UEENumber = getCitizenNumber(document)
	if citizen.UEENumber == -1 {
		return nil, ErrSpideringFailed
	}

	citizen.Handle = getCitizenHandle(document)
	if len(citizen.Handle) == 0 {
		return nil, ErrSpideringFailed
	}

	citizen.Moniker = getCitizenMoniker(document)
	if len(citizen.Moniker) == 0 {
		return nil, ErrSpideringFailed
	}

	citizen.Organizations, err = getOrgMemberships(client, profileHandle)
	if err != nil {
		return nil, err
	}

	return citizen, nil
}

// getCitizenNumber gets the citizen's number from their dossier. It returns
// -1 in case of error.
func getCitizenNumber(doc *goquery.Document) int64 {
	numText := doc.Find(".entry.citizen-record>.value").Text()
	if len(numText) == 0 {
		return -1
	}

	if strings.HasPrefix(numText, "#") {
		numText = numText[1:]
	}

	num, err := strconv.ParseInt(numText, 10, 64)
	if err != nil {
		return -1
	}

	return num
}

// getCitizenHandle gets the citizen's handle from their dossier. It returns
// an empty string in case of error.
func getCitizenHandle(doc *goquery.Document) string {
	const handlePrefix = "https://robertsspaceindustries.com/citizens/"

	overviewLink := doc.Find("a.holotab.active").AttrOr("href", "")
	if len(overviewLink) == 0 {
		return ""
	}

	overviewURL, err := rsiBaseURL.Parse(overviewLink)
	if err != nil {
		return ""
	}

	overviewURLText := overviewURL.String()

	if !strings.HasPrefix(overviewURLText, handlePrefix) {
		return ""
	}

	return overviewURLText[len(handlePrefix):]
}

// getCitizenMoniker gets the citizen's moniker from their dossier (in the
// correct case). It returns an empty string in case of error.
func getCitizenMoniker(doc *goquery.Document) string {
	return doc.Find(".thumb+.info>.entry>.value").Slice(0, 1).Text()
}
