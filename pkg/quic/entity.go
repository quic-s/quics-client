package quic

type MessageData interface {
	Encode() []byte
	Decode([]byte)
}
