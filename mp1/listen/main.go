package main

import(  
    "fmt"  
    "os"  
    "net"
    "bufio"
   // "strings"
    "os/exec" 
    "bytes"
    "encoding/json"
)

type Message struct{ 
    UserName string
    Address string
    Text string
    TimeStamp string
}


// a map to store the memberships in the chatroom
func main() {
    // To check the input
    checkArg()
    // To get the name, port and num of people from input
    name := os.Args[1]
    port := os.Args[2]
    localHost,err := os.Hostname()
    if err != nil {
        panic(err)
    }
    address := getdns(localHost)+ ":" + port
    // numMem := os.Args[3]w
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
                inBuf:=make([]byte,10000)
                size,err:=conn.Read(inBuf)
                if err!=nil{
                    fmt.Println("Read Error:",err.Error());
                    return
                }
                //fmt.Println("data from client:",string(buf),"size:",size)                                                                                                      
                var chatMsg Message
                err=json.Unmarshal(inBuf[:size],&chatMsg)
                if err!=nil{
                    fmt.Println("Unmarshal Error:",err.Error());
                    return
                }
                fmt.Println(chatMsg.UserName+":"+chatMsg.Text)
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
    
    memMap := make(map[string]bool)
    //memInfo := make(map[string]string) //key: host, value:name
    go func(){
        for{
            checkConnectAll(targets, name, memMap)
        }
    }()

    outMsg := make(chan string, 3)
    var msg string
    outInit := make(chan string,1)

    // multicast the messages
    go func() {
        for{ 
            select {
            case msg = <- outMsg:
                for item, _ := range memMap {
                    if(len(item) > 0 && item != address){
                        conn, err := net.Dial("tcp", item)
                        // handle the error 
                        if err != nil {
                            fmt.Println("checkConnect error")
                            os.Exit(1)
                        }
                        conn.Write([]byte(msg))
                    }
                }
            case msg = <- outInit:
                for item, _ := range memMap {
                    if(len(item) > 0 && item != address){
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
        var chatText string
        scanner := bufio.NewScanner(os.Stdin)
        if scanner.Scan() {
             chatText = scanner.Text()
        }
        outMsg := &Message{UserName:os.Args[1],Address:address,Text:chatText,TimeStamp:""}

        b,err := json.Marshal(outMsg)
        if err != nil {
            fmt.Println("encoding faild")
        }
        outInit <- string(b)
    }
}

func setupServer(port string) net.Listener {
    server, err := net.Listen("tcp", ":" + port)
    errHandler(err, "Can not start server!", true)
    return server
}


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
