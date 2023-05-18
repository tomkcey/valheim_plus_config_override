package processor

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
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
		// fmt.Print(k, "=")
		// fmt.Print(v, "\n")
		return &[2]string{k, v}
	}
	return nil
}

func PreProcess(path string) *map[string][][2]string {
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

func Process(sm map[string][][2]string, tm map[string][][2]string) map[string][][2]string {
	m := make(map[string][][2]string)
	for k := range sm {
		smKvs := sm[k]
		tmKvs := tm[k]
		tmpKvM := make(map[string]string)
		for i := 0; i < len(smKvs); i++ {
			t := smKvs[i]
			tmpKvM[t[0]] = t[1]
		}
		for j := 0; j < len(tmKvs); j++ {
			t := tmKvs[j]
			tmpKvM[t[0]] = t[1]
		}
		arr := make([][2]string, 1)
		for key := range tmpKvM {
			arr = append(arr, [2]string{key, tmpKvM[key]})
		}
		m[k] = arr
	}
	return m
}

func PostProcess(m map[string][][2]string) {
	v := strings.Join([]string{os.TempDir(), "valheim_plus.cfg"}, "/")
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
			if(i == len(m[k]) - 1) {
				w.WriteString("\n")
			}

		}
	}

	w.Flush()

	f.Close()
}
