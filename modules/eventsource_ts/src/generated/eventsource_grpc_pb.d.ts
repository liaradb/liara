// package: todo
// file: eventsource.proto

/* tslint:disable */
/* eslint-disable */

import * as grpc from "@grpc/grpc-js";
import * as eventsource_pb from "./eventsource_pb";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";

interface IEventSourceServiceService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
    append: IEventSourceServiceService_IAppend;
    get: IEventSourceServiceService_IGet;
    getByAggregateIDAndName: IEventSourceServiceService_IGetByAggregateIDAndName;
}

interface IEventSourceServiceService_IAppend extends grpc.MethodDefinition<eventsource_pb.AppendRequest, eventsource_pb.AppendResponse> {
    path: "/todo.EventSourceService/Append";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<eventsource_pb.AppendRequest>;
    requestDeserialize: grpc.deserialize<eventsource_pb.AppendRequest>;
    responseSerialize: grpc.serialize<eventsource_pb.AppendResponse>;
    responseDeserialize: grpc.deserialize<eventsource_pb.AppendResponse>;
}
interface IEventSourceServiceService_IGet extends grpc.MethodDefinition<eventsource_pb.GetRequest, eventsource_pb.Event> {
    path: "/todo.EventSourceService/Get";
    requestStream: false;
    responseStream: true;
    requestSerialize: grpc.serialize<eventsource_pb.GetRequest>;
    requestDeserialize: grpc.deserialize<eventsource_pb.GetRequest>;
    responseSerialize: grpc.serialize<eventsource_pb.Event>;
    responseDeserialize: grpc.deserialize<eventsource_pb.Event>;
}
interface IEventSourceServiceService_IGetByAggregateIDAndName extends grpc.MethodDefinition<eventsource_pb.GetByAggregateIDAndNameRequest, eventsource_pb.Event> {
    path: "/todo.EventSourceService/GetByAggregateIDAndName";
    requestStream: false;
    responseStream: true;
    requestSerialize: grpc.serialize<eventsource_pb.GetByAggregateIDAndNameRequest>;
    requestDeserialize: grpc.deserialize<eventsource_pb.GetByAggregateIDAndNameRequest>;
    responseSerialize: grpc.serialize<eventsource_pb.Event>;
    responseDeserialize: grpc.deserialize<eventsource_pb.Event>;
}

export const EventSourceServiceService: IEventSourceServiceService;

export interface IEventSourceServiceServer extends grpc.UntypedServiceImplementation {
    append: grpc.handleUnaryCall<eventsource_pb.AppendRequest, eventsource_pb.AppendResponse>;
    get: grpc.handleServerStreamingCall<eventsource_pb.GetRequest, eventsource_pb.Event>;
    getByAggregateIDAndName: grpc.handleServerStreamingCall<eventsource_pb.GetByAggregateIDAndNameRequest, eventsource_pb.Event>;
}

export interface IEventSourceServiceClient {
    append(request: eventsource_pb.AppendRequest, callback: (error: grpc.ServiceError | null, response: eventsource_pb.AppendResponse) => void): grpc.ClientUnaryCall;
    append(request: eventsource_pb.AppendRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: eventsource_pb.AppendResponse) => void): grpc.ClientUnaryCall;
    append(request: eventsource_pb.AppendRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: eventsource_pb.AppendResponse) => void): grpc.ClientUnaryCall;
    get(request: eventsource_pb.GetRequest, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<eventsource_pb.Event>;
    get(request: eventsource_pb.GetRequest, metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<eventsource_pb.Event>;
    getByAggregateIDAndName(request: eventsource_pb.GetByAggregateIDAndNameRequest, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<eventsource_pb.Event>;
    getByAggregateIDAndName(request: eventsource_pb.GetByAggregateIDAndNameRequest, metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<eventsource_pb.Event>;
}

export class EventSourceServiceClient extends grpc.Client implements IEventSourceServiceClient {
    constructor(address: string, credentials: grpc.ChannelCredentials, options?: Partial<grpc.ClientOptions>);
    public append(request: eventsource_pb.AppendRequest, callback: (error: grpc.ServiceError | null, response: eventsource_pb.AppendResponse) => void): grpc.ClientUnaryCall;
    public append(request: eventsource_pb.AppendRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: eventsource_pb.AppendResponse) => void): grpc.ClientUnaryCall;
    public append(request: eventsource_pb.AppendRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: eventsource_pb.AppendResponse) => void): grpc.ClientUnaryCall;
    public get(request: eventsource_pb.GetRequest, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<eventsource_pb.Event>;
    public get(request: eventsource_pb.GetRequest, metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<eventsource_pb.Event>;
    public getByAggregateIDAndName(request: eventsource_pb.GetByAggregateIDAndNameRequest, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<eventsource_pb.Event>;
    public getByAggregateIDAndName(request: eventsource_pb.GetByAggregateIDAndNameRequest, metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<eventsource_pb.Event>;
}
