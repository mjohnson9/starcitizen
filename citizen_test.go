package starcitizen

import (
	"net/http"
	"testing"
)

var citizenTests = []struct {
	Handle          string
	ExpectedCitizen *Citizen
}{
	{"ignis_vulpes", &Citizen{8857, "Ignis_Vulpes", "Fox", nil}},
	{"serrath", &Citizen{34783, "serrath", "serrath", nil}},
	{"the-derp-king", &Citizen{75493, "the-derp-king", "NightExcessive", nil}},
	{"rotorax", &Citizen{105376, "rotorax", "rotorax", nil}},
	{"illuvian", &Citizen{115918, "Illuvian", "Illuvian", nil}},
}

func TestCitizenRetrieves(t *testing.T) {
	if testing.Short() {
		t.Skip("not running in short mode")
	}

	for _, test := range citizenTests {
		profile, err := RetrieveCitizen(http.DefaultClient, test.Handle)
		if err != nil {
			t.Fatal(err)
		}

		if profile.UEENumber != test.ExpectedCitizen.UEENumber {
			t.Errorf("UEE citizen number doesn't match: expected %d, got %d", test.ExpectedCitizen.UEENumber, profile.UEENumber)
		}

		if profile.Handle != test.ExpectedCitizen.Handle {
			t.Errorf("Handle doesn't match: expected %q, got %q", test.ExpectedCitizen.Handle, profile.Handle)
		}

		if profile.Moniker != test.ExpectedCitizen.Moniker {
			t.Errorf("Moniker doesn't match: expected %q, got %q", test.ExpectedCitizen.Moniker, profile.Moniker)
		}
	}
}

func TestCitizenMissing(t *testing.T) {
	if testing.Short() {
		t.Skip("not running in short mode")
	}

	_, err := RetrieveCitizen(http.DefaultClient, "aishdf98ahysdf")
	if err == ErrMissing {
		return
	} else if err != nil {
		t.Fatal(err)
	}

	t.Fatal("should not have gotten here")
}
