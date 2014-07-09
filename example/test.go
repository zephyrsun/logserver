package main

import (
	"fmt"
	"time"
	"path/filepath"
)

func fibonacci(c, quit chan int) {
	x, y := 1, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

func main() {

	//err:= os.Mkdir("./log/a.log", os.FileMode(0666))
	//exec.LookPath("")

	e := filepath.Dir("./log/a.log")

	fmt.Printf("%v,%v", e)


}

func ticker() {
	c := time.Tick(2 * time.Second)
	i := 0
	for now := range c {
		fmt.Println(now)
		i++
		if i > 10 {
			break
		}
	}
}

func test2() {
	data := make([]int, 10, 20)
	data[0] = 1
	data[1] = 2
	dataappend := make([]int, 10, 20)//len <=10 则 	result[0] = 99 会 影响源Slice
	dataappend[0] = 1
	dataappend[1] = 2
	result := append(data, dataappend...)
	result[0] = 99
	result[11] = 98

	result = append(result, []int{44, 33}...)
	fmt.Println("length:", len(data), ":", data)
	fmt.Println("length:", len(result), ":", result)
	fmt.Println("length:", len(dataappend), ":", dataappend)
}
