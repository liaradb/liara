// package: todo
// file: eventsource.proto

/* tslint:disable */
/* eslint-disable */

import * as jspb from "google-protobuf";
import * as google_protobuf_timestamp_pb from "google-protobuf/google/protobuf/timestamp_pb";

export class Event extends jspb.Message { 
    getGlobalVersion(): number;
    setGlobalVersion(value: number): Event;
    getAggregateName(): string;
    setAggregateName(value: string): Event;
    getId(): string;
    setId(value: string): Event;
    getAggregateId(): string;
    setAggregateId(value: string): Event;
    getVersion(): number;
    setVersion(value: number): Event;
    getName(): string;
    setName(value: string): Event;
    getCorrelationId(): string;
    setCorrelationId(value: string): Event;
    getUserId(): string;
    setUserId(value: string): Event;

    hasTime(): boolean;
    clearTime(): void;
    getTime(): google_protobuf_timestamp_pb.Timestamp | undefined;
    setTime(value?: google_protobuf_timestamp_pb.Timestamp): Event;
    getSchema(): string;
    setSchema(value: string): Event;
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
        aggregateName: string,
        id: string,
        aggregateId: string,
        version: number,
        name: string,
        correlationId: string,
        userId: string,
        time?: google_protobuf_timestamp_pb.Timestamp.AsObject,
        schema: string,
        data: Uint8Array | string,
    }
}

export class AppendRequest extends jspb.Message { 
    clearEventsList(): void;
    getEventsList(): Array<Event>;
    setEventsList(value: Array<Event>): AppendRequest;
    addEvents(value?: Event, index?: number): Event;

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
        eventsList: Array<Event.AsObject>,
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
        limit: number,
    }
}
