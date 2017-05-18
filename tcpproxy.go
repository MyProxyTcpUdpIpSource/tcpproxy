/* =====================================================================
 * Copyright (C) 2017 mincore All Right Reserved.
 *      Author: mincore@163.com
 *    Filename: tcpproxy.go
 * Description:
 * =====================================================================
 */
package main

import (
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Println("usage: listenAddr remoteAddr")
		log.Println("usage: 127.0.0.1:80 192.168.1.1:80")
		os.Exit(-1)
	}

	localAddr := os.Args[1]
	remoteAddr := os.Args[2]

	ln, err := net.Listen("tcp", localAddr)
	if err != nil {
		log.Println("listen on ", localAddr, " failed")
		os.Exit(-1)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handleConn(conn, remoteAddr)
	}
}

func handleConn(localConn net.Conn, remoteAddr string) {
	defer localConn.Close()

	remoteConn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		log.Println("connect to", remoteAddr, "failed")
		return
	}
	defer remoteConn.Close()

	log.Println("connect from", localConn.RemoteAddr().String(), "to", remoteAddr)

	go copyConn(remoteConn, localConn)
	copyConn(localConn, remoteConn)

	log.Println("disconnect from", localConn.RemoteAddr().String(), "to", remoteAddr)
}

func copyConn(from, to net.Conn) {
	for {
		var buffer [4096]byte
		n, err := from.Read(buffer[0:])
		if err != nil {
			return
		}

		for i := 0; i < n; {
			j, err := to.Write(buffer[i:n])
			if err != nil {
				return
			}
			i += j
		}
	}
}
