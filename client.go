package main

import (
  "fmt"
  "net"
  "bufio"
  "os"
  "strings"
  "bytes"
)

const serverIP = "192.168.56.1"

func Log(v ...interface{}) {
  fmt.Println(v...)
}

func main() {
  conn, err := net.Dial("tcp", serverIP + ":6666")
  if err != nil {
    Log("Error dialing")
    return
  }
  
  quit := make(chan bool)
  // Send messages
  go func() {
    fmt.Printf("Enter username: ")
    reader := bufio.NewReader(os.Stdin)
    username, err := reader.ReadString('\n')
    if err != nil {
      Log("Invalid username")
    }
    conn.Write([]byte(username))
    for {
      //fmt.Printf("You: ")
      reader := bufio.NewReader(os.Stdin)
      message, err := reader.ReadString('\n')
      if err != nil {
        Log("\nExiting")
        quit <- true
      }
      if !bytes.Equal([]byte(message), []byte("\n")) {
        conn.Write([]byte(message))
      }
    }
  }()
  // Receive messages
  go func() {
    receivedMessageBuffer := make([]byte, 1024)
    for i := 0; i < 1024; i++ {
      receivedMessageBuffer[i] = ' '
    }
    _, err = conn.Read([]byte(receivedMessageBuffer))
    for err == nil {
      receivedMessage := strings.TrimSpace(string(receivedMessageBuffer))
      Log(receivedMessage)
      for i := 0; i < 1024; i++ {
        receivedMessageBuffer[i] = ' '
      }
      _, err = conn.Read([]byte(receivedMessageBuffer))
    }
    if err != nil {
      Log("\nExiting")
      quit <- true
    }
  }()
  if ( <- quit) {
    return
  }
}