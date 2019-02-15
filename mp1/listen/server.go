package main

import(  
    "fmt"  
    "os"  
    "net"
    "bufio"
)

func main(){
    checkArg()
    name := os.Args[1]
    port := os.Args[2]
    numMem := os.Args[3]
    server, err := net.Listen("tcp", ":" + port)
    errHandler(err, "Can not start server!", true)
    for {
        conn, err := server.Accept()
        errHandler(err, "Can not open connection!")
        go func(conn net.Conn) {
            defer conn.Close()
            input := bufio.NewReader(conn)
            output := bufio.NewWriter(conn)
            pattern, err := input.ReadString('\n')
            pattern = pattern[:len(pattern) - 1]
            fmt.Println("Pattern:" + pattern)
            s := strings.Split(pattern, ":")
            message := s[1]
        } (conn)
    }
}

