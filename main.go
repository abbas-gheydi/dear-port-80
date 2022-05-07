package main

import (
	"flag"
	"io"
	"log"
	"net"
)

var (
	listenAddr *string = flag.String("listen", "0.0.0.0:8080", "listen Address")
	ssh        *string = flag.String("ssh", "127.0.0.1:22", "ssh address")
	enableSSH  *bool   = flag.Bool("enable_ssh", true, "enable ssh routing")
	http       *string = flag.String("http", "127.0.0.1:80", "http server address")
)

func init() {
	flag.Parse()
	log.Printf("enableSSH: %v\n", *enableSSH)
	log.Printf("SSH_Proxying: %v\n", *ssh)
	log.Printf("HTTP_Proxying: %v\n", *http)
	log.Println("proxy server Started and Listening on:", *listenAddr)

}

func main() {

	listener, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		inputConnection, err := listener.Accept()
		if err != nil {
			log.Println("error Accepting Connection", err)
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

			defer inputConnection.Close()

			isSSH, data := sshDetector()
			if isSSH && *enableSSH {
				remoteService = *ssh
				log.Println(inputConnection.RemoteAddr(), "Connected as SSH Connection")
			} else {
				remoteService = *http
				log.Println(inputConnection.RemoteAddr(), "Connected as HTTP Connection")
			}

			outputConnection, err := net.Dial("tcp", remoteService)

			if err != nil {
				log.Println("error dialing remote addr", err)
				return

			}
			defer outputConnection.Close()
			outputConnection.Write(data)
			go io.Copy(outputConnection, inputConnection)
			io.Copy(inputConnection, outputConnection)

		}()
	}
}
