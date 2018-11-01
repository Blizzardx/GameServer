package Tcp

import (
	"fmt"
	"github.com/Blizzardx/GameServer/Core/Common"
	"github.com/Blizzardx/GameServer/Core/Common/Queue"
	"github.com/Blizzardx/GameServer/Core/Network/Codec/Core"
	"net"
)

func StartListen(port string, codeC *Core.NetworkCodeC, receiveQueue *Queue.NoneBlockingQueue, pingMsgId int32, pongMsgId int32) {
	l, err := net.Listen("tcp", port)
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
			NewPipeline(c, codeC, receiveQueue)
		})
	}
}
func StartDail(add string, port string, codeC *Core.NetworkCodeC, receiveQueue *Queue.NoneBlockingQueue, pingMsgId int32, pongMsgId int32) {
	conn, err := net.Dial("tcp", add+":"+port)
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}
	Common.SafeCall(func() {
		NewPipeline(conn, codeC, receiveQueue)
	})
}
