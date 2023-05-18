package input

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type InputHandler struct {
	r *bufio.Reader
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
		panic("Error: Couldn't find the file, make sure to copy the exact path (absolute path) to file")
	}
	fmt.Println("Found file at provided path", info.Name(), strings.Join([]string{"( ",  strconv.FormatInt(info.Size(), 10),  " bytes", " )"}, ""))
}

func (i InputHandler) source() string {
	i.prompt()
	path := i.read()
	i.eval(path)
	return path
}

type FileStore struct {
	Target string
	Source string
}

func Prepare() *FileStore  {
	i := initStdInReader()
	fmt.Println("Taking input for source file (file to be overriden)")
	source := i.source()
	fmt.Println("Taking input for target file (file to override with)")
	target := i.source()
	return &FileStore{Source:source, Target:target}
}