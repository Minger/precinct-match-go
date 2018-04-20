package precinct_test

import (
	"testing"

	"github.com/minger/precinct-match-go/precinct"
)

func TestWeakMatch(t *testing.T) {
	vp := precinct.Precinct{VFPrecinctName: "boBBy"}
	sp := precinct.Precinct{PrecinctName: "bob"}
	if !precinct.WeakMatch(&vp, &sp) {
		t.Errorf("Expected Source precinct name %s to match the VF precinct name %s", sp.PrecinctName, vp.VFPrecinctName)
	}

	vp = precinct.Precinct{VFPrecinctName: ")bo(BBy"}
	sp = precinct.Precinct{PrecinctName: "(bob"}
	if !precinct.WeakMatch(&vp, &sp) {
		t.Errorf("Expected Source precinct name %s to match the VF precinct name %s despite parens", sp.PrecinctName, vp.VFPrecinctName)
	}

	vp = precinct.Precinct{VFPrecinctName: ")bo(BBy"}
	sp = precinct.Precinct{PrecinctName: "rob"}
	if precinct.WeakMatch(&vp, &sp) {
		t.Errorf("Expected Source precinct name %s to NOT match the VF precinct name %s", sp.PrecinctName, vp.VFPrecinctName)
	}

	vp = precinct.Precinct{VFPrecinctName: ")14("}
	sp = precinct.Precinct{PrecinctName: "14"}
	if precinct.WeakMatch(&vp, &sp) {
		t.Errorf("Expected names without letters to return false - sp: %s, vp: %s", sp.PrecinctName, vp.VFPrecinctName)
	}

	vp = precinct.Precinct{VFPrecinctName: "bob"}
	sp = precinct.Precinct{PrecinctName: ""}
	if precinct.WeakMatch(&vp, &sp) {
		t.Errorf("Expected names with blanks return false - sp: %s, vp: %s", sp.PrecinctName, vp.VFPrecinctName)
	}
}
