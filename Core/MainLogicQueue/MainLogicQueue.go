package MainLogicQueue

import (
	"github.com/Blizzardx/GameServer/Core/Common"
	"github.com/Blizzardx/GameServer/Core/Common/Queue"
	"github.com/Blizzardx/GameServer/Core/Network/Pipeline"
)

type mainLogicQueue struct {
	Queue.NoneBlockingQueue
}

var msgHandler = map[int32]func(msgId int32, msgBody interface{}, pipeline Pipeline.Pipeline){}
var logicQueue = &mainLogicQueue{}

func RegisterMessageHandler(msgId int32, handler func(msgId int32, msgBody interface{}, pipeline Pipeline.Pipeline)) {
	if _, ok := msgHandler[msgId]; ok {
		return
	}
	msgHandler[msgId] = handler
}
func StartLogicQueue() {
	go Common.SafeCall(func() {
		for {
			var msgBuffer []interface{}
			for {
				logicQueue.Pick(&msgBuffer)
				// do send
				for _, msgElem := range msgBuffer {
					messageElement := msgElem.(*Pipeline.PipelineMessageElement)
					if handler, ok := msgHandler[messageElement.MsgId]; ok {
						handler(messageElement.MsgId, messageElement.MsgBody, messageElement.Session)
					} else {
						//todo log error
					}
				}
				// clear send buffer
				msgBuffer = msgBuffer[0:0]
			}
		}
	})
}
