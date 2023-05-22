package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Pair = [2]string
type MapSectionToPairs = map[string][]Pair

type InputHandler struct {
	r *bufio.Reader
}

type FileStore struct {
	Source string
	Target string
}

type MappedFileStore struct {
	Source *MapSectionToPairs
	Target *MapSectionToPairs
}
type SectionToPair struct {
	key   string
	value *MapSectionToPairs
}

type SectionMap struct {
	section string
	source  []Pair
	target  []Pair
}

type FilteredSectionMap struct {
	section string
	pairs   []Pair
}

func main() {
	store := Prepare()
	if store == nil {
		fmt.Println("A problem occured.")
		os.Exit(20)
	}

	now := time.Now()

	ms := PreProcess(store)
	m := Process(ms)

	PostProcess(*m)

	fmt.Println(time.Since(now))
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
	check(err)
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

func PreProcess(s *FileStore) MappedFileStore {
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

func mapSections(ms MappedFileStore) []SectionMap {
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

	r := make([]Pair, 0, 1)
	for _, pairA := range sm.source {
		f := false
		for _, pairB := range sm.target {
			if pairB[0] == pairA[0] {
				f = true
				r = append(r, Pair{pairB[0], pairB[1]})
			}
		}
		if f == false {
			r = append(r, Pair{pairA[0], pairA[1]})
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
		if f == false {
			r = append(r, Pair{pairA[0], pairA[1]})
		}
	}

	return FilteredSectionMap{section: sm.section, pairs: r}
}

func overrideSections(sms []SectionMap) *MapSectionToPairs {
	r := make(MapSectionToPairs)
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

func Process(ms MappedFileStore) *MapSectionToPairs {
	sms := mapSections(ms)
	return overrideSections(sms)
}

func PostProcess(m MapSectionToPairs) {
	t := strconv.FormatInt(time.Now().Unix(), 10)
	fn := strings.Join([]string{t, ".cfg"}, "")
	v := strings.Join([]string{os.TempDir(), fn}, string(os.PathSeparator))
	f, err := os.Create(v)
	check(err)

	w := bufio.NewWriter(f)

	for k, v := range m {
		s := strings.Join([]string{"[", k, "]", "\n"}, "")
		w.WriteString(s)
		w.WriteString("\n")
		for i := 0; i < len(v); i++ {
			x := v[i]
			n := make([]string, 0)
			for j := 0; j < len(x); j++ {
				fmtstr := strings.TrimSpace(x[j])
				if len(fmtstr) > 0 {
					n = append(n, fmtstr)
				}
			}
			if len(n) == 2 {
				kv := strings.Join([]string{n[0], "=", n[1], "\n"}, "")
				w.WriteString(kv)
			}
			if i == len(v)-1 {
				w.WriteString("\n")
			}

		}
	}

	w.Flush()

	f.Close()

	fmt.Println("Wrote new file at", v)
}

func (i InputHandler) read() string {
	text, err := i.r.ReadString('\n')
	check(err)
	text = strings.Replace(text, "\n", "", -1)
	text = strings.Replace(text, "\r", "", -1)
	return strings.Replace(text, "\\", "/", -1)
}

func (i InputHandler) eval(path string) {
	info, err := os.Stat(path)
	if err != nil {
		fmt.Println("No such file at provided path.")
		fmt.Println("Exiting...")
		os.Exit(20)
	}
	fmt.Println("Found file at provided path", info.Name(), strings.Join([]string{"( ", strconv.FormatInt(info.Size(), 10), " bytes", " )"}, ""))
}

func (i InputHandler) source() string {
	fmt.Print("Input file path: ")
	path := i.read()
	i.eval(path)
	return path
}

func Prepare() *FileStore {
	i := InputHandler{r: bufio.NewReader(os.Stdin)}
	fmt.Println("Taking input for source file (file to be overriden)")
	source := i.source()
	fmt.Println("Taking input for target file (file to override with)")
	target := i.source()
	return &FileStore{Source: source, Target: target}
}
