package Tcp

import (
	"errors"
	"fmt"
	"github.com/Blizzardx/GameServer/Core/Network/MessageRegister"
)

type pingHelper struct {
	PingMsgId   int32
	PongMsgId   int32
	PingMsgBody interface{}
	PongMsgBody interface{}
}

func createPingHelper(pingMsgId int32, pongMsgId int32) (*pingHelper, error) {
	pingMsgBody := MessageRegister.GetMessageInstanceByMsgId(pingMsgId)
	if nil == pingMsgBody {
		return nil, errors.New(fmt.Sprint("can't create ping msg body by msg id", pingMsgId))
	}
	pongMsgBody := MessageRegister.GetMessageInstanceByMsgId(pongMsgId)
	if nil == pongMsgBody {
		return nil, errors.New(fmt.Sprint("can't create pong msg body by msg id", pongMsgId))
	}
	return &pingHelper{PingMsgId: pingMsgId, PongMsgId: pongMsgId, PingMsgBody: pingMsgBody, PongMsgBody: pongMsgBody}, nil
}
func (self *pingHelper) CheckSendPong(msgId int32) (bool, interface{}) {
	if msgId == self.PingMsgId {
		return true, self.PingMsgBody
	}
	return false, nil
}
