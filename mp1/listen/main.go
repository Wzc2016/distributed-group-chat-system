package main

import(  
    "fmt"  
    "os"  
    "net"
    // "bufio"
    // "strings"
    "os/exec" 
    "bytes"
)

// a map to store the memberships in the chatroom
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
        targets[i] = result + ":" + port
    }
    server := setupServer(port)

    

    // listen and print to the prompt 
    inMsg := make(chan string)
    go func() {
        for {
        conn, err := server.Accept()
        errHandler(err, "Can not open connection!", true)
        // fmt.Println("成功连接，开始读取")
        go func(conn net.Conn) {
            defer conn.Close()
            for {
                buff := make([]byte, 256)
                length, err := conn.Read(buff)
                if err != nil {
                    break
                }
                inMsg <- string(buff[:length])
            }
        } (conn)
    }
    }()

    var hmsg string
    go func(){
        for {
            select {
            case hmsg = <- inMsg :
                fmt.Println(hmsg)
            default:
            }
        }
        }()
    
    // begin sending message
    // messageChan := make(chan []string)
    // return the machines that are alive
    // for _,target := range targets {
    //     fmt.Println(target)
    // }
    // var suc []string
    memMap := make(map[string]bool)
    go func(){
        for{
            checkConnectAll(targets, name, memMap)
        }
    }()

    // go func() {
    //     for{
    //         if(len(suc) >= 1){
    //             fmt.Println(suc[0])
    //         }
    //     }
    // }()
    // deal with the failures
    
    // memChan := make(chan []string)
    // checkItem(suc, memMap, memChan)
    // leaves := <-memChan
    // for leave := range leaves {
    //     fmt.Println(leave + " has left")
    // }
    
    outMsg := make(chan string, 3)
    var msg string

    // multicast the messages
    go func() {
        for{ 
            select {
            case msg = <- outMsg:
                for item, _ := range memMap {
                	// face, _ := net.InterfaceAddrs()
                    if(len(item) > 0 && item != targets[0]){
                        // fmt.Println(item)
                        conn, err := net.Dial("tcp", item)
                        // handle the error 
                        if err != nil {
                            fmt.Println("checkConnect error")
                            os.Exit(1)
                        }
                        conn.Write([]byte(msg))
                    }
                }
            default:
            }

            
        }
    }()

    for {
        var newMsg string
        // fmt.Printf("%s: ", name)
        fmt.Scan(&newMsg)
        outMsg <- name + ": " + newMsg
    }
}

func setupServer(port string) net.Listener {
    server, err := net.Listen("tcp", ":" + port)
    errHandler(err, "Can not start server!", true)
    return server
}

// func listener(server net.Listener, msgChan chan string){
//     for { 
//         conn, err := server.Accept()
//         // fmt.Println("hi")
//         errHandler(err, "Can not open connection!", true)
//         go func(conn net.Conn) {
//             defer conn.Close()
//             for {
//                 buff := make([]byte, 256)
//                 length, err := conn.Read(buff)
//                 if err != nil {
//                     break
//                 }
//                 msgChan <- string(buff[:length])
//             }
//         } (conn)
//     }
// }


func checkConnectAll(targets []string, name string, memMap map[string]bool) {
        // var results []string
        for _, target := range targets {
            // fmt.Println("hi")
            checkConnect(target, name, memMap)  
        }
}

// TODO add timestamps
func checkConnect(target string, name string, memMap map[string]bool) {
            // fmt.Println("hi")
            // results := make([]string, 0)
            // results[0] = name
            _, err := net.Dial("tcp", target)
            // no err
            // errHandler(err, "Can not connect to " + target, false)
            if err != nil {
                _, ok := memMap[target]
                if(ok == true) {
                    delete(memMap, target)
                    fmt.Println(target + " has left")
                }
                return
            }
            _, ok := memMap[target]
            if ok == false {
                memMap[target] = true;
            }

}


// Get DNS of the host
func getdns(vm string) string{
    cmd := exec.Command("/usr/bin/dig","+short", vm)
    resultsBytes, err := cmd.CombinedOutput()
    resultsBytes = bytes.Trim(resultsBytes, "\n")
    errHandler(err,"Couldn't lookup machine address:" + vm, true )
    return string(resultsBytes)
}


// TODO1: Check whether items in the addrList are in the list (or map), if not, add them into the map.
// TODO2: Check whether there are items not in the list but in the map, put them into an array and send the array into disconnects
func checkItem(addrList []string, memMap map[string]string, disconnects chan []string) {

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

// check the number of args
func checkArg() {
    if len(os.Args) != 4  {     
        fmt.Println("Usage: " + os.Args[0] + " <name> <port> <n>")  
        os.Exit(0)  
    }
}