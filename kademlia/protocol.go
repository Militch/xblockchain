package kademlia

import (
	"xblockchain/uint256"
)

const TypeFlagReqHello = 1
const TypeFlagAckHello = 1 << 1
type Protocol struct {
	Header *ProtocolHeader
	Data []byte
}

type ProtocolHeader struct {
	SrcID *uint256.UInt256
	DestID  *uint256.UInt256
	Version uint8
	TypeFlag uint8
}

func NewProtocol(byteData []byte) *Protocol {
	sourceIdBytes := byteData[:32]
	destIdBytes := byteData[32:64]
	versionByte := byteData[64:65][0]
	typeFlagByte := byteData[65:66][0]
	version := versionByte & 0xff
	typeFlag := typeFlagByte & 0xff
	sourceId := uint256.NewUInt256BS(sourceIdBytes)
	destId := uint256.NewUInt256BS(destIdBytes)
	header := &ProtocolHeader{
		SrcID: sourceId,
		DestID: destId,
		Version: version,
		TypeFlag: typeFlag,
	}
	return &Protocol{
		Header: header,
	}
}

func (protocol *Protocol) ToBytes() []byte {
	header := protocol.Header
	versionBytes := []byte{header.Version}
	typeFlagBytes := []byte{header.TypeFlag}
	headerBytes := append(header.SrcID.ToBytes(), header.DestID.ToBytes()...)
	headerBytes = append(headerBytes,versionBytes...)
	headerBytes = append(headerBytes,typeFlagBytes...)
	return headerBytes
}


func (protocol *Protocol) HasTypeFlag(typeFlag uint8) bool {
	header := protocol.Header
	return header.TypeFlag & typeFlag != 0
}