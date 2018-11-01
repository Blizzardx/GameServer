package Pipeline

type Pipeline interface {
	Send(msgId int32, msgBody interface{})
	Close()
}

type PipelineMessageElement struct {
	MsgId   int32
	MsgBody interface{}
	Session Pipeline
}
