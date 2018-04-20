package precinct

import (
	"regexp"
	"strings"
)

//SPrecinct for source file precincts
type SPrecinct struct {
	PrecinctID         string `csv:"precinct_id"`
	County             string `csv:"county"`
	Ward               string `csv:"ward"`
	PrecinctNumber     string `csv:"precinct_number"`
	PrecinctName       string `csv:"precinct_name"`
	PollingLocationIDS string `csv:"polling_location_ids"`
	Source             string `csv:"source"`
	InternalNotes      string `csv:"INTERNAL_notes"`
}

//VFPrecinct for voter file precincts
type VFPrecinct struct {
	VFPrecinctID     string `csv:"vf_precinct_id"`
	VFPrecinctCounty string `csv:"vf_precinct_county"`
	VFPrecinctWard   string `csv:"vf_precinct_ward"`
	VFPrecinctName   string `csv:"vf_precinct_name"`
	VFPrecinctCode   string `csv:"vf_precinct_code"`
	VFPrecinctCount  int32  `csv:"vf_precinct_count"`
}

// StrongMatch ...
func StrongMatch(vp *VFPrecinct, sp *SPrecinct) bool {
	if vp.VFPrecinctCode == "" || sp.PrecinctNumber == "" {
		return false
	}
	if vp.VFPrecinctCode == sp.PrecinctNumber ||
		strings.TrimLeft(vp.VFPrecinctCode, "0") == strings.TrimLeft(sp.PrecinctNumber, "0") {
		return true
	}
	if vp.VFPrecinctName == "" || sp.PrecinctName == "" {
		return false
	}
	if strings.ToLower(vp.VFPrecinctName) == strings.ToLower(sp.PrecinctName) {
		return true
	}
	return false
}

// WeakMatch ...
func WeakMatch(vp *VFPrecinct, sp *SPrecinct) bool {
	r, _ := regexp.Compile("[()]")
	shortName := r.ReplaceAllString(sp.PrecinctName, "")
	longName := r.ReplaceAllString(vp.VFPrecinctName, "")
	if shortName == "" || longName == "" {
		return false
	}
	s, _ := regexp.Compile("[a-zA-Z]")
	found := s.FindString(shortName)
	if found == "" {
		return false
	}
	if strings.Contains(strings.ToLower(longName), strings.ToLower(shortName)) {
		return true
	}
	return false
}

// VFPArr array of VFPrecinct
type VFPArr []*VFPrecinct

// SPArr array of VFPrecinct
type SPArr []*SPrecinct

// Len ...
func (p VFPArr) Len() int {
	return len(p)
}

// Swap ...
func (p VFPArr) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Less ...
func (p VFPArr) Less(i, j int) bool {
	return p[i].VFPrecinctCounty < p[j].VFPrecinctCounty
}

// Len ...
func (p SPArr) Len() int {
	return len(p)
}

// Swap ...
func (p SPArr) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Less ...
func (p SPArr) Less(i, j int) bool {
	return p[i].County < p[j].County
}
