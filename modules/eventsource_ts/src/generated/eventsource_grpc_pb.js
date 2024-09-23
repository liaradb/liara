// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('@grpc/grpc-js');
var eventsource_pb = require('./eventsource_pb.js');

function serialize_todo_ListEventsRequest(arg) {
  if (!(arg instanceof eventsource_pb.ListEventsRequest)) {
    throw new Error('Expected argument of type todo.ListEventsRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_todo_ListEventsRequest(buffer_arg) {
  return eventsource_pb.ListEventsRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_todo_ListEventsResponse(arg) {
  if (!(arg instanceof eventsource_pb.ListEventsResponse)) {
    throw new Error('Expected argument of type todo.ListEventsResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_todo_ListEventsResponse(buffer_arg) {
  return eventsource_pb.ListEventsResponse.deserializeBinary(new Uint8Array(buffer_arg));
}


var EventSourceServiceService = exports.EventSourceServiceService = {
  listEvents: {
    path: '/todo.EventSourceService/ListEvents',
    requestStream: false,
    responseStream: false,
    requestType: eventsource_pb.ListEventsRequest,
    responseType: eventsource_pb.ListEventsResponse,
    requestSerialize: serialize_todo_ListEventsRequest,
    requestDeserialize: deserialize_todo_ListEventsRequest,
    responseSerialize: serialize_todo_ListEventsResponse,
    responseDeserialize: deserialize_todo_ListEventsResponse,
  },
};

exports.EventSourceServiceClient = grpc.makeGenericClientConstructor(EventSourceServiceService);
