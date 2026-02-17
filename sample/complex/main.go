package main

import "fmt"

func main() {
	userApp, err := InitializeUserApp()
	if err != nil {
		panic(err)
	}

	orderApp, err := InitializeOrderApp()
	if err != nil {
		panic(err)
	}

	fmt.Println(userApp.UserHandler.Handle(1))
	fmt.Println(orderApp.ProductHandler.Handle(100))
	fmt.Println(orderApp.OrderHandler.Handle(1, 100))
}
