package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		log.Fatalf("Error connecting to server: %v", err)
	}
	defer conn.Close()

	// データ送信
	sendMyProtocolPacket(conn, 1, 1, []byte("Hello, Server!"))
	receiveResponse(conn)

	// 終了リクエスト送信
	sendMyProtocolPacket(conn, 1, 2, []byte{})
}

func sendMyProtocolPacket(conn net.Conn, version byte, opCode byte, data []byte) {
	packet := make([]byte, 4+len(data))
	packet[0] = version
	packet[1] = opCode
	binary.BigEndian.PutUint16(packet[2:4], uint16(len(data)))
	copy(packet[4:], data)

	_, err := conn.Write(packet)
	if err != nil {
		log.Fatalf("Error sending packet: %v", err)
	}
	fmt.Printf("Sent packet: Version=%d, OpCode=%d, Data=%s\n", version, opCode, string(data))
}

func receiveResponse(conn net.Conn) {
	header := make([]byte, 4)
	_, err := io.ReadFull(conn, header)
	if err != nil {
		log.Fatalf("Error reading response header: %v", err)
	}

	version := header[0]
	opCode := header[1]
	dataLength := binary.BigEndian.Uint16(header[2:4])

	data := make([]byte, dataLength)
	_, err = io.ReadFull(conn, data)
	if err != nil {
		log.Fatalf("Error reading response data: %v", err)
	}

	fmt.Printf("Received response: Version=%d, OpCode=%d, Data=%s\n", version, opCode, string(data))
}
