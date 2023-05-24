package preprocessor

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	input "example.com/tomkcey/m/input"
	utils "example.com/tomkcey/m/utils"
)

type Pair = [2]string
type MapSectionToPairs = map[string][]Pair

type SectionToPair struct {
	key   string
	value *MapSectionToPairs
}

type MappedFileStore struct {
	Source *MapSectionToPairs
	Target *MapSectionToPairs
}

func pullSection(l string) *string {
	rgxp := regexp.MustCompile(`[\[\]]+`)
	x := rgxp.Split(l, -1)
	if len(x) == 3 {
		return &x[1]
	}
	return nil
}

func pullComment(l string) *string {
	rgxp := regexp.MustCompile(`[;]+`)
	x := rgxp.Split(l, -1)
	if len(x) == 2 {
		return &x[1]
	}
	return nil
}

func pullKeyValue(l string) *Pair {
	rgxp := regexp.MustCompile(`[=]+`)
	x := rgxp.Split(l, -1)
	if len(x) == 2 {
		k := strings.TrimSpace(x[0])
		v := strings.TrimSpace(x[1])
		return &Pair{k, v}
	}
	return nil
}

func mapFileToIterable(path string) *MapSectionToPairs {
	f, err := os.Open(path)
	utils.Check(err)
	r := bufio.NewReader(f)

	m := make(map[string]([]Pair))

	cs := ""
	last := false
	l, err := r.ReadString('\n')
	for err == nil || !last {
		if err != nil && err == io.EOF {
			last = true
		}

		st := pullSection(l)
		if st != nil {
			cs = *st
		} else {
			cmt := pullComment(l)
			if cmt == nil {
				kv := pullKeyValue(l)
				if kv != nil {
					curKv := m[cs]
					newKv := append(curKv, *kv)
					m[cs] = newKv
				}
			}
		}
		l, err = r.ReadString('\n')
	}

	f.Close()
	return &m
}

func assignIterableToKey(k string, p string, c chan SectionToPair) {
	r := mapFileToIterable(p)
	if r == nil {
		fmt.Println("A problem occured.")
		os.Exit(20)
	}
	// question, would the line below be passed by value or reference?
	c <- SectionToPair{key: k, value: r}
}

func PreProcess(s *input.FileStore) MappedFileStore {
	c := make(chan SectionToPair)

	for _, p := range []Pair{{"source", s.Source}, {"target", s.Target}} {
		go assignIterableToKey(p[0], p[1], c)
	}

	ms := MappedFileStore{}

	for ms.Source == nil || ms.Target == nil {
		r := <-c

		if r.key == "source" {
			ms.Source = r.value
		} else if r.key == "target" {
			ms.Target = r.value
		}
	}

	return ms
}
