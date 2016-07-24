package starcitizen

import (
	"errors"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ErrSpideringFailed is returned whenever spidering a document formatted with
// humans, as opposed to machines, in mind. This is usually caused by a change
// in the format of an HTML page.
var ErrSpideringFailed = errors.New("document spidering failed, probably due to a change in document format")

func getEntryValue(el *goquery.Selection, label string) (string, bool) {
	value := ""
	found := false

	el.Find(".entry>.label").EachWithBreak(func(_ int, entryElement *goquery.Selection) bool {
		text := strings.ToLower(entryElement.Text())
		if text == label {
			found = true
			value = entryElement.NextFiltered(".value").Text()
			return false
		}

		return true
	})

	return value, found
}
