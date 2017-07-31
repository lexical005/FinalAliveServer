package tcpclient

import (
	"ffCommon/net/base"
	"ffProto"
	"fmt"
)

type clientNetEventData struct {
	eventType base.NetEventType

	data   base.NetEventData
	client *tcpClient
}

// Back 回收
func (c *clientNetEventData) Back() {
	eventType := c.NetEventType()
	if eventType == base.NetEventOff {
		c.client.onSessionClosed()
	} else if eventType == base.NetEventEnd {
		c.client.back()
	}
	c.client = nil

	// 回收data
	c.data.Back()

	// 回收自身
	eventDataPool.back(c)
}

// Client Client
func (c *clientNetEventData) Client() base.Client {
	return c.client
}

// NetEventType 获取事件类型
func (c *clientNetEventData) NetEventType() base.NetEventType {
	if c.eventType == base.NetEventInvalid {
		return c.data.NetEventType()
	}
	return c.eventType
}

// ManualClose 当NetEvent为NetEventOff时有效, 返回是不是主动关闭引发的Session断开
func (c *clientNetEventData) ManualClose() bool {
	return c.data.ManualClose()
}

// Proto 当NetEvent为NetEventProto时有效, 返回事件携带的协议
func (c *clientNetEventData) Proto() *ffProto.Proto {
	return c.data.Proto()
}

func (c *clientNetEventData) String() string {
	return fmt.Sprintf(`uuidClient[%v] dataSession[%v]`,
		c.client.uuid, c.data)
}

func newClientNetEventData() *clientNetEventData {
	return &clientNetEventData{}
}

func newClientNetEventDataFromSessionNetEventData(client *tcpClient, dataSession base.NetEventData) *clientNetEventData {
	dataClient := eventDataPool.apply()
	dataClient.client, dataClient.data, dataClient.eventType = client, dataSession, base.NetEventInvalid
	return dataClient
}

func newClientNetEventDataEnd(client *tcpClient) *clientNetEventData {
	dataClient := eventDataPool.apply()
	dataClient.client, dataClient.data, dataClient.eventType = client, nil, base.NetEventEnd
	return dataClient
}
