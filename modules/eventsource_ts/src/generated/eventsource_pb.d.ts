// package: todo
// file: eventsource.proto

/* tslint:disable */
/* eslint-disable */

import * as jspb from "google-protobuf";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";

export class Event extends jspb.Message { 
    getGlobalVersion(): number;
    setGlobalVersion(value: number): Event;
    getId(): string;
    setId(value: string): Event;
    getAggregateName(): string;
    setAggregateName(value: string): Event;
    getAggregateId(): string;
    setAggregateId(value: string): Event;
    getVersion(): number;
    setVersion(value: number): Event;
    getPartitionId(): number;
    setPartitionId(value: number): Event;
    getName(): string;
    setName(value: string): Event;
    getSchema(): string;
    setSchema(value: string): Event;

    hasMetadata(): boolean;
    clearMetadata(): void;
    getMetadata(): EventMetadata | undefined;
    setMetadata(value?: EventMetadata): Event;
    getData(): Uint8Array | string;
    getData_asU8(): Uint8Array;
    getData_asB64(): string;
    setData(value: Uint8Array | string): Event;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Event.AsObject;
    static toObject(includeInstance: boolean, msg: Event): Event.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Event, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Event;
    static deserializeBinaryFromReader(message: Event, reader: jspb.BinaryReader): Event;
}

export namespace Event {
    export type AsObject = {
        globalVersion: number,
        id: string,
        aggregateName: string,
        aggregateId: string,
        version: number,
        partitionId: number,
        name: string,
        schema: string,
        metadata?: EventMetadata.AsObject,
        data: Uint8Array | string,
    }
}

export class EventMetadata extends jspb.Message { 
    getCorrelationId(): string;
    setCorrelationId(value: string): EventMetadata;
    getUserId(): string;
    setUserId(value: string): EventMetadata;

    hasTime(): boolean;
    clearTime(): void;
    getTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setTime(value?: google_protobuf_timestamp_pb.Timestamp): EventMetadata;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): EventMetadata.AsObject;
    static toObject(includeInstance: boolean, msg: EventMetadata): EventMetadata.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: EventMetadata, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): EventMetadata;
    static deserializeBinaryFromReader(message: EventMetadata, reader: jspb.BinaryReader): EventMetadata;
}

export namespace EventMetadata {
    export type AsObject = {
        correlationId: string,
        userId: string,
        time?: google_protobuf_timestamp_pb.Timestamp.AsObject,
    }
}

export class AppendEvent extends jspb.Message { 
    getId(): string;
    setId(value: string): AppendEvent;
    getAggregateName(): string;
    setAggregateName(value: string): AppendEvent;
    getAggregateId(): string;
    setAggregateId(value: string): AppendEvent;
    getVersion(): number;
    setVersion(value: number): AppendEvent;
    getPartitionId(): number;
    setPartitionId(value: number): AppendEvent;
    getName(): string;
    setName(value: string): AppendEvent;
    getSchema(): string;
    setSchema(value: string): AppendEvent;

    hasMetadata(): boolean;
    clearMetadata(): void;
    getMetadata(): EventMetadata | undefined;
    setMetadata(value?: EventMetadata): AppendEvent;
    getData(): Uint8Array | string;
    getData_asU8(): Uint8Array;
    getData_asB64(): string;
    setData(value: Uint8Array | string): AppendEvent;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): AppendEvent.AsObject;
    static toObject(includeInstance: boolean, msg: AppendEvent): AppendEvent.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: AppendEvent, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): AppendEvent;
    static deserializeBinaryFromReader(message: AppendEvent, reader: jspb.BinaryReader): AppendEvent;
}

export namespace AppendEvent {
    export type AsObject = {
        id: string,
        aggregateName: string,
        aggregateId: string,
        version: number,
        partitionId: number,
        name: string,
        schema: string,
        metadata?: EventMetadata.AsObject,
        data: Uint8Array | string,
    }
}

export class AppendRequest extends jspb.Message { 
    getRequestId(): string;
    setRequestId(value: string): AppendRequest;
    clearEventsList(): void;
    getEventsList(): Array<AppendEvent>;
    setEventsList(value: Array<AppendEvent>): AppendRequest;
    addEvents(value?: AppendEvent, index?: number): AppendEvent;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): AppendRequest.AsObject;
    static toObject(includeInstance: boolean, msg: AppendRequest): AppendRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: AppendRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): AppendRequest;
    static deserializeBinaryFromReader(message: AppendRequest, reader: jspb.BinaryReader): AppendRequest;
}

export namespace AppendRequest {
    export type AsObject = {
        requestId: string,
        eventsList: Array<AppendEvent.AsObject>,
    }
}

export class AppendResponse extends jspb.Message { 

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): AppendResponse.AsObject;
    static toObject(includeInstance: boolean, msg: AppendResponse): AppendResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: AppendResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): AppendResponse;
    static deserializeBinaryFromReader(message: AppendResponse, reader: jspb.BinaryReader): AppendResponse;
}

export namespace AppendResponse {
    export type AsObject = {
    }
}

