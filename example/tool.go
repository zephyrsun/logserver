package main

import (
	ls "github.com/zephyrsun/logserver"
	"os"
	"flag"
	"bufio"
	"fmt"
)

func main() {

	file := flag.String("f", "", "name of check file")
	flag.Parse()

	if *file != "" {
		checkFile(*file)
	}

}

func checkFile(name string) {
	f, err := os.OpenFile(name, os.O_RDONLY, 0655)
	ls.DumpError(err, true)

	r := bufio.NewReader(f)

	n := 0
	i := 0
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			break
		}
		i++

		l := len(line)

		if n == 0 {
			n = l
		}else if n != l {
			fmt.Printf("not match %d", n)
		}
	}

	fmt.Printf("total line %d", i)
}
