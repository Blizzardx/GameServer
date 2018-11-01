package Core

type NetworkCodeC interface{
	Decode(msgId int32,buffer []byte)(interface{},error)
	Encode(msgId int32,msgBody interface{})([]byte,error)
}