export class GetRequest extends jspb.Message { 
    getAggregateId(): string;
    setAggregateId(value: string): GetRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GetRequest.AsObject;
    static toObject(includeInstance: boolean, msg: GetRequest): GetRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GetRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GetRequest;
    static deserializeBinaryFromReader(message: GetRequest, reader: jspb.BinaryReader): GetRequest;
}

export namespace GetRequest {
    export type AsObject = {
        aggregateId: string,
    }
}

export class GetByAggregateIDAndNameRequest extends jspb.Message { 
    getAggregateId(): string;
    setAggregateId(value: string): GetByAggregateIDAndNameRequest;
    getName(): string;
    setName(value: string): GetByAggregateIDAndNameRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GetByAggregateIDAndNameRequest.AsObject;
    static toObject(includeInstance: boolean, msg: GetByAggregateIDAndNameRequest): GetByAggregateIDAndNameRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GetByAggregateIDAndNameRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GetByAggregateIDAndNameRequest;
    static deserializeBinaryFromReader(message: GetByAggregateIDAndNameRequest, reader: jspb.BinaryReader): GetByAggregateIDAndNameRequest;
}

export namespace GetByAggregateIDAndNameRequest {
    export type AsObject = {
        aggregateId: string,
        name: string,
    }
}

export class GetAfterGlobalVersionRequest extends jspb.Message { 
    getGlobalVersion(): number;
    setGlobalVersion(value: number): GetAfterGlobalVersionRequest;
    getPartitionId(): number;
    setPartitionId(value: number): GetAfterGlobalVersionRequest;
    getLimit(): number;
    setLimit(value: number): GetAfterGlobalVersionRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GetAfterGlobalVersionRequest.AsObject;
    static toObject(includeInstance: boolean, msg: GetAfterGlobalVersionRequest): GetAfterGlobalVersionRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GetAfterGlobalVersionRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GetAfterGlobalVersionRequest;
    static deserializeBinaryFromReader(message: GetAfterGlobalVersionRequest, reader: jspb.BinaryReader): GetAfterGlobalVersionRequest;
}

export namespace GetAfterGlobalVersionRequest {
    export type AsObject = {
        globalVersion: number,
        partitionId: number,
        limit: number,
    }
}

export class GetOrCreateOutboxRequest extends jspb.Message { 
    getOutboxId(): string;
    setOutboxId(value: string): GetOrCreateOutboxRequest;
    getPartitionId(): string;
    setPartitionId(value: string): GetOrCreateOutboxRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GetOrCreateOutboxRequest.AsObject;
    static toObject(includeInstance: boolean, msg: GetOrCreateOutboxRequest): GetOrCreateOutboxRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GetOrCreateOutboxRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GetOrCreateOutboxRequest;
    static deserializeBinaryFromReader(message: GetOrCreateOutboxRequest, reader: jspb.BinaryReader): GetOrCreateOutboxRequest;
}

export namespace GetOrCreateOutboxRequest {
    export type AsObject = {
        outboxId: string,
        partitionId: string,
    }
}

export class GetOrCreateOutboxResponse extends jspb.Message { 
    getGlobalVersion(): number;
    setGlobalVersion(value: number): GetOrCreateOutboxResponse;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GetOrCreateOutboxResponse.AsObject;
    static toObject(includeInstance: boolean, msg: GetOrCreateOutboxResponse): GetOrCreateOutboxResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GetOrCreateOutboxResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GetOrCreateOutboxResponse;
    static deserializeBinaryFromReader(message: GetOrCreateOutboxResponse, reader: jspb.BinaryReader): GetOrCreateOutboxResponse;
}

export namespace GetOrCreateOutboxResponse {
    export type AsObject = {
        globalVersion: number,
    }
}

export class UpdateOutboxPositionRequest extends jspb.Message { 
    getOutboxId(): string;
    setOutboxId(value: string): UpdateOutboxPositionRequest;
    getGlobalVersion(): number;
    setGlobalVersion(value: number): UpdateOutboxPositionRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): UpdateOutboxPositionRequest.AsObject;
    static toObject(includeInstance: boolean, msg: UpdateOutboxPositionRequest): UpdateOutboxPositionRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: UpdateOutboxPositionRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): UpdateOutboxPositionRequest;
    static deserializeBinaryFromReader(message: UpdateOutboxPositionRequest, reader: jspb.BinaryReader): UpdateOutboxPositionRequest;
}

export namespace UpdateOutboxPositionRequest {
    export type AsObject = {
        outboxId: string,
        globalVersion: number,
    }
}

export class UpdateOutboxPositionResponse extends jspb.Message { 

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): UpdateOutboxPositionResponse.AsObject;
    static toObject(includeInstance: boolean, msg: UpdateOutboxPositionResponse): UpdateOutboxPositionResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: UpdateOutboxPositionResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): UpdateOutboxPositionResponse;
    static deserializeBinaryFromReader(message: UpdateOutboxPositionResponse, reader: jspb.BinaryReader): UpdateOutboxPositionResponse;
}

export namespace UpdateOutboxPositionResponse {
    export type AsObject = {
    }
}
