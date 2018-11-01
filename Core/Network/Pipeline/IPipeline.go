package Pipeline

type Pipeline interface {
	Send(msgId int32, msgBody interface{})
	Close()
	SetParameter(arg interface{})
	GetParameter() interface{}
}

type PipelineMessageElement struct {
	MsgId   int32
	MsgBody interface{}
	Session Pipeline
}
