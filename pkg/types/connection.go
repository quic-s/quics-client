package types

type MessageData interface {
	Encode() []byte
	Decode([]byte)
}
