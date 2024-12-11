// package: liara
// file: eventsource.proto

/* tslint:disable */
/* eslint-disable */

import * as grpc from "@grpc/grpc-js";
import * as eventsource_pb from "./eventsource_pb";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";

interface IEventSourceServiceService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
    append: IEventSourceServiceService_IAppend;
    testIdempotency: IEventSourceServiceService_ITestIdempotency;
    get: IEventSourceServiceService_IGet;
    getByAggregateIDAndName: IEventSourceServiceService_IGetByAggregateIDAndName;
    getAfterGlobalVersion: IEventSourceServiceService_IGetAfterGlobalVersion;
    getByOutbox: IEventSourceServiceService_IGetByOutbox;
    createOutbox: IEventSourceServiceService_ICreateOutbox;
    getOutbox: IEventSourceServiceService_IGetOutbox;
    updateOutboxPosition: IEventSourceServiceService_IUpdateOutboxPosition;
    createTenant: IEventSourceServiceService_ICreateTenant;
    deleteTenant: IEventSourceServiceService_IDeleteTenant;
    renameTenant: IEventSourceServiceService_IRenameTenant;
    getTenant: IEventSourceServiceService_IGetTenant;
    listTenants: IEventSourceServiceService_IListTenants;
}

interface IEventSourceServiceService_IAppend extends grpc.MethodDefinition<eventsource_pb.AppendRequest, eventsource_pb.AppendResponse> {
    path: "/liara.EventSourceService/Append";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<eventsource_pb.AppendRequest>;
    requestDeserialize: grpc.deserialize<eventsource_pb.AppendRequest>;
    responseSerialize: grpc.serialize<eventsource_pb.AppendResponse>;
    responseDeserialize: grpc.deserialize<eventsource_pb.AppendResponse>;
}
interface IEventSourceServiceService_ITestIdempotency extends grpc.MethodDefinition<eventsource_pb.TestIdempotencyRequest, eventsource_pb.TestIdempotencyResponse> {
    path: "/liara.EventSourceService/TestIdempotency";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<eventsource_pb.TestIdempotencyRequest>;
    requestDeserialize: grpc.deserialize<eventsource_pb.TestIdempotencyRequest>;
    responseSerialize: grpc.serialize<eventsource_pb.TestIdempotencyResponse>;
    responseDeserialize: grpc.deserialize<eventsource_pb.TestIdempotencyResponse>;
}
interface IEventSourceServiceService_IGet extends grpc.MethodDefinition<eventsource_pb.GetRequest, eventsource_pb.Event> {
    path: "/liara.EventSourceService/Get";
    requestStream: false;
    responseStream: true;
    requestSerialize: grpc.serialize<eventsource_pb.GetRequest>;
    requestDeserialize: grpc.deserialize<eventsource_pb.GetRequest>;
    responseSerialize: grpc.serialize<eventsource_pb.Event>;
    responseDeserialize: grpc.deserialize<eventsource_pb.Event>;
}
interface IEventSourceServiceService_IGetByAggregateIDAndName extends grpc.MethodDefinition<eventsource_pb.GetByAggregateIDAndNameRequest, eventsource_pb.Event> {
    path: "/liara.EventSourceService/GetByAggregateIDAndName";
    requestStream: false;
    responseStream: true;
    requestSerialize: grpc.serialize<eventsource_pb.GetByAggregateIDAndNameRequest>;
    requestDeserialize: grpc.deserialize<eventsource_pb.GetByAggregateIDAndNameRequest>;
    responseSerialize: grpc.serialize<eventsource_pb.Event>;
    responseDeserialize: grpc.deserialize<eventsource_pb.Event>;
}
interface IEventSourceServiceService_IGetAfterGlobalVersion extends grpc.MethodDefinition<eventsource_pb.GetAfterGlobalVersionRequest, eventsource_pb.Event> {
    path: "/liara.EventSourceService/GetAfterGlobalVersion";
    requestStream: false;
    responseStream: true;
    requestSerialize: grpc.serialize<eventsource_pb.GetAfterGlobalVersionRequest>;
    requestDeserialize: grpc.deserialize<eventsource_pb.GetAfterGlobalVersionRequest>;
    responseSerialize: grpc.serialize<eventsource_pb.Event>;
    responseDeserialize: grpc.deserialize<eventsource_pb.Event>;
}
interface IEventSourceServiceService_IGetByOutbox extends grpc.MethodDefinition<eventsource_pb.GetByOutboxRequest, eventsource_pb.Event> {
    path: "/liara.EventSourceService/GetByOutbox";
    requestStream: false;
    responseStream: true;
    requestSerialize: grpc.serialize<eventsource_pb.GetByOutboxRequest>;
    requestDeserialize: grpc.deserialize<eventsource_pb.GetByOutboxRequest>;
    responseSerialize: grpc.serialize<eventsource_pb.Event>;
    responseDeserialize: grpc.deserialize<eventsource_pb.Event>;
}
interface IEventSourceServiceService_ICreateOutbox extends grpc.MethodDefinition<eventsource_pb.CreateOutboxRequest, eventsource_pb.CreateOutboxResponse> {
    path: "/liara.EventSourceService/CreateOutbox";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<eventsource_pb.CreateOutboxRequest>;
    requestDeserialize: grpc.deserialize<eventsource_pb.CreateOutboxRequest>;
    responseSerialize: grpc.serialize<eventsource_pb.CreateOutboxResponse>;
    responseDeserialize: grpc.deserialize<eventsource_pb.CreateOutboxResponse>;
}
interface IEventSourceServiceService_IGetOutbox extends grpc.MethodDefinition<eventsource_pb.GetOutboxRequest, eventsource_pb.GetOutboxResponse> {
    path: "/liara.EventSourceService/GetOutbox";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<eventsource_pb.GetOutboxRequest>;
    requestDeserialize: grpc.deserialize<eventsource_pb.GetOutboxRequest>;
    responseSerialize: grpc.serialize<eventsource_pb.GetOutboxResponse>;
    responseDeserialize: grpc.deserialize<eventsource_pb.GetOutboxResponse>;
}
interface IEventSourceServiceService_IUpdateOutboxPosition extends grpc.MethodDefinition<eventsource_pb.UpdateOutboxPositionRequest, eventsource_pb.UpdateOutboxPositionResponse> {
    path: "/liara.EventSourceService/UpdateOutboxPosition";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<eventsource_pb.UpdateOutboxPositionRequest>;
    requestDeserialize: grpc.deserialize<eventsource_pb.UpdateOutboxPositionRequest>;
    responseSerialize: grpc.serialize<eventsource_pb.UpdateOutboxPositionResponse>;
    responseDeserialize: grpc.deserialize<eventsource_pb.UpdateOutboxPositionResponse>;
}
interface IEventSourceServiceService_ICreateTenant extends grpc.MethodDefinition<eventsource_pb.CreateTenantRequest, eventsource_pb.CreateTenantReponse> {
    path: "/liara.EventSourceService/CreateTenant";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<eventsource_pb.CreateTenantRequest>;
    requestDeserialize: grpc.deserialize<eventsource_pb.CreateTenantRequest>;
    responseSerialize: grpc.serialize<eventsource_pb.CreateTenantReponse>;
    responseDeserialize: grpc.deserialize<eventsource_pb.CreateTenantReponse>;
}
interface IEventSourceServiceService_IDeleteTenant extends grpc.MethodDefinition<eventsource_pb.DeleteTenantRequest, eventsource_pb.DeleteTenantResponse> {
    path: "/liara.EventSourceService/DeleteTenant";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<eventsource_pb.DeleteTenantRequest>;
    requestDeserialize: grpc.deserialize<eventsource_pb.DeleteTenantRequest>;
    responseSerialize: grpc.serialize<eventsource_pb.DeleteTenantResponse>;
    responseDeserialize: grpc.deserialize<eventsource_pb.DeleteTenantResponse>;
}
interface IEventSourceServiceService_IRenameTenant extends grpc.MethodDefinition<eventsource_pb.RenameTenantRequest, eventsource_pb.RenameTenantResponse> {
    path: "/liara.EventSourceService/RenameTenant";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<eventsource_pb.RenameTenantRequest>;
    requestDeserialize: grpc.deserialize<eventsource_pb.RenameTenantRequest>;
    responseSerialize: grpc.serialize<eventsource_pb.RenameTenantResponse>;
    responseDeserialize: grpc.deserialize<eventsource_pb.RenameTenantResponse>;
}
interface IEventSourceServiceService_IGetTenant extends grpc.MethodDefinition<eventsource_pb.GetTenantRequest, eventsource_pb.GetTenantResponse> {
    path: "/liara.EventSourceService/GetTenant";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<eventsource_pb.GetTenantRequest>;
    requestDeserialize: grpc.deserialize<eventsource_pb.GetTenantRequest>;
    responseSerialize: grpc.serialize<eventsource_pb.GetTenantResponse>;
    responseDeserialize: grpc.deserialize<eventsource_pb.GetTenantResponse>;
}
interface IEventSourceServiceService_IListTenants extends grpc.MethodDefinition<eventsource_pb.ListTenantsRequest, eventsource_pb.Tenant> {
    path: "/liara.EventSourceService/ListTenants";
    requestStream: false;
    responseStream: true;
    requestSerialize: grpc.serialize<eventsource_pb.ListTenantsRequest>;
    requestDeserialize: grpc.deserialize<eventsource_pb.ListTenantsRequest>;
    responseSerialize: grpc.serialize<eventsource_pb.Tenant>;
    responseDeserialize: grpc.deserialize<eventsource_pb.Tenant>;
}

