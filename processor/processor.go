package processor

import (
	preprocessor "example.com/tomkcey/m/preprocessor"
)

type SectionMap struct {
	section string
	source  []preprocessor.Pair
	target  []preprocessor.Pair
}

type FilteredSectionMap struct {
	section string
	pairs   []preprocessor.Pair
}

func mapSections(ms preprocessor.MappedFileStore) []SectionMap {
	r := make([]SectionMap, 0, 1)
	for k, v := range *ms.Source {
		sm := SectionMap{section: k, source: v}
		t := *ms.Target
		sm.target = t[k]
		r = append(r, sm)
	}
	return r
}

func overrideSection(sm SectionMap) FilteredSectionMap {
	if sm.target == nil || len(sm.target) == 0 {
		return FilteredSectionMap{section: sm.section, pairs: sm.source}
	}

	r := make([]preprocessor.Pair, 0, 1)
	for _, pairA := range sm.source {
		f := false
		for _, pairB := range sm.target {
			if pairB[0] == pairA[0] {
				f = true
				r = append(r, preprocessor.Pair{pairB[0], pairB[1]})
			}
		}
		if !f {
			r = append(r, preprocessor.Pair{pairA[0], pairA[1]})
		}
	}

	// adding new ones that weren't there in source
	for _, pairA := range sm.target {
		f := false
		for _, pairB := range sm.source {
			if pairA[0] == pairB[0] {
				f = true
			}
		}
		if !f {
			r = append(r, preprocessor.Pair{pairA[0], pairA[1]})
		}
	}

	return FilteredSectionMap{section: sm.section, pairs: r}
}

func overrideSections(sms []SectionMap) *preprocessor.MapSectionToPairs {
	r := make(preprocessor.MapSectionToPairs)
	c := make(chan FilteredSectionMap)

	for _, sm := range sms {
		go func(smap SectionMap) {
			c <- overrideSection(smap)
		}(sm)
	}

	for i := 0; i < len(sms); i++ {
		o := <-c
		r[o.section] = o.pairs
	}

	return &r
}

func Process(ms preprocessor.MappedFileStore) *preprocessor.MapSectionToPairs {
	sms := mapSections(ms)
	return overrideSections(sms)
}
