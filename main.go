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

type InputHandler struct {
	r *bufio.Reader
}

type FileStore struct {
	Target string
	Source string
}

type MappedFileStore struct {
	Source *map[string][][2]string
	Target *map[string][][2]string
}
type Pair struct {
	key   string
	value *map[string][][2]string
}

func assignIterableToKey(k string, p string, c chan Pair) {
	r := mapFileToIterable(p)
	if r == nil {
		fmt.Println("A problem occured.")
		os.Exit(20)
	}
	// question, would the line below be passed by value or reference?
	c <- Pair{key: k, value: r}
}

func PreProcess(s *FileStore) MappedFileStore {
	c := make(chan Pair)

	for _, p := range [][2]string{{"source", s.Source}, {"target", s.Target}} {
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

func pullKeyValue(l string) *[2]string {
	rgxp := regexp.MustCompile(`[=]+`)
	x := rgxp.Split(l, -1)
	if len(x) == 2 {
		k := strings.TrimSpace(x[0])
		v := strings.TrimSpace(x[1])
		return &[2]string{k, v}
	}
	return nil
}

func mapFileToIterable(path string) *map[string][][2]string {
	f, err := os.Open(path)
	check(err)
	r := bufio.NewReader(f)

	m := make(map[string]([][2]string))

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

type SectionMap struct {
	section string
	source  [][2]string
	target  [][2]string
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

type FilteredSectionMap struct {
	section string
	pairs   [][2]string
}

func overrideSection(sm SectionMap) FilteredSectionMap {
	if sm.target == nil || len(sm.target) == 0 {
		return FilteredSectionMap{section: sm.section, pairs: sm.source}
	}

	r := make([][2]string, 0, 1)
	for _, pairA := range sm.source {
		f := false
		for _, pairB := range sm.target {
			if pairB[0] == pairA[0] {
				f = true
				r = append(r, [2]string{pairB[0], pairB[1]})
			}
		}
		if f == false {
			r = append(r, [2]string{pairA[0], pairA[1]})
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
			r = append(r, [2]string{pairA[0], pairA[1]})
		}
	}

	return FilteredSectionMap{section: sm.section, pairs: r}
}

func overrideSections(sms []SectionMap) *map[string][][2]string {
	r := make(map[string][][2]string)
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

func Process(ms MappedFileStore) *map[string][][2]string {
	sms := mapSections(ms)
	return overrideSections(sms)
}

func PostProcess(m map[string][][2]string) {
	t := strconv.FormatInt(time.Now().Unix(), 10)
	fn := strings.Join([]string{t, ".cfg"}, "")
	v := strings.Join([]string{os.TempDir(), fn}, string(os.PathSeparator))
	f, err := os.Create(v)
	check(err)

	w := bufio.NewWriter(f)

	for k := range m {
		s := strings.Join([]string{"[", k, "]", "\n"}, "")
		w.WriteString(s)
		w.WriteString("\n")
		for i := 0; i < len(m[k]); i++ {
			x := m[k][i]
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
			if i == len(m[k])-1 {
				w.WriteString("\n")
			}

		}
	}

	w.Flush()

	f.Close()

	fmt.Println("Wrote new file at", v)
}

func initStdInReader() InputHandler {
	return InputHandler{r: bufio.NewReader(os.Stdin)}
}

func (i InputHandler) prompt() {
	fmt.Print("Input file path: ")
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
		fmt.Println("A problem occured.")
		fmt.Println("Exiting...")
		os.Exit(20)
	}
	fmt.Println("Found file at provided path", info.Name(), strings.Join([]string{"( ", strconv.FormatInt(info.Size(), 10), " bytes", " )"}, ""))
}

func (i InputHandler) source() string {
	i.prompt()
	path := i.read()
	i.eval(path)
	return path
}

func Prepare() *FileStore {
	i := initStdInReader()
	fmt.Println("Taking input for source file (file to be overriden)")
	source := i.source()
	fmt.Println("Taking input for target file (file to override with)")
	target := i.source()
	return &FileStore{Source: source, Target: target}
}
