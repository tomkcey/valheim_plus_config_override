package m

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	utils "example.com/tomkcey/m/utils"
)

type InputHandler struct {
	r *bufio.Reader
}

type FileStore struct {
	Source string
	Target string
}

func (i InputHandler) read() string {
	text, err := i.r.ReadString('\n')
	utils.Check(err)
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