export const EventSourceServiceService: IEventSourceServiceService;

export interface IEventSourceServiceServer extends grpc.UntypedServiceImplementation {
    append: grpc.handleUnaryCall<eventsource_pb.AppendRequest, eventsource_pb.AppendResponse>;
    testIdempotency: grpc.handleUnaryCall<eventsource_pb.TestIdempotencyRequest, eventsource_pb.TestIdempotencyResponse>;
    get: grpc.handleServerStreamingCall<eventsource_pb.GetRequest, eventsource_pb.Event>;
    getByAggregateIDAndName: grpc.handleServerStreamingCall<eventsource_pb.GetByAggregateIDAndNameRequest, eventsource_pb.Event>;
    getAfterGlobalVersion: grpc.handleServerStreamingCall<eventsource_pb.GetAfterGlobalVersionRequest, eventsource_pb.Event>;
    getByOutbox: grpc.handleServerStreamingCall<eventsource_pb.GetByOutboxRequest, eventsource_pb.Event>;
    createOutbox: grpc.handleUnaryCall<eventsource_pb.CreateOutboxRequest, eventsource_pb.CreateOutboxResponse>;
    getOutbox: grpc.handleUnaryCall<eventsource_pb.GetOutboxRequest, eventsource_pb.GetOutboxResponse>;
    updateOutboxPosition: grpc.handleUnaryCall<eventsource_pb.UpdateOutboxPositionRequest, eventsource_pb.UpdateOutboxPositionResponse>;
    createTenant: grpc.handleUnaryCall<eventsource_pb.CreateTenantRequest, eventsource_pb.CreateTenantReponse>;
    deleteTenant: grpc.handleUnaryCall<eventsource_pb.DeleteTenantRequest, eventsource_pb.DeleteTenantResponse>;
    renameTenant: grpc.handleUnaryCall<eventsource_pb.RenameTenantRequest, eventsource_pb.RenameTenantResponse>;
    getTenant: grpc.handleUnaryCall<eventsource_pb.GetTenantRequest, eventsource_pb.GetTenantResponse>;
    listTenants: grpc.handleServerStreamingCall<eventsource_pb.ListTenantsRequest, eventsource_pb.Tenant>;
}

