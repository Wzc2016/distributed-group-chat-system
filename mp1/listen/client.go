package main

import (
"fmt"
"os"
// "strings"
)

func main() {
	checkArg()
    name := os.Args[1]
    port := os.Args[2]
    // numMem := os.Args[3]

    messageChan := make(chan []string)

    hosts := []string{"sp19-cs425-g08-01.cs.illinois.edu", "sp19-cs425-g08-02.cs.illinois.edu", "sp19-cs425-g08-03.cs.illinois.edu", "sp19-cs425-g08-04.cs.illinois.edu", "sp19-cs425-g08-05.cs.illinois.edu", "sp19-cs425-g08-06.cs.illinois.edu", "sp19-cs425-g08-07.cs.illinois.edu", "sp19-cs425-g08-08.cs.illinois.edu", "sp19-cs425-g08-09.cs.illinois.edu", "sp19-cs425-g08-10.cs.illinois.edu"}
    targets := make([]string, len(hosts))
    for i, host := range hosts {
    	result := getdns(host)
    	targets[i] = string(result) + ":" + port
    	fmt.Println("Added " + host + "(" + targets[i] + ") to list of machines.")
    }

    handleMsgAll(targets, name, messageChan)

    for i := 0; i < len(targets); i ++ {
    	// message := <-messageChan
    	//TODO
    }
}