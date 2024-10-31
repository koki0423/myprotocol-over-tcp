package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
)

type EthernetFrame struct {
	Preamble   [7]byte // プリアンブル
	DstMacAddr [6]byte // 宛先MACアドレス
	SrcMacAddr [6]byte // 送信元MACアドレス
	Type       [2]byte // イーサタイプ
	Payload    []byte  // ペイロード
}

type Arp struct {
	HardwareType     [2]byte // ハードウェアタイプ（イーサネット: 0x0001）
	ProtocolType     [2]byte // プロトコルタイプ（IPv4: 0x0800）
	HardwareSize     byte    // ハードウェアサイズ（通常6）
	ProtocolSize     byte    // プロトコルサイズ（通常4）
	OperationCode    [2]byte // オペレーションコード（リクエスト: 0x0001、リプライ: 0x0002）
	SenderMacAddress [6]byte // 送信元MACアドレス
	SenderIPAddress  [4]byte // 送信元IPアドレス
	TargetMacAddress [6]byte // 宛先MACアドレス
	TargetIPAddress  [4]byte // 宛先IPアドレス
}

type ICMP struct {
	Type     byte    // タイプ
	Code     byte    // コード
	Checksum [2]byte // チェックサム
	ID       [2]byte // 識別子
	Seq      [2]byte // シーケンス番号
	Data     []byte  // データ
}

type LocalIpMacAddr struct {
	LocalMacAddr [6]byte
	LocalIpAddr  [4]byte
}

var (
	EtherTypeIPv4       = [2]byte{0x08, 0x00}
	EtherTypeARP        = [2]byte{0x08, 0x06}
	EtherTypeAppleTalk  = [2]byte{0x80, 0x9b}
	EtherTypeIEEE802_1q = [2]byte{0x81, 0x00}
	EtherTypeIPv6       = [2]byte{0x86, 0xdd}
)

func NewEthernetFrame(dstMacAddr, srcMacAddr, ethType string) EthernetFrame {
	var ethernet EthernetFrame

	// プリアンブルの設定
	ethernet.Preamble = [7]byte{0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xab}

	// MACアドレスをバイト配列に変換して設定
	dstMac, err := hex.DecodeString(dstMacAddr)
	if err != nil || len(dstMac) != 6 {
		log.Fatalf("Invalid destination MAC address: %v", err)
	}
	srcMac, err := hex.DecodeString(srcMacAddr)
	if err != nil || len(srcMac) != 6 {
		log.Fatalf("Invalid source MAC address: %v", err)
	}
	copy(ethernet.DstMacAddr[:], dstMac)
	copy(ethernet.SrcMacAddr[:], srcMac)

	// EtherTypeの設定
	switch ethType {
	case "IPv4":
		ethernet.Type = [2]byte{0x08, 0x00}
	case "IPv6":
		ethernet.Type = [2]byte{0x86, 0xDD}
	case "ARP":
		ethernet.Type = [2]byte{0x08, 0x06}
	default:
		log.Fatalf("Unknown EtherType: %s", ethType)
	}

	return ethernet
}

func NewArpRequest(localif LocalIpMacAddr, targetip string) Arp {
	return Arp{
		HardwareType:     [2]byte{0x00, 0x01},
		ProtocolType:     [2]byte{0x08, 0x00},
		HardwareSize:     6,
		ProtocolSize:     4,
		OperationCode:    [2]byte{0x00, 0x01},
		SenderMacAddress: localif.LocalMacAddr,
		SenderIPAddress:  localif.LocalIpAddr,
		TargetMacAddress: [6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		TargetIPAddress:  iptobyte(targetip),
	}
}

func NewICMPRequest(id, seq uint16, data []byte) ICMP {
	icmp := ICMP{
		Type: 8,
		Code: 0,
		ID:   [2]byte{byte(id >> 8), byte(id & 0xff)},
		Seq:  [2]byte{byte(seq >> 8), byte(seq & 0xff)},
		Data: data,
	}
	icmp.Checksum = calculateChecksum(icmp)
	return icmp
}
func calculateChecksum(icmp ICMP) [2]byte {
	// ダミーのチェックサムを0に設定しデータをシリアライズ
	icmp.Checksum = [2]byte{0x00, 0x00}
	var buf []byte
	buf = append(buf, icmp.Type, icmp.Code)
	buf = append(buf, icmp.Checksum[:]...)
	buf = append(buf, icmp.ID[:]...)
	buf = append(buf, icmp.Seq[:]...)
	buf = append(buf, icmp.Data...)

	// チェックサム計算
	sum := 0
	for i := 0; i < len(buf)-1; i += 2 {
		sum += int(binary.BigEndian.Uint16(buf[i : i+2]))
	}
	if len(buf)%2 == 1 {
		sum += int(buf[len(buf)-1]) << 8
	}
	for (sum >> 16) > 0 {
		sum = (sum & 0xFFFF) + (sum >> 16)
	}
	checksum := ^uint16(sum)

	return [2]byte{byte(checksum >> 8), byte(checksum & 0xff)}
}

func iptobyte(ipStr string) [4]byte {
	ip := net.ParseIP(ipStr).To4()
	if ip == nil {
		log.Fatalf("Invalid IP address: %s", ipStr)
	}
	var ipBytes [4]byte
	copy(ipBytes[:], ip)
	return ipBytes
}

func main() {
	//RAWソケット作成
	conn, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		log.Fatalf("Error creating RAW Socket: %v", err)
	}

	defer conn.Close()

	// ARPリクエストを送信

	//ICMPパケット作成
	//icmp := NewICMPRequest(1, 1, []byte("ping"))
}

func ArpSend
