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
    "strconv"
    "time"
)

type Message struct{ 
    UserName string
    // Host string
    Address string
    Text string
    TimeStamp map[string]int
}

// a map to store the memberships in the chatroom
func main() {
    // To check the input
    checkArg()
    // To get the name, port and num of people from input
    name := os.Args[1]
    port := os.Args[2]
    numMem, _ := strconv.Atoi(os.Args[3])
    localHost, err := os.Hostname()
    receivedMsg := make(map[string][]int)
    if err != nil {
        panic(err)
    }
    address := getdns(localHost)+ ":" + port
    // localId := getId(localHost)
    // fmt.Print
    memMap := make(map[string]string)
    // addressHostMap := make(map[string]string)
    localTimeStamp := make(map[string]int)
    localTimeStamp[address] = 0
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
        errHandler(err, "#Can not open connection!", true)
        if err != nil {
            return
        }

        go func(conn net.Conn) {
            defer conn.Close()
            for {
                inBuf := make([]byte, 512)
                size, err:=conn.Read(inBuf)
                if err != nil{
                    // fmt.Println("Read Error:",err.Error());
                    return
                }
                //fmt.Println("data from client:",string(buf),"size:",size)                                                                                                      
                var chatMsg Message
                err = json.Unmarshal(inBuf[:size], &chatMsg)
                if err!=nil{
                    fmt.Println("#Unmarshal Error:", err.Error());
                    return
                }
                _, ok := memMap[chatMsg.Address]
                if(ok == false || len(memMap[chatMsg.Address]) == 0) {
                    memMap[chatMsg.Address] = chatMsg.UserName
                    _, timeOk := localTimeStamp[chatMsg.Address]
                    if(timeOk == false){
                        localTimeStamp[chatMsg.Address] = 0
                    }
                }
                if(len(chatMsg.Text) != 0){
                        received := contains(receivedMsg[chatMsg.Address], chatMsg.TimeStamp[chatMsg.Address])
                        if(received == true){
                            continue
                        }
                        receivedMsg[chatMsg.Address] = append(receivedMsg[chatMsg.Address], chatMsg.TimeStamp[chatMsg.Address])
                        go func(){// keep judging whether a message can be delivered or not 
                            for {
                                var toDeliver = checkDeliver(chatMsg, localTimeStamp, memMap)
                                if toDeliver == true{
                                    fmt.Println(chatMsg.UserName + ": " + chatMsg.Text)
                                    localTimeStamp[chatMsg.Address] = localTimeStamp[chatMsg.Address] + 1
                                    // all the tests
                                    // fmt.Println("localTimeStamp after deliver:",localTimeStamp)
                                    // fmt.Println("toDeliver",toDeliver)
                                    // fmt.Println("chatMsg.Timestamp:",chatMsg.TimeStamp)
                                    // fmt.Println("localTimeStamp before deliver:",localTimeStamp)
                                    // fmt.Println("receivedMsg after deliver:",receivedMsg)
                                    return
                                }
                            }
                            
                        }()
                        // var toDeliver = checkDeliver(chatMsg, localTimeStamp, memMap)
                        
                        
                    
                }
            }
        } (conn)
    }
    }()

    var hmsg string // the message to print into the console (the receiver)
    go func(){
        for {
            select {
            case hmsg = <- inMsg :
                fmt.Println(hmsg)

            default:
            }
        }
        }()
    
    
     //key: host, value:name
    go func(){
        for{
            checkConnectAll(targets, name, memMap)
        }
    }()

    // outMsg := make(chan []byte, 3)
    var msg []byte
    outInit := make(chan []byte, 3)

    // multicast the messages
    go func() {
        for{ 
            select {
           
            case msg = <- outInit:
                for item, _ := range memMap {
                    if(len(item) > 0 && item != address){
                        conn, err := net.DialTimeout("tcp", item, 100 * time.Millisecond)
                        // handle the error 
                        if err != nil {
                            // fmt.Println("checkConnect error")
                            // fmt.Println("len of memmap",len(memMap))
                            continue
                        }
                        conn.Write(msg)
                    }
                }
            default:
            }   
        }
    }()


    ready := 0
    current_mem := 0
    for {
        if(ready == 0) {// not ready
            for {
                var chatText string
                // var TimeStamp map[string]int
                outMsg := &Message{UserName:name, Address:address, Text:chatText, TimeStamp:localTimeStamp}
                b, err := json.Marshal(outMsg)
                if err != nil {
                    fmt.Println("#encoding faild")
                }
                if(current_mem < len(memMap)) {
                    outInit <- b
                    current_mem += 1
                }
                
                if(len(memMap) ==  numMem){
                    ready = 1
                    fmt.Println("READY")
                    break
                }
            }
        } else {// ready
            for {
                var chatText string
            
                scanner := bufio.NewScanner(os.Stdin)
                if scanner.Scan() {
                     chatText = scanner.Text()
                }
                // msgTimeStamp:= make([]int,len(localTimeStamp))
                // copy(msgTimeStamp,localTimeStamp)
                localTimeStamp[address] = localTimeStamp[address] + 1
                outMsg := &Message{UserName:name, Address:address, Text:chatText, TimeStamp:localTimeStamp}
                b, err := json.Marshal(outMsg)
                if err != nil {
                    fmt.Println("#encoding faild")
                }
                outInit <- b
            }
        }
    }

    
}

// implement the algorithm of vector timestamps
func checkDeliver(chatMsg Message,localTimeStamp map[string]int,memMap map[string]string)bool{
    var toDeliver bool = true
    fromAddress := chatMsg.Address
    if localTimeStamp[fromAddress]+1<chatMsg.TimeStamp[fromAddress]{
        toDeliver = false
    }
    for k:=range localTimeStamp{
        if(k == fromAddress){
            continue
        }
        if chatMsg.TimeStamp[k]>localTimeStamp[k]{
            toDeliver = false
        }     
    }
    return toDeliver
}

// set up the listener
func setupServer(port string) net.Listener {
    server, err := net.Listen("tcp", ":" + port)
    errHandler(err, "#Can not start server!", true)
    return server
}

// dial all the 10 VMs to check whether they can connect successfully (run in the background) 
func checkConnectAll(targets []string, name string, memMap map[string]string) {
        // var results []string
        for _, target := range targets {
            // fmt.Println("hi")
            checkConnect(target, name, memMap)  
        }
}

// TODO add timestamps
func checkConnect(target string, name string, memMap map[string]string) {
            // fmt.Println("hi")
            // results := make([]string, 0)
            // results[0] = name
            _, err := net.DialTimeout("tcp", target, 100 * time.Millisecond)
            // no err
            // errHandler(err, "Can not connect to " + target, false)
            if err != nil {
                _, ok := memMap[target]
                if(ok == true && len(memMap[target]) > 0) {// if cannot dial, then print left information
                    fmt.Println(memMap[target] + " has left")
                    delete(memMap, target)
                }
                return
            }
            _, ok := memMap[target]
            if ok == false {// add the node into the memMap
                memMap[target] = "";
            }

}

// Get DNS of the host
func getdns(vm string) string{
    cmd := exec.Command("/usr/bin/dig","+short", vm)
    resultsBytes, err := cmd.CombinedOutput()
    resultsBytes = bytes.Trim(resultsBytes, "\n")
    errHandler(err,"#Couldn't lookup machine address:" + vm, true )
    return string(resultsBytes)
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
        fmt.Println("#Usage: " + os.Args[0] + " <name> <port> <n>")  
        os.Exit(0)  
    }
}

// check whether a slice of int contains an element
func contains(s []int, e int) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}
