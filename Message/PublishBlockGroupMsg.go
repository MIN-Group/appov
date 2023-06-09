package Message

import (
	"fmt"
	"ppov/MetaData"
)

//go:generate msgp

type PublishBlockGroupMsg struct {
	Height   int                 `msg:"Height"`
	Group    MetaData.BlockGroup `msg:"-"`
	MinicNum int                 `msg:"MinicNum"`
	Data     []byte              `msg:"BlockGroup"`
}

func (msg *PublishBlockGroupMsg) ToByteArray() ([]byte, error) {
	temp_data, err := msg.Group.ToBytes(nil)
	msg.Data = temp_data
	data, err := msg.MarshalMsg(nil)
	return data, err
}
func (msg *PublishBlockGroupMsg) FromByteArray(data []byte) error {
	data, err := msg.UnmarshalMsg(data)
	data, err = msg.Group.FromBytes(msg.Data)
	if err != nil {
		fmt.Println("PublishBlockGroup-FromByteArray, err=", err)
	}
	return err
}

func (manager *MessagerManager) CreatePublishBlockGroupMsg(receiver uint64, height int, minicNum int) (header MessageHeader, msg PublishBlockGroupMsg) {
	msg.Height = height
	msg.MinicNum = minicNum
	header = manager.CreateHeader(receiver, PublishBlockGroup, 0, 0)
	return
}
