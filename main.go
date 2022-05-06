package main

import (
	"flag"
	"io"
	"log"
	"net"
)

var (
	listenAddr *string = flag.String("listen", "0.0.0.0:8080", "listen Address")

	ssh *string = flag.String("ssh", "localhost:22", "ssh address")

	http *string = flag.String("http", "localhost:80", "http server address")
)

func init() {
	flag.Parse()

	log.Printf("Listening: %v ,SSH_Proxying:%v ,HTTP_Proxying:%v\n", *listenAddr, *ssh, *http)
}

func main() {

	listener, err := net.Listen("tcp", *listenAddr)

	if err != nil {
		log.Fatal(err)
	}

	for {

		inputConnection, err := listener.Accept()
		if err != nil {
			log.Println("error accepting connection", err)
			continue
		}

		go func() {
			var remoteService string

			sshDetector := func() (bool, []byte) {
				data := make([]byte, 256)
				isSSH := false
				dataSize, err := inputConnection.Read(data)
				if err != nil {
					log.Println("failed to read input Connection", err)
				}

				if data[0] == 'S' && data[1] == 'S' && data[2] == 'H' {

					isSSH = true
					data = data[:dataSize]

				}
				return isSSH, data

			}

			isSSH, data := sshDetector()

			if isSSH {
				remoteService = *ssh
				log.Println(inputConnection.RemoteAddr(), "Connected as SSH conenction")
			} else {
				remoteService = *http
				log.Println(inputConnection.RemoteAddr(), "Connected as HTTP conenction")
			}

			outputConnection, err := net.Dial("tcp", remoteService)
			if err != nil {
				log.Println("error dialing remote addr", err)
				return
			}

			outputConnection.Write(data)
			go io.Copy(outputConnection, inputConnection)
			io.Copy(inputConnection, outputConnection)
			outputConnection.Close()
			inputConnection.Close()
		}()
	}
}
