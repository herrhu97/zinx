package znet

import "github.com/herrhu97/zinx/ziface"

type Message struct {
	DataLen uint32
	MsgID   uint32
	Data    []byte
}

func NewMsgPackage(id uint32, data []byte) ziface.IMessage {
	m := &Message{
		DataLen: uint32(len(data)),
		MsgID:   id,
		Data:    data,
	}
	return m
}

func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

func (m *Message) GetMsgID() uint32 {
	return m.MsgID
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}

func (m *Message) SetDataLen(l uint32) {
	m.DataLen = l
}

func (m *Message) SetMsgID(id uint32) {
	m.MsgID = id
}
