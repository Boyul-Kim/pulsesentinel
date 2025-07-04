// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: proto/sentinel/event.proto

package sentinel

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	EventIngestor_StreamEvents_FullMethodName = "/sentinel.EventIngestor/StreamEvents"
)

// EventIngestorClient is the client API for EventIngestor service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// gRPC Service
type EventIngestorClient interface {
	StreamEvents(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[SecurityEvent, EventResponse], error)
}

type eventIngestorClient struct {
	cc grpc.ClientConnInterface
}

func NewEventIngestorClient(cc grpc.ClientConnInterface) EventIngestorClient {
	return &eventIngestorClient{cc}
}

func (c *eventIngestorClient) StreamEvents(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[SecurityEvent, EventResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &EventIngestor_ServiceDesc.Streams[0], EventIngestor_StreamEvents_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[SecurityEvent, EventResponse]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type EventIngestor_StreamEventsClient = grpc.BidiStreamingClient[SecurityEvent, EventResponse]

// EventIngestorServer is the server API for EventIngestor service.
// All implementations must embed UnimplementedEventIngestorServer
// for forward compatibility.
//
// gRPC Service
type EventIngestorServer interface {
	StreamEvents(grpc.BidiStreamingServer[SecurityEvent, EventResponse]) error
	mustEmbedUnimplementedEventIngestorServer()
}

// UnimplementedEventIngestorServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedEventIngestorServer struct{}

func (UnimplementedEventIngestorServer) StreamEvents(grpc.BidiStreamingServer[SecurityEvent, EventResponse]) error {
	return status.Errorf(codes.Unimplemented, "method StreamEvents not implemented")
}
func (UnimplementedEventIngestorServer) mustEmbedUnimplementedEventIngestorServer() {}
func (UnimplementedEventIngestorServer) testEmbeddedByValue()                       {}

// UnsafeEventIngestorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EventIngestorServer will
// result in compilation errors.
type UnsafeEventIngestorServer interface {
	mustEmbedUnimplementedEventIngestorServer()
}

func RegisterEventIngestorServer(s grpc.ServiceRegistrar, srv EventIngestorServer) {
	// If the following call pancis, it indicates UnimplementedEventIngestorServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&EventIngestor_ServiceDesc, srv)
}

func _EventIngestor_StreamEvents_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(EventIngestorServer).StreamEvents(&grpc.GenericServerStream[SecurityEvent, EventResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type EventIngestor_StreamEventsServer = grpc.BidiStreamingServer[SecurityEvent, EventResponse]

// EventIngestor_ServiceDesc is the grpc.ServiceDesc for EventIngestor service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EventIngestor_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sentinel.EventIngestor",
	HandlerType: (*EventIngestorServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamEvents",
			Handler:       _EventIngestor_StreamEvents_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/sentinel/event.proto",
}
