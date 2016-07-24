package starcitizen

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Organization struct {
	SID string

	NumMembers int
	Members    []*OrgMember
}

// RetrieveOrganization gets an RSI organization's profile from the
// robertsspaceindustries.com website.
func RetrieveOrganization(client *http.Client, spectrumID string) (*Organization, error) {
	const profileURLFormat = "https://robertsspaceindustries.com/orgs/%s"

	profileURL := fmt.Sprintf(profileURLFormat, spectrumID)

	resp, err := client.Get(profileURL)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, ErrMissing
	}

	document, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}

	organization := &Organization{}

	organization.SID = getSpectrumID(document)

	return organization, nil
}

func getSpectrumID(doc *goquery.Document) string {
	return doc.Find("h1>.symbol").Text()
}
