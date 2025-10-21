package pii

import "regexp"

// Detector provides simple PII detection helpers built upon Go's regexp
// package. Patterns are compiled once and re-used to minimise allocations.
type Detector struct {
	email  *regexp.Regexp
	phone  *regexp.Regexp
	ip     *regexp.Regexp
	credit *regexp.Regexp
}

// New creates a detector with sensible defaults.
func New() *Detector {
	return &Detector{
		email:  regexp.MustCompile(`(?i)[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}`),
		phone:  regexp.MustCompile(`\+?[0-9]{1,3}[\s-]?(?:\([0-9]{1,4}\)[\s-]?)?[0-9\s-]{5,}`),
		ip:     regexp.MustCompile(`\b(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(?:\.|$)){4}`),
		credit: regexp.MustCompile(`\b(?:\d[ -]*?){13,16}\b`),
	}
}

// ContainsPII reports whether any of the detector's patterns match the input.
func (d *Detector) ContainsPII(input string) bool {
	return d.email.MatchString(input) ||
		d.phone.MatchString(input) ||
		d.ip.MatchString(input) ||
		d.credit.MatchString(input)
}
