package randomizer

import (
	"bufio"
	"io"
	"math/rand"
	"os"
	"strings"
)

type (
	Randomizer struct {
		names    []string
		surnames []string
		streets  []string
	}
)

func New() *Randomizer {
	return &Randomizer{
		names:    readFileList("randomizer/names"),
		surnames: readFileList("randomizer/surnames"),
		streets:  readFileList("randomizer/streets"),
	}
}

func readFileList(name string) []string {
	var list = make([]string, 0, 100)
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		lineStr := strings.TrimSpace(string(line))
		if lineStr != "" {
			list = append(list, lineStr)
		}
	}
	return list
}

func (r *Randomizer) RandomName() string {
	return r.names[rand.Intn(len(r.names))]
}

func (r *Randomizer) RandomSurname() string {
	return r.surnames[rand.Intn(len(r.surnames))]
}

func (r *Randomizer) RandomStreet() string {
	return r.streets[rand.Intn(len(r.streets))]
}
