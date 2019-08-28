package dump

import (
	"github.com/bradleyjkemp/grpc-tools/internal/dump_format"
	"google.golang.org/grpc"
	"sync"
	"time"
)

// recordedServerStream wraps a grpc.ServerStream and allows the dump interceptor to record all sent/received messages
type recordedServerStream struct {
	sync.Mutex
	grpc.ServerStream
	events []*dump_format.Message
}

func (ss *recordedServerStream) SendMsg(m interface{}) error {
	message := m.([]byte)
	ss.Lock()
	ss.events = append(ss.events, &dump_format.Message{
		MessageOrigin: dump_format.ServerMessage,
		RawMessage:    message,
		Timestamp:     time.Now(),
	})
	ss.Unlock()
	return ss.ServerStream.SendMsg(m)
}

func (ss *recordedServerStream) RecvMsg(m interface{}) error {
	err := ss.ServerStream.RecvMsg(m)
	if err != nil {
		return err
	}
	// now m is populated
	message := m.(*[]byte)
	ss.Lock()
	ss.events = append(ss.events, &dump_format.Message{
		MessageOrigin: dump_format.ClientMessage,
		RawMessage:    *message,
		Timestamp:     time.Now(),
	})
	ss.Unlock()
	return nil
}
