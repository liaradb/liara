// package: todo
// file: eventsource.proto

/* tslint:disable */
/* eslint-disable */

import * as jspb from "google-protobuf";

export class ListEventsRequest extends jspb.Message { 

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ListEventsRequest.AsObject;
    static toObject(includeInstance: boolean, msg: ListEventsRequest): ListEventsRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ListEventsRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ListEventsRequest;
    static deserializeBinaryFromReader(message: ListEventsRequest, reader: jspb.BinaryReader): ListEventsRequest;
}

export namespace ListEventsRequest {
    export type AsObject = {
    }
}

export class ListEventsResponse extends jspb.Message { 

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ListEventsResponse.AsObject;
    static toObject(includeInstance: boolean, msg: ListEventsResponse): ListEventsResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ListEventsResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ListEventsResponse;
    static deserializeBinaryFromReader(message: ListEventsResponse, reader: jspb.BinaryReader): ListEventsResponse;
}

export namespace ListEventsResponse {
    export type AsObject = {
    }
}
