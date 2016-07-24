package starcitizen

import (
	"net/http"
	"testing"
)

var orgTests = []struct {
	SID         string
	ExpectedOrg *Organization
}{
	{"sun", &Organization{"SUN", 0, nil}},
}

func TestOrgs(t *testing.T) {
	if testing.Short() {
		t.Skip("not running in short mode")
	}

	for _, test := range orgTests {
		t.Logf("Retrieving organization: %s", test.SID)

		profile, err := RetrieveOrganization(http.DefaultClient, test.SID)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("%#v", profile)

		if profile.SID != test.ExpectedOrg.SID {
			t.Errorf("Spectrum ID doesn't match: expected %q, got %q", test.ExpectedOrg.SID, profile.SID)
		}
	}
}

func TestOrgMembers(t *testing.T) {
	if testing.Short() {
		t.Skip("not running in short mode")
	}

	for _, test := range orgTests {
		t.Logf("Retrieving organization members: %s", test.SID)

		members, err := RetrieveOrganizationMembers(http.DefaultClient, test.SID)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("Got %d members", len(members))
		for _, member := range members {
			t.Logf("%#v", member)
		}
	}
}

func TestOrgMissing(t *testing.T) {
	if testing.Short() {
		t.Skip("not running in short mode")
	}

	_, err := RetrieveOrganization(http.DefaultClient, "aishdf98ahysdf")
	if err == ErrMissing {
		return
	} else if err != nil {
		t.Fatal(err)
	}

	t.Fatal("should not have gotten here")
}
