package main

import "pages" //generated package

func main(){
	pages.Run(":9999") //start server on localhost:9999
//	pages.Render()
}

/*
package main

import (
	"pages"
	"http"
)

func main(){
               // http.Handle("/manual", http.HandlerFunc(pages.Rendersrchelloghtml));
		http.Handle("/manual", http.HandlerFunc(pages.Render));
                pages.Run(":9999")
}
*/
