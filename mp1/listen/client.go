package main

import (
"fmt"
"os"
// "strings"
"net"
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
    	// fmt.Println("Added " + host + "(" + targets[i] + ") to list of machines.")
    }
    go handleMsgAll(targets, name, messageChan)

    suc := <-messageChan
    // fmt.Println(suc[1])
    conn, err := net.Dial("tcp", suc[1])

	errHandler(err, "无法连接到: " + suc[1], false)
	// no err

	if err != nil {
		fmt.Println("服务没打开")
		os.Exit(1)
	}
	// defer conn.Close()
	for {
		var msg string
		fmt.Scan(&msg)
		fmt.Print("<" + name + ">" + "说:")
		fmt.Println(msg)
		b := []byte("<" + name + ">" + "说：" + msg)
		conn.Write(b)
	}

    va := <-messageChan
    fmt.Println(va[0])
}

// handle the error message and print to the terminal
func errHandler(err error, message string, exit bool) {
    if err != nil {
        fmt.Println(message)
        if exit {
            os.Exit(1)
        }
    }
}