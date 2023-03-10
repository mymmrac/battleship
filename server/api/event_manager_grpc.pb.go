// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: event_manager.proto

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	EventManager_Events_FullMethodName = "/api.EventManager/Events"
)

// EventManagerClient is the client API for EventManager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EventManagerClient interface {
	Events(ctx context.Context, opts ...grpc.CallOption) (EventManager_EventsClient, error)
}

type eventManagerClient struct {
	cc grpc.ClientConnInterface
}

func NewEventManagerClient(cc grpc.ClientConnInterface) EventManagerClient {
	return &eventManagerClient{cc}
}

func (c *eventManagerClient) Events(ctx context.Context, opts ...grpc.CallOption) (EventManager_EventsClient, error) {
	stream, err := c.cc.NewStream(ctx, &EventManager_ServiceDesc.Streams[0], EventManager_Events_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &eventManagerEventsClient{stream}
	return x, nil
}

type EventManager_EventsClient interface {
	Send(*Event) error
	Recv() (*Event, error)
	grpc.ClientStream
}

type eventManagerEventsClient struct {
	grpc.ClientStream
}

func (x *eventManagerEventsClient) Send(m *Event) error {
	return x.ClientStream.SendMsg(m)
}

func (x *eventManagerEventsClient) Recv() (*Event, error) {
	m := new(Event)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// EventManagerServer is the server API for EventManager service.
// All implementations must embed UnimplementedEventManagerServer
// for forward compatibility
type EventManagerServer interface {
	Events(EventManager_EventsServer) error
	mustEmbedUnimplementedEventManagerServer()
}

// UnimplementedEventManagerServer must be embedded to have forward compatible implementations.
type UnimplementedEventManagerServer struct {
}

func (UnimplementedEventManagerServer) Events(EventManager_EventsServer) error {
	return status.Errorf(codes.Unimplemented, "method Events not implemented")
}
func (UnimplementedEventManagerServer) mustEmbedUnimplementedEventManagerServer() {}

// UnsafeEventManagerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EventManagerServer will
// result in compilation errors.
type UnsafeEventManagerServer interface {
	mustEmbedUnimplementedEventManagerServer()
}

func RegisterEventManagerServer(s grpc.ServiceRegistrar, srv EventManagerServer) {
	s.RegisterService(&EventManager_ServiceDesc, srv)
}

func _EventManager_Events_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(EventManagerServer).Events(&eventManagerEventsServer{stream})
}

type EventManager_EventsServer interface {
	Send(*Event) error
	Recv() (*Event, error)
	grpc.ServerStream
}

type eventManagerEventsServer struct {
	grpc.ServerStream
}

func (x *eventManagerEventsServer) Send(m *Event) error {
	return x.ServerStream.SendMsg(m)
}

func (x *eventManagerEventsServer) Recv() (*Event, error) {
	m := new(Event)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// EventManager_ServiceDesc is the grpc.ServiceDesc for EventManager service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EventManager_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.EventManager",
	HandlerType: (*EventManagerServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Events",
			Handler:       _EventManager_Events_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "event_manager.proto",
}
