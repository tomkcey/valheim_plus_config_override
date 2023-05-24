package postprocessor

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	preprocessor "example.com/tomkcey/m/preprocessor"
	utils "example.com/tomkcey/m/utils"
)

func PostProcess(m preprocessor.MapSectionToPairs) {
	t := strconv.FormatInt(time.Now().Unix(), 10)
	fn := strings.Join([]string{t, ".cfg"}, "")
	v := strings.Join([]string{os.TempDir(), fn}, string(os.PathSeparator))
	f, err := os.Create(v)
	utils.Check(err)

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
