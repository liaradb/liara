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

function serialize_todo_CreateOutboxRequest(arg) {
  if (!(arg instanceof eventsource_pb.CreateOutboxRequest)) {
    throw new Error('Expected argument of type todo.CreateOutboxRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_todo_CreateOutboxRequest(buffer_arg) {
  return eventsource_pb.CreateOutboxRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_todo_CreateOutboxResponse(arg) {
  if (!(arg instanceof eventsource_pb.CreateOutboxResponse)) {
    throw new Error('Expected argument of type todo.CreateOutboxResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_todo_CreateOutboxResponse(buffer_arg) {
  return eventsource_pb.CreateOutboxResponse.deserializeBinary(new Uint8Array(buffer_arg));
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

function serialize_todo_GetByOutboxRequest(arg) {
  if (!(arg instanceof eventsource_pb.GetByOutboxRequest)) {
    throw new Error('Expected argument of type todo.GetByOutboxRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_todo_GetByOutboxRequest(buffer_arg) {
  return eventsource_pb.GetByOutboxRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_todo_GetOutboxRequest(arg) {
  if (!(arg instanceof eventsource_pb.GetOutboxRequest)) {
    throw new Error('Expected argument of type todo.GetOutboxRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_todo_GetOutboxRequest(buffer_arg) {
  return eventsource_pb.GetOutboxRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_todo_GetOutboxResponse(arg) {
  if (!(arg instanceof eventsource_pb.GetOutboxResponse)) {
    throw new Error('Expected argument of type todo.GetOutboxResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_todo_GetOutboxResponse(buffer_arg) {
  return eventsource_pb.GetOutboxResponse.deserializeBinary(new Uint8Array(buffer_arg));
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

function serialize_todo_ListTenantsRequest(arg) {
  if (!(arg instanceof eventsource_pb.ListTenantsRequest)) {
    throw new Error('Expected argument of type todo.ListTenantsRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_todo_ListTenantsRequest(buffer_arg) {
  return eventsource_pb.ListTenantsRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_todo_Tenant(arg) {
  if (!(arg instanceof eventsource_pb.Tenant)) {
    throw new Error('Expected argument of type todo.Tenant');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_todo_Tenant(buffer_arg) {
  return eventsource_pb.Tenant.deserializeBinary(new Uint8Array(buffer_arg));
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
  getByOutbox: {
    path: '/todo.EventSourceService/GetByOutbox',
    requestStream: false,
    responseStream: true,
    requestType: eventsource_pb.GetByOutboxRequest,
    responseType: eventsource_pb.Event,
    requestSerialize: serialize_todo_GetByOutboxRequest,
    requestDeserialize: deserialize_todo_GetByOutboxRequest,
    responseSerialize: serialize_todo_Event,
    responseDeserialize: deserialize_todo_Event,
  },
  createOutbox: {
    path: '/todo.EventSourceService/CreateOutbox',
    requestStream: false,
    responseStream: false,
    requestType: eventsource_pb.CreateOutboxRequest,
    responseType: eventsource_pb.CreateOutboxResponse,
    requestSerialize: serialize_todo_CreateOutboxRequest,
    requestDeserialize: deserialize_todo_CreateOutboxRequest,
    responseSerialize: serialize_todo_CreateOutboxResponse,
    responseDeserialize: deserialize_todo_CreateOutboxResponse,
  },
  getOutbox: {
    path: '/todo.EventSourceService/GetOutbox',
    requestStream: false,
    responseStream: false,
    requestType: eventsource_pb.GetOutboxRequest,
    responseType: eventsource_pb.GetOutboxResponse,
    requestSerialize: serialize_todo_GetOutboxRequest,
    requestDeserialize: deserialize_todo_GetOutboxRequest,
    responseSerialize: serialize_todo_GetOutboxResponse,
    responseDeserialize: deserialize_todo_GetOutboxResponse,
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
  listTenants: {
    path: '/todo.EventSourceService/ListTenants',
    requestStream: false,
    responseStream: true,
    requestType: eventsource_pb.ListTenantsRequest,
    responseType: eventsource_pb.Tenant,
    requestSerialize: serialize_todo_ListTenantsRequest,
    requestDeserialize: deserialize_todo_ListTenantsRequest,
    responseSerialize: serialize_todo_Tenant,
    responseDeserialize: deserialize_todo_Tenant,
  },
};

exports.EventSourceServiceClient = grpc.makeGenericClientConstructor(EventSourceServiceService);
