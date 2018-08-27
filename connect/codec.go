package connect

import (
	"encoding/binary"
	"time"
	"net"
	"fmt"
)

const (
	TypeLen = 2         // 消息类型字节数组长度
	LenLen  = 2         // 消息长度字节数组长度
	HeadLen = 4         // 消息头部字节数组长度（消息类型字节数组长度+消息长度字节数组长度）
	BufLen  = 65536 + 4 // 缓冲buffer字节数组长度
)

type Codec struct {
	Conn     net.Conn
	ReadBuf  buffer // 读缓冲
	WriteBuf []byte // 写缓冲
}

// newCodec 创建一个解码器
func NewCodec(conn net.Conn) *Codec {
	return &Codec{
		Conn:     conn,
		ReadBuf:  newBuffer(conn, BufLen),
		WriteBuf: make([]byte, BufLen),
	}
}

// Read 从conn里面读取数据，当conn发生阻塞，这个方法也会阻塞
func (c *Codec) Read() (int, error) {
	return c.ReadBuf.readFromReader()
}

// Decode 解码数据
func (c *Codec) Decode() (*Message, bool) {
	var err error
	// 读取数据类型
	typeBuf, err := c.ReadBuf.seek(0, TypeLen)
	if err != nil {
		return nil, false
	}

	// 读取数据长度
	lenBuf, err := c.ReadBuf.seek(TypeLen, HeadLen)
	if err != nil {
		return nil, false
	}

	// 读取数据内容
	valueType := int(binary.BigEndian.Uint16(typeBuf))
	valueLen := int(binary.BigEndian.Uint16(lenBuf))

	valueBuf, err := c.ReadBuf.read(HeadLen, valueLen)
	if err != nil {
		return nil, false
	}
	message := Message{Code: valueType, Content: valueBuf}
	return &message, true
}

// Eecode 编码数据
func (c *Codec) Encode(message Message, duration time.Duration) error {
	contentLen := len(message.Content)

	binary.BigEndian.PutUint16(c.WriteBuf[0:TypeLen], uint16(message.Code))

	binary.BigEndian.PutUint16(c.WriteBuf[TypeLen:HeadLen], uint16(len(message.Content)))

	copy(c.WriteBuf[HeadLen:], message.Content[:contentLen])

	c.Conn.SetWriteDeadline(time.Now().Add(duration))

	_, err := c.Conn.Write(c.WriteBuf[:HeadLen+contentLen])

	fmt.Println(c.WriteBuf[:HeadLen+contentLen])
	//spew.Dump("消息 >>>>", c.WriteBuf[:HeadLen+contentLen])
	if err != nil {
		return err
	}
	return nil
}
