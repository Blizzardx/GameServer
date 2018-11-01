package Impl

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Blizzardx/GameServer/Core/Network/MessageRegister"
)

type Codec_Json struct {
}

func (self *Codec_Json) Decode(msgId int32, buffer []byte) (interface{}, error) {
	// create msg body instance
	msgBody := MessageRegister.GetMessageInstanceByMsgId(msgId)
	if nil == msgBody {
		return nil, errors.New(fmt.Sprint("msg not found", msgId))
	}
	err := json.Unmarshal(buffer, msgBody)

	return msgBody, err
}
func (self *Codec_Json) Encode(msgId int32, msgBody interface{}) ([]byte, error) {
	return json.Marshal(msgBody)
}
