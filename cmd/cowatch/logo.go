package main

import (
	"fmt"

	codes "github.com/avearmin/stylecodes"
)

func printLogo() {
	logo := `
 _   _           _  _|_  _  |    
/   / \ |  |  | / |  |  /   |/\  
\__ \_/  \/ \/  \/|_ |_ \__ |  | 
`

	fmt.Print(codes.ColorBrightMagenta)
	fmt.Print(logo)
	fmt.Println(codes.ResetColor)
}
