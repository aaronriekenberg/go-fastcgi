package request

import "sync/atomic"

type RequestID int64

type requestIDContextKey struct {
}

var RequestIDContextKey = &requestIDContextKey{}

type RequestIDFactory interface {
	NextRequestID() RequestID
}

type requestIDFactory struct {
	nextRequestID atomic.Int64
}

func (rf *requestIDFactory) NextRequestID() RequestID {
	return RequestID(rf.nextRequestID.Add(1))
}

var requestIDFactoryInstance = &requestIDFactory{}

func RequestIDFactoryInstance() RequestIDFactory {
	return requestIDFactoryInstance
}
