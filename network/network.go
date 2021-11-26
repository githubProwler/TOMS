package network

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

type Server struct {
	listener  net.Listener
	handlerFn func(string)
}

func (s *Server) Init(handlerFn func(string)) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("There was a problem listening to tcp socket ", err)
	}

	log.Printf("Initiallized on port %d", ln.Addr().(*net.TCPAddr).Port)

	s.listener = ln
	s.handlerFn = handlerFn
}

func (s *Server) Serve() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatal("There was a problem accepting connection ", err)
		}

		go func(conn net.Conn) {
			str, err := readResponse(conn)
			if err != nil {
				log.Fatal("There was a problem reading from client ", err)
			}

			length := 0
			if len(str) > 0 {
				length = len(str) - 1
			}
			s.handlerFn(str[:length])
		}(conn)
	}
}

func readResponse(conn net.Conn) (response string, err error) {
	reader := bufio.NewReader(conn)
	// _, err = reader.Discard(8)
	if err != nil {
		log.Fatal("There was a problem reading the request ", err)
		return
	}
	response, err = reader.ReadString('\n')
	// fmt.Print(response)
	return
}

// func handleConnection(conn net.Conn) {
// 	slice := make([]byte, 20)
// 	_, err := conn.Read(slice)
// 	if err != nil {
// 		log.Fatal("There was a problem reading from client ", err)
// 	}
// }

func (s *Server) GetAddress() string {
	conn, err := net.Dial("udp", "2.2.2.2:2220")
	if err != nil {
		log.Fatal("Failed to get my IP address ", err)
	}

	var addr, port, result string
	addr = strings.Split(conn.LocalAddr().String(), ":")[0]
	port = strconv.Itoa(s.listener.Addr().(*net.TCPAddr).Port)
	result = addr + ":" + port
	log.Println("[Network][GetAddress] Address: ", addr, " Port: ", port)
	return result
}

func SendMessage(message string, addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal("There was a problem connecting to server ", err)
	}

	data := []byte(message)
	_, err = conn.Write(data)
	if err != nil {
		log.Fatal("There was a problem writing to server ", err)
	}
	fmt.Println(conn.LocalAddr().String())
	conn.Close()
}
