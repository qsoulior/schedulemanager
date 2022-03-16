package main

import "fmt"

func main() {
	err := parseSchedules()
	if err != nil {
		fmt.Println(err)
	}
}
