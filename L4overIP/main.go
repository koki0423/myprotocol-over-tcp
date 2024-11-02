package main

import (
    "fmt"
    "golang.org/x/net/icmp"
    "golang.org/x/net/ipv4"
    "net"
    "os"
    "time"
)

func main() {
    // 対象のIPアドレスを指定
    target := "8.8.8.8" // Google DNSのIPアドレス

    // ICMP Echo Requestメッセージの作成
    icmpMessage := icmp.Message{
        Type: ipv4.ICMPTypeEcho, Code: 0,
        Body: &icmp.Echo{
            ID: os.Getpid() & 0xffff, Seq: 1,
            Data: []byte("HELLO-R-U-THERE"),
        },
    }
    messageData, err := icmpMessage.Marshal(nil)
    if err != nil {
        fmt.Printf("メッセージのマーシャルに失敗しました: %v\n", err)
        return
    }

    // ICMPプロトコルを使用したネットワーク接続の作成
    conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
    if err != nil {
        fmt.Printf("接続の作成に失敗しました: %v\n", err)
        return
    }
    defer conn.Close()

    // 送信先アドレスの設定
    dst, err := net.ResolveIPAddr("ip4", target)
    if err != nil {
        fmt.Printf("アドレスの解決に失敗しました: %v\n", err)
        return
    }

    // パケットの送信
    startTime := time.Now()
    n, err := conn.WriteTo(messageData, dst)
    if err != nil {
        fmt.Printf("パケットの送信に失敗しました: %v\n", err)
        return
    } else if n != len(messageData) {
        fmt.Printf("送信バイト数が一致しません: %v\n", err)
        return
    }
    fmt.Printf("Pingを送信しました。対象IP: %s\n", target)

    // パケットの受信
    reply := make([]byte, 1500)
    err = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
    if err != nil {
        fmt.Printf("読み取りタイムアウトの設定に失敗しました: %v\n", err)
        return
    }
    n, peer, err := conn.ReadFrom(reply)
    if err != nil {
        fmt.Printf("パケットの受信に失敗しました: %v\n", err)
        return
    }
    rtt := time.Since(startTime)

    // 受信したICMPメッセージの解析
    receivedMessage, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), reply[:n])
    if err != nil {
        fmt.Printf("メッセージの解析に失敗しました: %v\n", err)
        return
    }

    switch receivedMessage.Type {
    case ipv4.ICMPTypeEchoReply:
        fmt.Printf("Pingの応答あり。応答元IP: %s 時間: %v\n", peer.String(), rtt)
    default:
        fmt.Printf("予期しないICMPメッセージを受信しました: %+v\n", receivedMessage)
    }
}
