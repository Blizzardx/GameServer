package Tcp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/Blizzardx/GameServer/Core/Common"
	"github.com/Blizzardx/GameServer/Core/Common/Queue"
	"github.com/Blizzardx/GameServer/Core/Network/Codec/Core"
	"github.com/Blizzardx/GameServer/Core/Network/Pipeline"
	"net"
	"sync"
)

//解码状态
type DecodeStatus int32

const (
	decodeStatus_Id         DecodeStatus = iota + 1 // value --> 1 id 4 bit
	decodeStatus_BodyLength                         // value --> 2 body length 4 length
	decodeStatus_Body                               // value --> 3 body
)

type tcpPipeline struct {
	connection           net.Conn
	codeC                Core.NetworkCodeC
	sendQueue            *Queue.NoneBlockingQueue
	receivedQueue        *Queue.NoneBlockingQueue
	decodeStatus         DecodeStatus
	receivedBuffer       []byte
	currentMsgId         int32
	currentMsgBodyLength int32
	isClose              bool
	closeMutex           sync.Mutex
}

func NewPipeline(c net.Conn, codeC *Core.NetworkCodeC, receiveQueue *Queue.NoneBlockingQueue) *tcpPipeline {
	piepline := &tcpPipeline{}
	piepline.start(c, codeC, receiveQueue)
	return piepline
}
func (self *tcpPipeline) Send(msgId int32, msgBody interface{}) {
	if msgBody == nil {
		fmt.Println("can't send nil msg ")
		return
	}
	if self.isClose {
		fmt.Println("pipeline is closed")
		return
	}
	// add to send queue
	self.sendQueue.Add(&Pipeline.PipelineMessageElement{MsgId: msgId, MsgBody: msgBody})
}
func (self *tcpPipeline) Close() {
	self.doClose()
}

func (self *tcpPipeline) start(c net.Conn, codeC *Core.NetworkCodeC, mainLogicMsgQueue *Queue.NoneBlockingQueue) {
	self.connection = c
	self.codeC = codeC
	self.sendQueue = Queue.NewNoneBlockingQueue()
	self.receivedQueue = mainLogicMsgQueue
	self.isClose = false

	go Common.SafeCallWithCrashCallback(self.beginReceive, func() {
		Common.SafeCall(self.Close)
	})
	go Common.SafeCallWithCrashCallback(self.beginSend, func() {
		Common.SafeCall(self.Close)
	})
}
func (self *tcpPipeline) beginSend() {
	var sendBuffer []interface{}
	for {
		self.sendQueue.Pick(&sendBuffer)
		if self.isClose {
			return
		}
		// do send
		for _, msgElem := range sendBuffer {
			messageElement := msgElem.(*pieplineMessageElement)
			sendBuffer, err := self.getSendBuffer(messageElement)
			if err != nil {
				continue
			}
			if self.isClose {
				return
			}
			self.connection.Write(sendBuffer)
		}
		// clear send buffer
		sendBuffer = sendBuffer[0:0]
	}
	fmt.Sprintf("exit send ")
}
func (self *tcpPipeline) beginReceive() {
	var buffer []byte
	for {
		if self.isClose {
			return
		}
		len, err := self.connection.Read(buffer)
		if err != nil {
			// error
			self.Close()
			return
		}
		err = self.onReceived(buffer[0:len])
		if err != nil {
			// error
			self.Close()
			return
		}
	}
	fmt.Sprintf("exit receive ")
}
func (self *tcpPipeline) doClose() {
	if self.isClose {
		return
	}
	self.closeMutex.Lock()
	defer self.closeMutex.Unlock()

	if self.isClose {
		return
	}
	self.isClose = true
	self.connection.Close()
	self.sendQueue.Add(&pieplineMessageElement{msgId: -1, msgBody: nil})
}
func (self *tcpPipeline) getSendBuffer(messageElement *Pipeline.PipelineMessageElement) ([]byte, error) {
	sendBuffer, err := self.codeC.Encode(messageElement.MsgId, messageElement.MsgBody)
	if err != nil {
		return nil, err
	}
	//4字节 msgid
	var buffer bytes.Buffer
	err = binary.Write(&buffer, binary.BigEndian, messageElement.MsgId)
	if nil != err {
		return nil, err
	}
	//4字节 body length
	var bodyLength int32 = int32(len(sendBuffer))
	err = binary.Write(&buffer, binary.BigEndian, bodyLength)
	if nil != err {
		return nil, err
	}
	// body
	err = binary.Write(&buffer, binary.BigEndian, sendBuffer)
	if nil != err {
		return nil, err
	}

	return sendBuffer, nil
}
func (self *tcpPipeline) onReceived(buffer []byte) error {
	self.receivedBuffer = append(self.receivedBuffer, buffer...)
	for {
		switch self.decodeStatus {
		case decodeStatus_Id:
			if len(self.receivedBuffer) < 4 {
				return nil
			}
			// read msg id
			buffer := bytes.NewBuffer(self.receivedBuffer[0:4])
			var msgId int32
			err := binary.Read(buffer, binary.BigEndian, &msgId)
			if nil != err {
				return err
			}
			self.currentMsgId = msgId
			self.receivedBuffer = self.receivedBuffer[4:]
			self.decodeStatus = decodeStatus_BodyLength
		case decodeStatus_BodyLength:
			if len(self.receivedBuffer) < 4 {
				return nil
			}
			// read body length
			buffer := bytes.NewBuffer(self.receivedBuffer[0:4])
			var bodyLength int32
			err := binary.Read(buffer, binary.BigEndian, &bodyLength)
			if nil != err {
				return err
			}
			self.currentMsgBodyLength = bodyLength
			if self.currentMsgBodyLength <= 0 {
				return errors.New(fmt.Sprint("error on parser body length", self.currentMsgBodyLength))
			}
			self.receivedBuffer = self.receivedBuffer[4:]
			self.decodeStatus = decodeStatus_Body
		case decodeStatus_Body:
			if int32(len(self.receivedBuffer)) < self.currentMsgBodyLength {
				return nil
			}
			// read body
			buffer := self.receivedBuffer[:self.currentMsgBodyLength]
			msgBody, err := self.codeC.Decode(self.currentMsgId, buffer)
			if err != nil {
				return err
			}
			self.receivedQueue.Add(&Pipeline.PipelineMessageElement{MsgId: self.currentMsgId, MsgBody: msgBody, Session: self})
			self.receivedBuffer = self.receivedBuffer[self.currentMsgBodyLength:]
			self.decodeStatus = decodeStatus_Id
		}
	}
	return nil
}
