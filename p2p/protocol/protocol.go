package protocol

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
)


type Type uint32

var HeaderLength = 12
var VersionT Type = 0x01
var VersionAckT Type = 0x02

type Header struct {
	Length  uint32
	Version uint32
	Type    Type
}

type Message struct {
	Header Header
	Body   []byte
}

func NewMessageString(t Type, version uint32, data string) *Message {
	dataLength := HeaderLength + len(data)
	header := &Header{
		Type: t,
		Version: version,
		Length: uint32(dataLength),
	}
	dataBytes := []byte(data)
	return &Message{
		Header: *header,
		Body: dataBytes,
	}
}

func NewMessage(t Type, version uint32, data []byte) *Message {
	dataLength := HeaderLength + len(data)
	header := &Header{
		Type: t,
		Version: version,
		Length: uint32(dataLength),
	}
	return &Message{
		Header: *header,
		Body: data,
	}
}

func (p Message) Marshal() []byte {
	var data []byte = nil
	hbs,err := p.Header.Marshal()
	if err != nil {
		return nil
	}
	data = append(hbs, p.Body...)
	return data
}

func UnmarshalMessage(bs []byte) (*Message,error) {
	header,err := UnmarshalHeader(bs)
	if err != nil {
		return nil, err
	}
	bodyLength := int(header.Length) - HeaderLength
	data := bs[HeaderLength : HeaderLength+bodyLength]
	return &Message{
		Header: header,
		Body: data,
	},nil
}

func UnmarshalHeader(bs []byte) (Header,error) {
	if len(bs) < HeaderLength {
		return Header{}, errors.New("illegal header structure")
	}
	header := Header{}
	header.Length = binary.BigEndian.Uint32(bs[0:4])
	header.Version = binary.BigEndian.Uint32(bs[4:8])
	t := binary.BigEndian.Uint32(bs[8:12])
	header.Type = Type(t)
	return header,nil
}

func (h *Header) Marshal() ([]byte,error) {
	bytesBuffer := bytes.NewBuffer([]byte{})
	err := binary.Write(bytesBuffer, binary.BigEndian, h.Length)
	if err != nil {
		return nil, err
	}
	err = binary.Write(bytesBuffer, binary.BigEndian, h.Version)
	if err != nil {
		return nil, err
	}
	err = binary.Write(bytesBuffer, binary.BigEndian, h.Type)
	if err != nil {
		return nil, err
	}
	return bytesBuffer.Bytes(), nil
}

func (h *Header) String() string {
	bs,err := json.Marshal(h)
	if err != nil {
		return ""
	}
	return string(bs)
}