package main

import(  
    "fmt"  
    "os"  
    "net"
    // "bufio"
    // "strings"
)

func main(){
    checkArg()
    // name := os.Args[1]
    port := os.Args[2]
    // numMem := os.Args[3]
    server, err := net.Listen("tcp", ":" + port)
    errHandler(err, "Can not start server!", true)
    
    for {
        conn, err := server.Accept()
        fmt.Println("hi")
        errHandler(err, "Can not open connection!", true)
        go func(conn net.Conn) {
            defer conn.Close()
            for {
                buff := make([]byte, 256)
                length, err := conn.Read(buff)
                if err != nil {
                    break
                }
                fmt.Println(string(buff[:length]))
            }
        } (conn)
    }
}

