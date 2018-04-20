package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/gocarina/gocsv"
	"github.com/yanatan16/itertools"

	"github.com/minger/precinct-match-go/precinct"
)

func main() {
	// unmarshal data files
	sourcedPrecinctsFile, err := os.OpenFile("sourced_precincts.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer sourcedPrecinctsFile.Close()

	sPrecincts := []*precinct.SPrecinct{}

	if err := gocsv.UnmarshalFile(sourcedPrecinctsFile, &sPrecincts); err != nil {
		panic(err)
	}

	sort.Sort(precinct.SPArr(sPrecincts))
	sp := make([]interface{}, len(sPrecincts))
	for i, p := range sPrecincts {
		sp[i] = p
	}

	sIter := itertools.New(sp...)

	type CountyMap map[string][]*precinct.SPrecinct
	countyMap := make(CountyMap)
	reducer1 := func(memo interface{}, element interface{}) interface{} {
		// assignment is not functional
		if countyMap, ok := memo.(CountyMap); ok {
			if sp, ok := element.(*precinct.SPrecinct); ok {
				if spSlice, ok := countyMap[sp.County]; ok {
					countyMap[sp.County] = append(spSlice, sp)
				} else {
					countyMap[sp.County] = []*precinct.SPrecinct{sp}
				}
				return countyMap
			}
		}
		return memo.(CountyMap)
	}

	sourced := itertools.Reduce(sIter, reducer1, countyMap)

	vfPrecinctsFile, err := os.OpenFile("vf_precincts.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer vfPrecinctsFile.Close()

	vfPrecincts := []*precinct.VFPrecinct{}
	if err := gocsv.UnmarshalFile(vfPrecinctsFile, &vfPrecincts); err != nil {
		panic(err)
	}
	sort.Sort(precinct.VFPArr(vfPrecincts))
	vp := make([]interface{}, len(vfPrecincts))
	for i, p := range vfPrecincts {
		vp[i] = p
	}

	vpIter := itertools.New(vp...)

	type Aggregates struct {
		Matched          precinct.VFPArr
		UnMatched        precinct.VFPArr
		SourcedCountyMap CountyMap
		CurrentVP        *precinct.VFPrecinct
	}

	matched := precinct.VFPArr{}
	unmatched := precinct.VFPArr{}
	p := vfPrecincts[0]
	aggregates := Aggregates{
		Matched:          matched,
		UnMatched:        unmatched,
		SourcedCountyMap: countyMap,
		CurrentVP:        p,
	}

	ctyPred := func(p interface{}) bool {
		cty := p.(*precinct.VFPrecinct).VFPrecinctCounty
		if s, ok := sourced.(CountyMap); ok {
			if _, ok := s[cty]; ok {
				return true
			}
		}
		return false
	}

	strongPred := func(vp *precinct.VFPrecinct) itertools.Predicate {
		return func(p interface{}) bool {
			if sp, ok := p.(*precinct.SPrecinct); ok {
				return precinct.StrongMatch(vp, sp)
			}
			return false
		}
	}

	reducer3 := func(memo interface{}, element interface{}) interface{} {
		if aggregates, ok := memo.(Aggregates); ok {
			if sp, ok := element.(*precinct.SPrecinct); ok {
				vp := aggregates.CurrentVP
				if vp.VFPrecinctCounty == "" {
					return aggregates
				}
				matched := append(aggregates.Matched, vp)
				s := aggregates.SourcedCountyMap[vp.VFPrecinctCounty]
				loc := len(s)
				for i, pc := range s {
					if pc == sp {
						loc = i
						break
					}
				}
				o := aggregates.SourcedCountyMap[p.VFPrecinctCounty]
				if loc < len(s) {
					o = append(s[:loc], s[loc+1:]...)
				}
				aggregates.SourcedCountyMap[p.VFPrecinctCounty] = o

				return Aggregates{
					Matched:          matched,
					UnMatched:        aggregates.UnMatched,
					SourcedCountyMap: aggregates.SourcedCountyMap,
					CurrentVP:        vp,
				}
			}
		}
		return memo.(Aggregates)
	}

	reducer2 := func(memo interface{}, element interface{}) interface{} {
		if aggregates, ok := memo.(Aggregates); ok {
			if vp, ok := element.(*precinct.VFPrecinct); ok {
				scp := aggregates.SourcedCountyMap[vp.VFPrecinctCounty]
				ip := make([]interface{}, len(scp))
				for i, p := range scp {
					ip[i] = p
				}
				a := Aggregates{
					Matched:          aggregates.Matched,
					UnMatched:        aggregates.UnMatched,
					SourcedCountyMap: aggregates.SourcedCountyMap,
					CurrentVP:        vp,
				}
				spIter := itertools.New(ip...)
				return itertools.Reduce(itertools.Filter(strongPred(vp), spIter), reducer3, a)
			}
			return aggregates
		}
		return memo.(Aggregates)
	}

	a := itertools.Reduce(itertools.Filter(ctyPred, vpIter), reducer2, aggregates)
	fmt.Println(len(a.(Aggregates).Matched))
}
