package main

import "fmt"

func main() {
	err := scrap()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Scrapping complete")
}
