package esphomeproto

import "google.golang.org/protobuf/proto"

type EsphomeMessageTyper interface {
	proto.Message
	EsphomeMessageType() int
	EsphomeSource() APISourceType
}

type MessageFactory func() EsphomeMessageTyper

func messageFactory[T any, PT interface {
	EsphomeMessageTyper
	*T
}]() MessageFactory {
	return func() EsphomeMessageTyper {
		return PT(new(T))
	}
}

func MessageByType(t int) EsphomeMessageTyper {
	c, ok := MessageTypeToType[t]
	if !ok {
		return nil
	}
	return c()
}
