package protocol

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	MaxBodySize = int32(1 << 12)
)

const (
	_packSize      = 4
	_headerSize    = 2
	_verSize       = 2
	_opSize        = 4
	_seqSize       = 4
	_heartSize     = 4
	_rawHeaderSize = _packSize + _headerSize + _verSize + _opSize + _seqSize
	_maxPackSize   = MaxBodySize + int32(_rawHeaderSize)
	_packOffset    = 0
	_headerOffset  = _packOffset + _packSize
	_verOffset     = _headerOffset + _headerSize
	_opOffset      = _verOffset + _verSize
	_seqOffset     = _opOffset + _opSize
	_heartOffset   = _seqOffset + _seqSize

	rawHeaderLen = uint16(16)
)

var (
	// ErrProtoPackLen proto packet len error
	ErrProtoPackLen = errors.New("default server codec pack length error")
	// ErrProtoHeaderLen proto header len error
	ErrProtoHeaderLen = errors.New("default server codec header length error")
)

// Proto 协议定义，以及协议编解码
type Proto struct {
	Ver  int32
	Op   uint32
	Seq  uint32
	Body []byte
}

// ReadTCP 借助io.ReadFull读取指定大小的包（该方法会一直读取到指定大小才返回，除非遇到读取错误或EOF
// 先读取Head大小的字节流，解出Head，再根据PackLen和HeadLen计算出Body大小，读取Body大上的字节流
func (p *Proto) ReadTCP(rr *bufio.Reader) (err error) {
	var (
		bodyLen   int
		headerLen uint16
		packLen   uint32
	)
	buf := make([]byte, _rawHeaderSize)

	if _, err = io.ReadFull(rr, buf); err != nil {
		fmt.Println("read error:", err)
		return
	}

	packLen = binary.BigEndian.Uint32(buf[_packOffset:_headerOffset])
	headerLen = binary.BigEndian.Uint16(buf[_headerOffset:_verOffset])
	p.Ver = int32(binary.BigEndian.Uint16(buf[_verOffset:_opOffset]))
	p.Op = binary.BigEndian.Uint32(buf[_opOffset:_seqOffset])
	p.Seq = binary.BigEndian.Uint32(buf[_seqOffset:])

	if packLen > uint32(_maxPackSize) {
		return ErrProtoPackLen
	}

	if headerLen != _rawHeaderSize {
		return ErrProtoHeaderLen
	}

	if bodyLen = int(packLen - uint32(headerLen)); bodyLen > 0 {
		p.Body = make([]byte, bodyLen)
		if _, err = io.ReadFull(rr, p.Body); err != nil {
			fmt.Println("read error:", err)
			return
		}
	} else {
		p.Body = nil
	}

	return
}

func (p *Proto) WriteTCP(wr *bufio.Writer) (err error) {

	if err = binary.Write(wr, binary.BigEndian, uint32(rawHeaderLen)+uint32(len(p.Body))); err != nil {
		return
	}
	if err = binary.Write(wr, binary.BigEndian, rawHeaderLen); err != nil {
		return
	}
	if err = binary.Write(wr, binary.BigEndian, p.Ver); err != nil {
		return
	}
	if err = binary.Write(wr, binary.BigEndian, p.Op); err != nil {
		return
	}
	if err = binary.Write(wr, binary.BigEndian, p.Seq); err != nil {
		return
	}
	if p.Body != nil {
		if err = binary.Write(wr, binary.BigEndian, p.Body); err != nil {
			return
		}
	}
	err = wr.Flush()
	return
}
