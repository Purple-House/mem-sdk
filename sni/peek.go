package sni

import (
	"fmt"
	"log"
	"net"
)

// func PeekSNI(conn net.Conn) (string, net.Conn, error) {
// 	buf := make([]byte, 512)
// 	n, err := conn.Read(buf)
// 	if err != nil {
// 		return "", nil, err
// 	}
// 	serverName, err := SniStream(buf[:n])
// 	if err != nil {
// 		log.Println("[SNI] error on found the with TLS servername", err)
// 		serverName, err = ExtractHostFromStream(buf[:n])
// 		if err != nil {
// 			log.Println("[SNI] error on found the servername", err)
// 		}
// 	}

// 	log.Printf("found the server name domain: %s", serverName)

// 	if err != nil {
// 		return "", nil, err
// 	}
// 	return serverName, &ConnBuffer{buf: buf[:n], conn: conn}, nil
// }

func PeekSNI(conn net.Conn) (string, net.Conn, error) {
	var buf []byte
	tmp := make([]byte, 4096)

	for {
		n, err := conn.Read(tmp)
		if err != nil {
			return "", nil, err
		}

		buf = append(buf, tmp[:n]...)

		sni, err := SniStream(buf)
		if err == nil && sni != "" {
			log.Printf("found server name: %s", sni)
			return sni, &ConnBuffer{buf: buf, conn: conn}, nil
		}

		// Prevent unbounded buffering
		if len(buf) > 64*1024 {
			return "", nil, fmt.Errorf("ClientHello too large or SNI not found")
		}
	}
}
