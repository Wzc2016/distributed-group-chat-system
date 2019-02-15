package main

import (
"bufio"
"os"
"os/exec"
"net"
"fmt"
// "cmd"
"bytes"
)

// check the number of args
func checkArg() {
    if len(os.Args) != 4  {     
        fmt.Println("Usage: " + os.Args[0] + " <name> <port> <n>")  
        os.Exit(0)  
    }
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

func getdns(vm string) []byte{
	cmd := exec.Command("/usr/bin/dig","+short", vm)
	resultsBytes, err := cmd.CombinedOutput()
	resultsBytes = bytes.Trim(resultsBytes, "\n")
	errHandler(err,"Couldn't lookup machine address:" + vm, true )
	return resultsBytes
}

func handleMsgAll(targets []string, name string, msgChan chan []string) {
	for _, target := range targets {
		handleMsg(target, name, msgChan)
	}
}

// TODO add timestamps
func handleMsg(target string, name string, msgChan chan []string) {
		go func() {
			results := make([]string, 1)
			results[0] = name
			conn, err := net.Dial("tcp", target)
			errHandler(err, "Can not connect to the server: " + target, false)
			// no err
			if err != nil {
				msgChan <- append(results, "has left")
				return
			}

			fmt.Fprintf(conn, name + ":")
			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				results = append(results, scanner.Text())
			}
			err = scanner.Err()
			errHandler(err, "Can not read server response!", true)
			msgChan <- results
		}()
}