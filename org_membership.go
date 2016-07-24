package starcitizen

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// OrgMembership represents a citizen's membership in an organization.
type OrgMembership struct {
	SID     string
	Name    string
	Rank    string
	RankNum int8

	Visibility string

	Main bool
}

func getOrgMemberships(client *http.Client, profileHandle string) ([]*OrgMembership, error) {
	const profileURLFormat = "https://robertsspaceindustries.com/citizens/%s/organizations"

	profileURL := fmt.Sprintf(profileURLFormat, profileHandle)

	resp, err := client.Get(profileURL)
	if err != nil {
		return nil, err
	}

	document, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}

	orgElements := document.Find(".org")

	orgs := make([]*OrgMembership, 0, orgElements.Size())

	orgElements.Each(func(_ int, orgElement *goquery.Selection) {
		org := convertOrgMembershipElement(orgElement)
		if org == nil {
			panic(ErrSpideringFailed)
		}

		orgs = append(orgs, org)
	})

	return orgs, nil
}

func convertOrgMembershipElement(el *goquery.Selection) *OrgMembership {
	membership := &OrgMembership{}

	membership.Main = el.Is(".main")

	membership.Visibility = getOrgMembershipVisibility(el)
	if membership.Visibility == "" {
		return nil
	}

	if membership.Visibility == "visible" {
		membership.SID = getOrgMembershipSID(el)
		if len(membership.SID) == 0 {
			return nil
		}

		membership.Name = getOrgMembershipName(el)
		if len(membership.Name) == 0 {
			return nil
		}

		membership.Rank, membership.RankNum = getOrgMembershipRank(el)
		if membership.Rank == "" || membership.RankNum == -1 {
			return nil
		}
	}

	return membership
}

func getOrgMembershipVisibility(el *goquery.Selection) string {
	if el.Is(".visibility-R") {
		return "redacted"
	} else if el.Is(".visibility-V") {
		return "visible"
	} else {
		return ""
	}
}

func getOrgMembershipSID(el *goquery.Selection) string {
	sid, _ := getEntryValue(el, "spectrum identification (sid)")

	return sid
}

func getOrgMembershipName(el *goquery.Selection) string {
	return el.Find(".entry>a.value").Text()
}

func getOrgMembershipRank(el *goquery.Selection) (string, int8) {
	rank, found := getEntryValue(el, "organization rank")
	if !found {
		return "", -1
	}

	rankNum := int8(el.Find(".ranking>.active").Size())

	return rank, rankNum
}
