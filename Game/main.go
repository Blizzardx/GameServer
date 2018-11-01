package main

import (
	"fmt"
	"github.com/Blizzardx/GameServer/Core/Common"
	"github.com/Blizzardx/GameServer/Core/Network/Tcp"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			break
		}
		fmt.Println("xx")

		Common.SafeCall(func() {
			Tcp.NewPipeline(c, nil, nil)
		})
	}
}
