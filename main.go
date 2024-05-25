package main

import "fmt"

func main() {

	miurl := "https://www.google.com/search?oq=%D8%AC%D9%86%D9%8A%D9%86+%D8%A8%D9%88%D9%86%D8%B3+%D8%A7%D9%8A%D8%B1%D8%B3&q=%D8%AC%D9%86%D9%8A%D9%86%20%D8%A8%D9%88%D9%86%D8%B3%20%D8%A7%D9%8A%D8%B1%D8%B3&start=1&num=10"
	getData(miurl)
	fmt.Println()
	miurl = "https://www.google.com/search?q=pablo+petrecca&start=11&num=10"
	getData(miurl)
	miurl = "https://www.google.com/search?q=pablo+petrecca&start=21&num=10"
	getData(miurl)

}
