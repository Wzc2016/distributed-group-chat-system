package main

import(  
    "fmt"  
    "os"  
    "net"
    // "bufio"
    // "strings"
)

// a map to store the memberships in the chatroom
var memMap map[string]string


func main() {
    // To check the input
    checkArg()
    // To get the name, port and num of people from input
    name := os.Args[1]
    port := os.Args[2]
    // numMem := os.Args[3]
    // To initialize the dns array of all the vm
    hosts := []string{"sp19-cs425-g08-01.cs.illinois.edu", "sp19-cs425-g08-02.cs.illinois.edu", "sp19-cs425-g08-03.cs.illinois.edu", "sp19-cs425-g08-04.cs.illinois.edu", "sp19-cs425-g08-05.cs.illinois.edu", "sp19-cs425-g08-06.cs.illinois.edu", "sp19-cs425-g08-07.cs.illinois.edu", "sp19-cs425-g08-08.cs.illinois.edu", "sp19-cs425-g08-09.cs.illinois.edu", "sp19-cs425-g08-10.cs.illinois.edu"}
    targets := make([]string, len(hosts))
    for i, host := range hosts {
        result := getdns(host)
        targets[i] = string(result) + ":" + port
    }
    // begin listen to the port
    msgChan := make(chan string, 10)
    go listener(port, msgChan)
    // begin sending message
    messageChan := make(chan []string)
    
    // return the machines that are alive
    go checkConnectAll(targets, name, messageChan)
    
    suc := <-messageChan
    
    for _, item := range suc {

        checkItem(item)

        conn, err := net.Dial("tcp", item)
        // handle the error 
        if err != nil {
        fmt.Println("checkConnect error")
        os.Exit(1)
    }

    }

    if err != nil {
        fmt.Println("服务没打开")
        os.Exit(1)
    }
    for {
        var msg string
        fmt.Scan(&msg)
        fmt.Print("<" + name + ">" + "说:")
        fmt.Println(msg)
        b := []byte("<" + name + ">" + "说：" + msg)
        conn.Write(b)

        select {
            case i := <- msgChan :
                fmt.Println(i)
        }
    }
}

func listener(port string, msgChan chan string){
    server, err := net.Listen("tcp", ":" + port)
    errHandler(err, "Can not start server!", true)
    conn, err := server.Accept()
        errHandler(err, "Can not open connection!", true)
    for {
        go func(conn net.Conn) {
            defer conn.Close()
            for {
                buff := make([]byte, 256)
                length, err := conn.Read(buff)
                if err != nil {
                    break
                }
                msgChan <- string(buff[:length])
            }
        } (conn)
    }
}


func checkConnectAll(targets []string, name string, msgChan chan []string) {
        results := make([]string, 0)
        for _, target := range targets {
            conChan := make(chan string)
            go checkConnect(target, name, conChan)
            i := <- conChan
            append(results, i)
        }
        msgChan <- results
}

// TODO add timestamps
func checkConnect(target string, name string, msgChan chan string) {
            // fmt.Println("hi")
            results := make([]string, 1)
            results[0] = name
            conn, err := net.Dial("tcp", target)
            // no err
            if err != nil {
                return
            }
            msgChan <- append(results, target)
            return
}

// TODO1: Check whether items in the addrList are in the list (or map), if not, add them into the map.
// TODO2: Check whether there are items not in the list but in the map, put them into an array and send the array into disconnects
func checkItem(addrList []string, memMap map[string]string, disconnects chan []string) {

}