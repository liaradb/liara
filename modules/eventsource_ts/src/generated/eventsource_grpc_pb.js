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

function serialize_todo_GetByAggregateIDAndNameRequest(arg) {
  if (!(arg instanceof eventsource_pb.GetByAggregateIDAndNameRequest)) {
    throw new Error('Expected argument of type todo.GetByAggregateIDAndNameRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_todo_GetByAggregateIDAndNameRequest(buffer_arg) {
  return eventsource_pb.GetByAggregateIDAndNameRequest.deserializeBinary(new Uint8Array(buffer_arg));
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
};

exports.EventSourceServiceClient = grpc.makeGenericClientConstructor(EventSourceServiceService);
