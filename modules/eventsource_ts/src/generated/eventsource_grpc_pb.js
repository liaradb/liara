// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('@grpc/grpc-js');
var eventsource_pb = require('./eventsource_pb.js');
var google_protobuf_timestamp_pb = require('google-protobuf/google/protobuf/timestamp_pb.js');

function serialize_todo_AppendRequest(arg) {
  if (!(arg instanceof eventsource_pb.AppendRequest)) {
    throw new Error('Expected argument of type todo.AppendRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_todo_AppendRequest(buffer_arg) {
  return eventsource_pb.AppendRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_todo_AppendResponse(arg) {
  if (!(arg instanceof eventsource_pb.AppendResponse)) {
    throw new Error('Expected argument of type todo.AppendResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_todo_AppendResponse(buffer_arg) {
  return eventsource_pb.AppendResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_todo_Event(arg) {
  if (!(arg instanceof eventsource_pb.Event)) {
    throw new Error('Expected argument of type todo.Event');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_todo_Event(buffer_arg) {
  return eventsource_pb.Event.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_todo_GetAfterGlobalVersionRequest(arg) {
  if (!(arg instanceof eventsource_pb.GetAfterGlobalVersionRequest)) {
    throw new Error('Expected argument of type todo.GetAfterGlobalVersionRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_todo_GetAfterGlobalVersionRequest(buffer_arg) {
  return eventsource_pb.GetAfterGlobalVersionRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_todo_GetByAggregateIDAndNameRequest(arg) {
  if (!(arg instanceof eventsource_pb.GetByAggregateIDAndNameRequest)) {
    throw new Error('Expected argument of type todo.GetByAggregateIDAndNameRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_todo_GetByAggregateIDAndNameRequest(buffer_arg) {
  return eventsource_pb.GetByAggregateIDAndNameRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_todo_GetOrCreateOutboxRequest(arg) {
  if (!(arg instanceof eventsource_pb.GetOrCreateOutboxRequest)) {
    throw new Error('Expected argument of type todo.GetOrCreateOutboxRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_todo_GetOrCreateOutboxRequest(buffer_arg) {
  return eventsource_pb.GetOrCreateOutboxRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_todo_GetOrCreateOutboxResponse(arg) {
  if (!(arg instanceof eventsource_pb.GetOrCreateOutboxResponse)) {
    throw new Error('Expected argument of type todo.GetOrCreateOutboxResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_todo_GetOrCreateOutboxResponse(buffer_arg) {
  return eventsource_pb.GetOrCreateOutboxResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_todo_GetRequest(arg) {
  if (!(arg instanceof eventsource_pb.GetRequest)) {
    throw new Error('Expected argument of type todo.GetRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_todo_GetRequest(buffer_arg) {
  return eventsource_pb.GetRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_todo_UpdateOutboxPositionRequest(arg) {
  if (!(arg instanceof eventsource_pb.UpdateOutboxPositionRequest)) {
    throw new Error('Expected argument of type todo.UpdateOutboxPositionRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_todo_UpdateOutboxPositionRequest(buffer_arg) {
  return eventsource_pb.UpdateOutboxPositionRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_todo_UpdateOutboxPositionResponse(arg) {
  if (!(arg instanceof eventsource_pb.UpdateOutboxPositionResponse)) {
    throw new Error('Expected argument of type todo.UpdateOutboxPositionResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_todo_UpdateOutboxPositionResponse(buffer_arg) {
  return eventsource_pb.UpdateOutboxPositionResponse.deserializeBinary(new Uint8Array(buffer_arg));
}


var EventSourceServiceService = exports.EventSourceServiceService = {
  append: {
    path: '/todo.EventSourceService/Append',
    requestStream: false,
    responseStream: false,
    requestType: eventsource_pb.AppendRequest,
    responseType: eventsource_pb.AppendResponse,
    requestSerialize: serialize_todo_AppendRequest,
    requestDeserialize: deserialize_todo_AppendRequest,
    responseSerialize: serialize_todo_AppendResponse,
    responseDeserialize: deserialize_todo_AppendResponse,
  },
  get: {
    path: '/todo.EventSourceService/Get',
    requestStream: false,
    responseStream: true,
    requestType: eventsource_pb.GetRequest,
    responseType: eventsource_pb.Event,
    requestSerialize: serialize_todo_GetRequest,
    requestDeserialize: deserialize_todo_GetRequest,
    responseSerialize: serialize_todo_Event,
    responseDeserialize: deserialize_todo_Event,
  },
  getByAggregateIDAndName: {
    path: '/todo.EventSourceService/GetByAggregateIDAndName',
    requestStream: false,
    responseStream: true,
    requestType: eventsource_pb.GetByAggregateIDAndNameRequest,
    responseType: eventsource_pb.Event,
    requestSerialize: serialize_todo_GetByAggregateIDAndNameRequest,
    requestDeserialize: deserialize_todo_GetByAggregateIDAndNameRequest,
    responseSerialize: serialize_todo_Event,
    responseDeserialize: deserialize_todo_Event,
  },
  getAfterGlobalVersion: {
    path: '/todo.EventSourceService/GetAfterGlobalVersion',
    requestStream: false,
    responseStream: true,
    requestType: eventsource_pb.GetAfterGlobalVersionRequest,
    responseType: eventsource_pb.Event,
    requestSerialize: serialize_todo_GetAfterGlobalVersionRequest,
    requestDeserialize: deserialize_todo_GetAfterGlobalVersionRequest,
    responseSerialize: serialize_todo_Event,
    responseDeserialize: deserialize_todo_Event,
  },
  getOrCreateOutbox: {
    path: '/todo.EventSourceService/GetOrCreateOutbox',
    requestStream: false,
    responseStream: false,
    requestType: eventsource_pb.GetOrCreateOutboxRequest,
    responseType: eventsource_pb.GetOrCreateOutboxResponse,
    requestSerialize: serialize_todo_GetOrCreateOutboxRequest,
    requestDeserialize: deserialize_todo_GetOrCreateOutboxRequest,
    responseSerialize: serialize_todo_GetOrCreateOutboxResponse,
    responseDeserialize: deserialize_todo_GetOrCreateOutboxResponse,
  },
  updateOutboxPosition: {
    path: '/todo.EventSourceService/UpdateOutboxPosition',
    requestStream: false,
    responseStream: false,
    requestType: eventsource_pb.UpdateOutboxPositionRequest,
    responseType: eventsource_pb.UpdateOutboxPositionResponse,
    requestSerialize: serialize_todo_UpdateOutboxPositionRequest,
    requestDeserialize: deserialize_todo_UpdateOutboxPositionRequest,
    responseSerialize: serialize_todo_UpdateOutboxPositionResponse,
    responseDeserialize: deserialize_todo_UpdateOutboxPositionResponse,
  },
};

exports.EventSourceServiceClient = grpc.makeGenericClientConstructor(EventSourceServiceService);
