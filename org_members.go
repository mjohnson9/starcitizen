package starcitizen

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type OrgMember struct {
	Handle string
}

func RetrieveOrganizationMembers(client *http.Client, spectrumID string) ([]*OrgMember, error) {
	document, err := getOrgMembersDoc(client, spectrumID, 1)
	if err != nil {
		return nil, err
	}

	numMembers := getOrgNumMembers(document)
	if numMembers == -1 {
		return nil, ErrSpideringFailed
	}

	numPages := int(numMembers/32 + 1)

	orgMembers := make([]*OrgMember, 0, numMembers)

	pageMembers := getOrgMembersFromDoc(document)
	if pageMembers == nil {
		return nil, ErrSpideringFailed
	}
	orgMembers = append(orgMembers, pageMembers...)

	for i := 2; i <= numPages; i++ {
		pageMembers, err = getOrgMembers(client, spectrumID, i)
		if err != nil {
			return nil, err
		}

		orgMembers = append(orgMembers, pageMembers...)
	}

	return orgMembers, nil
}

func getOrgNumMembers(doc *goquery.Document) int64 {
	totalRows := doc.Find(".totalrows").Text()

	const totalRowsSuffix = " members"
	if strings.HasSuffix(totalRows, totalRowsSuffix) {
		totalRows = totalRows[:len(totalRows)-len(totalRowsSuffix)]
	}

	num, err := strconv.ParseInt(totalRows, 10, 64)
	if err != nil {
		return -1
	}

	return num
}

func getOrgMembersDoc(client *http.Client, sid string, page int) (*goquery.Document, error) {
	const orgMembersURLFormat = "https://robertsspaceindustries.com/orgs/%s/members?page=%d"

	orgMembersURL := fmt.Sprintf(orgMembersURLFormat, sid, page)

	resp, err := client.Get(orgMembersURL)
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

	return document, nil
}

func getOrgMembers(client *http.Client, sid string, page int) ([]*OrgMember, error) {
	document, err := getOrgMembersDoc(client, sid, page)
	if err != nil {
		return nil, err
	}

	pageMembers := getOrgMembersFromDoc(document)
	if pageMembers == nil {
		return nil, ErrSpideringFailed
	}
	return pageMembers, nil
}

func getOrgMembersFromDoc(doc *goquery.Document) []*OrgMember {
	members := make([]*OrgMember, 0, 32)

	memberElements := doc.Find(".member-item")
	memberElements.Each(func(_ int, memberElement *goquery.Selection) {
		if memberElement.Is(".org-visibility-V") {
			handle := memberElement.Find(".nick").Text()
			if len(handle) > 0 {
				members = append(members, &OrgMember{handle})
			}
		}
	})

	return members
}
