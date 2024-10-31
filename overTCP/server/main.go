package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}
	defer listener.Close()
	fmt.Println("Server is listening on port 9000...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		// ヘッダーを受信
		header := make([]byte, 4)
		_, err := io.ReadFull(conn, header)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading header: %v", err)
			}
			return
		}

		// プロトコルバージョンとオペコードを解析
		version := header[0]
		opCode := header[1]
		dataLength := binary.BigEndian.Uint16(header[2:4])

		fmt.Printf("Received packet: Version=%d, OpCode=%d, DataLength=%d\n", version, opCode, dataLength)

		// データ長に基づいてデータを受信
		data := make([]byte, dataLength)
		_, err = io.ReadFull(conn, data)
		if err != nil {
			log.Printf("Error reading data: %v", err)
			return
		}
		fmt.Printf("Data: %s\n", string(data))

		// オペコードに基づいて応答を送信
		if opCode == 1 { // データ受信応答
			response := constructMyProtocolPacket(1, 1, []byte("Data received"))
			conn.Write(response)
		} else if opCode == 2 { // 終了
			fmt.Println("Received termination request. Closing connection.")
			return
		}
	}
}

// MyProtocolパケットの構築
func constructMyProtocolPacket(version byte, opCode byte, data []byte) []byte {
	packet := make([]byte, 4+len(data))
	packet[0] = version
	packet[1] = opCode
	binary.BigEndian.PutUint16(packet[2:4], uint16(len(data)))
	copy(packet[4:], data)
	return packet
}
