package util

import (
	"strings"

	"github.com/nyaruka/phonenumbers"
)

// NormalizeE164 converts input phone to E.164 format using the given defaultRegion (e.g., "ID").
// Returns normalized "+<country><national>" string and ok=false if invalid.
func NormalizeE164(input string, defaultRegion string) (string, bool) {
	s := strings.TrimSpace(input)
	if s == "" {
		return "", false
	}
	num, err := phonenumbers.Parse(s, defaultRegion)
	if err != nil {
		return "", false
	}
	if !phonenumbers.IsValidNumber(num) {
		return "", false
	}
	return phonenumbers.Format(num, phonenumbers.E164), true
}

// CountryCodeForRegion returns the phone country code for a given ISO region.
// Returns code and ok=false if the region is unknown.
func CountryCodeForRegion(region string) (int32, bool) {
	cc := phonenumbers.GetCountryCodeForRegion(region)
	if cc == 0 {
		return 0, false
	}
	return int32(cc), true
}
