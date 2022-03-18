package main

import "fmt"

func main() {
	// err := parseSchedules()
	err := scrap()
	if err != nil {
		fmt.Println(err)
	}
}
