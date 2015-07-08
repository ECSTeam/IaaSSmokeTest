package main

import (
	"fmt"
	"encoding/json"
	"os"
	"net"
	"strconv"
	"bufio"
)

var configuration struct {
	RemoteHost string `json:"remoteHost"`
	RemoteConnectionDetail[] connectionDetail `json:"remoteConnectionDetails"`
	LocalConnectionDetail[] connectionDetail `json:"localConnectionDetails"`
}


type connectionDetail struct {
	Port string `json:"port"`
	Protocol string `json:"protocol"`
	Description string `json:"description"`
	
}

var resultSet struct {
	TCPListeners int
	UDPListeners int
	TCPRequestsReceived int
	UDPRequestsReceived int
	TCPResponsesReceived int
}


func main () {
	
	
    initialize()
    openLocalPorts()
    waitForInput("Hit Enter when all ports are listening: \n\n")
    connectToRemotePorts()
    waitForInput("\nHit Enter to shutdown listeners: \n\n")
    writeOutput()
}

func initialize () {
	
	configDir := os.Getenv("IAASTESTCONFIGDIR")
	configFileName := configDir + "config.json"
	
	configFile, err := os.Open(configFileName)
    if err != nil {
        fmt.Println("error in opening config file", err.Error())
    	os.Exit(1)
    }

    jsonParser := json.NewDecoder(configFile)
    if err = jsonParser.Decode(&configuration); err != nil {
        fmt.Println("error in parsing config file", err.Error())
    	os.Exit(1)
    }
    
    resultSet.TCPListeners = 0
    resultSet.UDPListeners = 0
    resultSet.TCPRequestsReceived = 0
    resultSet.UDPRequestsReceived = 0
    resultSet.TCPResponsesReceived = 0
}

func waitForInput(inStr string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(inStr)
	reader.ReadString('\n')
}

func openLocalPorts () {
	for _, detail := range configuration.LocalConnectionDetail {
    	go openLocalPort (detail.Port, detail.Protocol, detail.Description)
    }
}

func openLocalPort(port string, protocol string, description string) {
	fmt.Printf("Setting up a %s listener on port %s used for %s\n", protocol, port, description)

	if protocol=="tcp" {
		openLocalTCP(port, description)
	} else if protocol=="udp" {
		openLocalUDP(port, description)
	}
}

func openLocalTCP(port string, description string) {
	
	ln, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		fmt.Printf("error in setting up TCP listener on port %s used for %s.\n", port, description)
		fmt.Println("error is ", err.Error())
	} else {
	
		resultSet.TCPListeners = resultSet.TCPListeners + 1
		
		defer ln.Close()
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Printf("error in accepting TCP connections on port %s used for %s.\n", port, description)
				fmt.Println("error is ", err.Error())
			} else {
				fmt.Printf("Received message %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())
				go handleConnection(conn)
				resultSet.TCPRequestsReceived = resultSet.TCPRequestsReceived + 1
			}
		}
	}
	
}

func openLocalUDP(port string, description string) {
	
	udpAddr := "0.0.0.0:" + port
	serverAddr,_ := net.ResolveUDPAddr("udp", udpAddr)
	ln, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		fmt.Printf("error in setting up UDP listener on port %s used for %s.\n", port, description)
		fmt.Println("error is ", err.Error())
	} else {
	
		resultSet.UDPListeners = resultSet.UDPListeners + 1
		
		defer ln.Close()
	
		buf := make([]byte, 1024)
		for {
			reqLen, addr, err := ln.ReadFromUDP(buf)
			if err != nil {
				fmt.Printf("error in accepting UDP connections on port %s used for %s.\n", port, description)
				fmt.Println("error is ", err.Error())
			} else {
				fmt.Printf("Received %d byte UDP message from %s \n", reqLen, addr )
				resultSet.UDPRequestsReceived = resultSet.UDPRequestsReceived + 1
			}
		}
	}
	
}

func handleConnection(conn net.Conn)  {
	buf := make([]byte, 1024)
	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading", err.Error())
	} else {
		message := "Recived message "
	    message += strconv.Itoa(reqLen)
	    message += " bytes long\n"
	    conn.Write([]byte(message))
    }
	conn.Close()
}

func connectToRemotePorts() {
	for _, detail := range configuration.RemoteConnectionDetail {
    	connectToRemotePort (detail.Port, detail.Protocol, detail.Description)
    }
}

func connectToRemotePort (port string, protocol string, description string) {
	host := configuration.RemoteHost
	fmt.Printf("connecting via %s to host:port %s:%s used for %s\n", protocol, host, port, description)

	if protocol=="tcp" {
		connectViaTCP(port, host, description)
	} else if protocol=="udp" {
		connectViaUDP(port, host, description)
	}
}

func connectViaTCP(port string, host string, description string) {
	url := host + ":" + port
	conn, err := net.Dial("tcp", url)
	if err != nil {
		fmt.Printf("error in connecting to TCP listener at %s:%s\n", host, port)
		fmt.Printf("error is : %s", err.Error())
	} else {
		fmt.Fprintf(conn, "HEAD\n")
		status, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("error in reading return stream from %s:%s\n", host, port)
			fmt.Printf("error is : %s", err.Error())
		} else {
			fmt.Println("response = ", status)
			resultSet.TCPResponsesReceived = resultSet.TCPResponsesReceived + 1
		}
	}
}

func connectViaUDP (port string, host string, description string) {
	url := host + ":" + port
	serverAddr,err := net.ResolveUDPAddr("udp",url)
    if err != nil {
		fmt.Printf("error in resolving remote UDP address %s:%s\n", host, port)
		fmt.Printf("error is : %s", err.Error())
	} else {
 	 
	    conn, err := net.DialUDP("udp", nil, serverAddr)
	    if err != nil {
			fmt.Printf("error in connecting to UDP listener at %s:%s\n", host, port)
			fmt.Printf("error is : %s", err.Error())
		} else {
	 
		    defer conn.Close()
		    
		    msg := "HEAD\n"
		    buf := []byte(msg)
		    _,err = conn.Write(buf)
		    if err != nil {
		        fmt.Printf("error in writing to UDP listener at %s:%s\n", host, port)
		        fmt.Println("error is :", err.Error())
		    } else {
		    	fmt.Println("wrote the buffer in UDP land")
		    }
	    }
	    
	}
	
}

func writeOutput() {
	fmt.Printf("Total TCP Listeners Set up:  %d\n", resultSet.TCPListeners)
	fmt.Printf("Total UDP Listeners Set up:  %d\n", resultSet.UDPListeners)
	fmt.Println("")
	fmt.Printf("Total TCP Requests Received:  %d\n", resultSet.TCPRequestsReceived)
	fmt.Printf("Total UDP Requests Received:  %d\n", resultSet.UDPRequestsReceived)
	fmt.Println("")
	fmt.Printf("Total TCP Responses Recieved:  %d\n", resultSet.TCPResponsesReceived)
}