export interface IEventSourceServiceClient {
    append(request: eventsource_pb.AppendRequest, callback: (error: grpc.ServiceError | null, response: eventsource_pb.AppendResponse) => void): grpc.ClientUnaryCall;
    append(request: eventsource_pb.AppendRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: eventsource_pb.AppendResponse) => void): grpc.ClientUnaryCall;
    append(request: eventsource_pb.AppendRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: eventsource_pb.AppendResponse) => void): grpc.ClientUnaryCall;
    testIdempotency(request: eventsource_pb.TestIdempotencyRequest, callback: (error: grpc.ServiceError | null, response: eventsource_pb.TestIdempotencyResponse) => void): grpc.ClientUnaryCall;
    testIdempotency(request: eventsource_pb.TestIdempotencyRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: eventsource_pb.TestIdempotencyResponse) => void): grpc.ClientUnaryCall;
    testIdempotency(request: eventsource_pb.TestIdempotencyRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: eventsource_pb.TestIdempotencyResponse) => void): grpc.ClientUnaryCall;
    get(request: eventsource_pb.GetRequest, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<eventsource_pb.Event>;
    get(request: eventsource_pb.GetRequest, metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<eventsource_pb.Event>;
    getByAggregateIDAndName(request: eventsource_pb.GetByAggregateIDAndNameRequest, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<eventsource_pb.Event>;
    getByAggregateIDAndName(request: eventsource_pb.GetByAggregateIDAndNameRequest, metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<eventsource_pb.Event>;
    getAfterGlobalVersion(request: eventsource_pb.GetAfterGlobalVersionRequest, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<eventsource_pb.Event>;
    getAfterGlobalVersion(request: eventsource_pb.GetAfterGlobalVersionRequest, metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<eventsource_pb.Event>;
    getByOutbox(request: eventsource_pb.GetByOutboxRequest, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<eventsource_pb.Event>;
    getByOutbox(request: eventsource_pb.GetByOutboxRequest, metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<eventsource_pb.Event>;
    createOutbox(request: eventsource_pb.CreateOutboxRequest, callback: (error: grpc.ServiceError | null, response: eventsource_pb.CreateOutboxResponse) => void): grpc.ClientUnaryCall;
    createOutbox(request: eventsource_pb.CreateOutboxRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: eventsource_pb.CreateOutboxResponse) => void): grpc.ClientUnaryCall;
    createOutbox(request: eventsource_pb.CreateOutboxRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: eventsource_pb.CreateOutboxResponse) => void): grpc.ClientUnaryCall;
    getOutbox(request: eventsource_pb.GetOutboxRequest, callback: (error: grpc.ServiceError | null, response: eventsource_pb.GetOutboxResponse) => void): grpc.ClientUnaryCall;
    getOutbox(request: eventsource_pb.GetOutboxRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: eventsource_pb.GetOutboxResponse) => void): grpc.ClientUnaryCall;
    getOutbox(request: eventsource_pb.GetOutboxRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: eventsource_pb.GetOutboxResponse) => void): grpc.ClientUnaryCall;
    updateOutboxPosition(request: eventsource_pb.UpdateOutboxPositionRequest, callback: (error: grpc.ServiceError | null, response: eventsource_pb.UpdateOutboxPositionResponse) => void): grpc.ClientUnaryCall;
    updateOutboxPosition(request: eventsource_pb.UpdateOutboxPositionRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: eventsource_pb.UpdateOutboxPositionResponse) => void): grpc.ClientUnaryCall;
    updateOutboxPosition(request: eventsource_pb.UpdateOutboxPositionRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: eventsource_pb.UpdateOutboxPositionResponse) => void): grpc.ClientUnaryCall;
    createTenant(request: eventsource_pb.CreateTenantRequest, callback: (error: grpc.ServiceError | null, response: eventsource_pb.CreateTenantReponse) => void): grpc.ClientUnaryCall;
    createTenant(request: eventsource_pb.CreateTenantRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: eventsource_pb.CreateTenantReponse) => void): grpc.ClientUnaryCall;
    createTenant(request: eventsource_pb.CreateTenantRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: eventsource_pb.CreateTenantReponse) => void): grpc.ClientUnaryCall;
    deleteTenant(request: eventsource_pb.DeleteTenantRequest, callback: (error: grpc.ServiceError | null, response: eventsource_pb.DeleteTenantResponse) => void): grpc.ClientUnaryCall;
    deleteTenant(request: eventsource_pb.DeleteTenantRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: eventsource_pb.DeleteTenantResponse) => void): grpc.ClientUnaryCall;
    deleteTenant(request: eventsource_pb.DeleteTenantRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: eventsource_pb.DeleteTenantResponse) => void): grpc.ClientUnaryCall;
    renameTenant(request: eventsource_pb.RenameTenantRequest, callback: (error: grpc.ServiceError | null, response: eventsource_pb.RenameTenantResponse) => void): grpc.ClientUnaryCall;
    renameTenant(request: eventsource_pb.RenameTenantRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: eventsource_pb.RenameTenantResponse) => void): grpc.ClientUnaryCall;
    renameTenant(request: eventsource_pb.RenameTenantRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: eventsource_pb.RenameTenantResponse) => void): grpc.ClientUnaryCall;
    getTenant(request: eventsource_pb.GetTenantRequest, callback: (error: grpc.ServiceError | null, response: eventsource_pb.GetTenantResponse) => void): grpc.ClientUnaryCall;
    getTenant(request: eventsource_pb.GetTenantRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: eventsource_pb.GetTenantResponse) => void): grpc.ClientUnaryCall;
    getTenant(request: eventsource_pb.GetTenantRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: eventsource_pb.GetTenantResponse) => void): grpc.ClientUnaryCall;
    listTenants(request: eventsource_pb.ListTenantsRequest, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<eventsource_pb.Tenant>;
    listTenants(request: eventsource_pb.ListTenantsRequest, metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<eventsource_pb.Tenant>;
}

export class EventSourceServiceClient extends grpc.Client implements IEventSourceServiceClient {
    constructor(address: string, credentials: grpc.ChannelCredentials, options?: Partial<grpc.ClientOptions>);
    public append(request: eventsource_pb.AppendRequest, callback: (error: grpc.ServiceError | null, response: eventsource_pb.AppendResponse) => void): grpc.ClientUnaryCall;
    public append(request: eventsource_pb.AppendRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: eventsource_pb.AppendResponse) => void): grpc.ClientUnaryCall;
    public append(request: eventsource_pb.AppendRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: eventsource_pb.AppendResponse) => void): grpc.ClientUnaryCall;
    public testIdempotency(request: eventsource_pb.TestIdempotencyRequest, callback: (error: grpc.ServiceError | null, response: eventsource_pb.TestIdempotencyResponse) => void): grpc.ClientUnaryCall;
    public testIdempotency(request: eventsource_pb.TestIdempotencyRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: eventsource_pb.TestIdempotencyResponse) => void): grpc.ClientUnaryCall;
    public testIdempotency(request: eventsource_pb.TestIdempotencyRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: eventsource_pb.TestIdempotencyResponse) => void): grpc.ClientUnaryCall;
    public get(request: eventsource_pb.GetRequest, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<eventsource_pb.Event>;
    public get(request: eventsource_pb.GetRequest, metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<eventsource_pb.Event>;
    public getByAggregateIDAndName(request: eventsource_pb.GetByAggregateIDAndNameRequest, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<eventsource_pb.Event>;
    public getByAggregateIDAndName(request: eventsource_pb.GetByAggregateIDAndNameRequest, metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<eventsource_pb.Event>;
    public getAfterGlobalVersion(request: eventsource_pb.GetAfterGlobalVersionRequest, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<eventsource_pb.Event>;
    public getAfterGlobalVersion(request: eventsource_pb.GetAfterGlobalVersionRequest, metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<eventsource_pb.Event>;
    public getByOutbox(request: eventsource_pb.GetByOutboxRequest, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<eventsource_pb.Event>;
    public getByOutbox(request: eventsource_pb.GetByOutboxRequest, metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<eventsource_pb.Event>;
    public createOutbox(request: eventsource_pb.CreateOutboxRequest, callback: (error: grpc.ServiceError | null, response: eventsource_pb.CreateOutboxResponse) => void): grpc.ClientUnaryCall;
    public createOutbox(request: eventsource_pb.CreateOutboxRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: eventsource_pb.CreateOutboxResponse) => void): grpc.ClientUnaryCall;
    public createOutbox(request: eventsource_pb.CreateOutboxRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: eventsource_pb.CreateOutboxResponse) => void): grpc.ClientUnaryCall;
    public getOutbox(request: eventsource_pb.GetOutboxRequest, callback: (error: grpc.ServiceError | null, response: eventsource_pb.GetOutboxResponse) => void): grpc.ClientUnaryCall;
    public getOutbox(request: eventsource_pb.GetOutboxRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: eventsource_pb.GetOutboxResponse) => void): grpc.ClientUnaryCall;
    public getOutbox(request: eventsource_pb.GetOutboxRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: eventsource_pb.GetOutboxResponse) => void): grpc.ClientUnaryCall;
    public updateOutboxPosition(request: eventsource_pb.UpdateOutboxPositionRequest, callback: (error: grpc.ServiceError | null, response: eventsource_pb.UpdateOutboxPositionResponse) => void): grpc.ClientUnaryCall;
    public updateOutboxPosition(request: eventsource_pb.UpdateOutboxPositionRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: eventsource_pb.UpdateOutboxPositionResponse) => void): grpc.ClientUnaryCall;
    public updateOutboxPosition(request: eventsource_pb.UpdateOutboxPositionRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: eventsource_pb.UpdateOutboxPositionResponse) => void): grpc.ClientUnaryCall;
    public createTenant(request: eventsource_pb.CreateTenantRequest, callback: (error: grpc.ServiceError | null, response: eventsource_pb.CreateTenantReponse) => void): grpc.ClientUnaryCall;
    public createTenant(request: eventsource_pb.CreateTenantRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: eventsource_pb.CreateTenantReponse) => void): grpc.ClientUnaryCall;
    public createTenant(request: eventsource_pb.CreateTenantRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: eventsource_pb.CreateTenantReponse) => void): grpc.ClientUnaryCall;
    public deleteTenant(request: eventsource_pb.DeleteTenantRequest, callback: (error: grpc.ServiceError | null, response: eventsource_pb.DeleteTenantResponse) => void): grpc.ClientUnaryCall;
    public deleteTenant(request: eventsource_pb.DeleteTenantRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: eventsource_pb.DeleteTenantResponse) => void): grpc.ClientUnaryCall;
    public deleteTenant(request: eventsource_pb.DeleteTenantRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: eventsource_pb.DeleteTenantResponse) => void): grpc.ClientUnaryCall;
    public renameTenant(request: eventsource_pb.RenameTenantRequest, callback: (error: grpc.ServiceError | null, response: eventsource_pb.RenameTenantResponse) => void): grpc.ClientUnaryCall;
    public renameTenant(request: eventsource_pb.RenameTenantRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: eventsource_pb.RenameTenantResponse) => void): grpc.ClientUnaryCall;
    public renameTenant(request: eventsource_pb.RenameTenantRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: eventsource_pb.RenameTenantResponse) => void): grpc.ClientUnaryCall;
    public getTenant(request: eventsource_pb.GetTenantRequest, callback: (error: grpc.ServiceError | null, response: eventsource_pb.GetTenantResponse) => void): grpc.ClientUnaryCall;
    public getTenant(request: eventsource_pb.GetTenantRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: eventsource_pb.GetTenantResponse) => void): grpc.ClientUnaryCall;
    public getTenant(request: eventsource_pb.GetTenantRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: eventsource_pb.GetTenantResponse) => void): grpc.ClientUnaryCall;
    public listTenants(request: eventsource_pb.ListTenantsRequest, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<eventsource_pb.Tenant>;
    public listTenants(request: eventsource_pb.ListTenantsRequest, metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<eventsource_pb.Tenant>;
}
