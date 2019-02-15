package main

import "fmt"
import "strings"

func main() {
	resultBytes := getdns("sp19-cs425-g08-01.cs.illinois.edu")
	var strBuilder strings.Builder
		strBuilder.Write(resultBytes)
		strBuilder.WriteString(":500")
		result := strBuilder.String()
	fmt.Println(result)
	
}