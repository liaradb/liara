// package: todo
// file: eventsource.proto

/* tslint:disable */
/* eslint-disable */

import * as grpc from "@grpc/grpc-js";
import * as eventsource_pb from "./eventsource_pb";

interface IEventSourceServiceService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
    listEvents: IEventSourceServiceService_IListEvents;
}

interface IEventSourceServiceService_IListEvents extends grpc.MethodDefinition<eventsource_pb.ListEventsRequest, eventsource_pb.ListEventsResponse> {
    path: "/todo.EventSourceService/ListEvents";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<eventsource_pb.ListEventsRequest>;
    requestDeserialize: grpc.deserialize<eventsource_pb.ListEventsRequest>;
    responseSerialize: grpc.serialize<eventsource_pb.ListEventsResponse>;
    responseDeserialize: grpc.deserialize<eventsource_pb.ListEventsResponse>;
}

export const EventSourceServiceService: IEventSourceServiceService;

export interface IEventSourceServiceServer extends grpc.UntypedServiceImplementation {
    listEvents: grpc.handleUnaryCall<eventsource_pb.ListEventsRequest, eventsource_pb.ListEventsResponse>;
}

export interface IEventSourceServiceClient {
    listEvents(request: eventsource_pb.ListEventsRequest, callback: (error: grpc.ServiceError | null, response: eventsource_pb.ListEventsResponse) => void): grpc.ClientUnaryCall;
    listEvents(request: eventsource_pb.ListEventsRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: eventsource_pb.ListEventsResponse) => void): grpc.ClientUnaryCall;
    listEvents(request: eventsource_pb.ListEventsRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: eventsource_pb.ListEventsResponse) => void): grpc.ClientUnaryCall;
}

export class EventSourceServiceClient extends grpc.Client implements IEventSourceServiceClient {
    constructor(address: string, credentials: grpc.ChannelCredentials, options?: Partial<grpc.ClientOptions>);
    public listEvents(request: eventsource_pb.ListEventsRequest, callback: (error: grpc.ServiceError | null, response: eventsource_pb.ListEventsResponse) => void): grpc.ClientUnaryCall;
    public listEvents(request: eventsource_pb.ListEventsRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: eventsource_pb.ListEventsResponse) => void): grpc.ClientUnaryCall;
    public listEvents(request: eventsource_pb.ListEventsRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: eventsource_pb.ListEventsResponse) => void): grpc.ClientUnaryCall;
}
