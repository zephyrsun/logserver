package main

import (
	"fmt"
	"time"
)

func main() {

	//v, _ := os.OpenFile("config.json", os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0666))

	//go ticker()

	//time.Sleep(10 * time.Second)
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
