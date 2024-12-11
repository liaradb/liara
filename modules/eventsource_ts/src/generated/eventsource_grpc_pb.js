// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('@grpc/grpc-js');
var eventsource_pb = require('./eventsource_pb.js');
var google_protobuf_timestamp_pb = require('google-protobuf/google/protobuf/timestamp_pb.js');

function serialize_liara_AppendRequest(arg) {
  if (!(arg instanceof eventsource_pb.AppendRequest)) {
    throw new Error('Expected argument of type liara.AppendRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_liara_AppendRequest(buffer_arg) {
  return eventsource_pb.AppendRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_liara_AppendResponse(arg) {
  if (!(arg instanceof eventsource_pb.AppendResponse)) {
    throw new Error('Expected argument of type liara.AppendResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_liara_AppendResponse(buffer_arg) {
  return eventsource_pb.AppendResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_liara_CreateOutboxRequest(arg) {
  if (!(arg instanceof eventsource_pb.CreateOutboxRequest)) {
    throw new Error('Expected argument of type liara.CreateOutboxRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_liara_CreateOutboxRequest(buffer_arg) {
  return eventsource_pb.CreateOutboxRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_liara_CreateOutboxResponse(arg) {
  if (!(arg instanceof eventsource_pb.CreateOutboxResponse)) {
    throw new Error('Expected argument of type liara.CreateOutboxResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_liara_CreateOutboxResponse(buffer_arg) {
  return eventsource_pb.CreateOutboxResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_liara_CreateTenantReponse(arg) {
  if (!(arg instanceof eventsource_pb.CreateTenantReponse)) {
    throw new Error('Expected argument of type liara.CreateTenantReponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_liara_CreateTenantReponse(buffer_arg) {
  return eventsource_pb.CreateTenantReponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_liara_CreateTenantRequest(arg) {
  if (!(arg instanceof eventsource_pb.CreateTenantRequest)) {
    throw new Error('Expected argument of type liara.CreateTenantRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_liara_CreateTenantRequest(buffer_arg) {
  return eventsource_pb.CreateTenantRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_liara_DeleteTenantRequest(arg) {
  if (!(arg instanceof eventsource_pb.DeleteTenantRequest)) {
    throw new Error('Expected argument of type liara.DeleteTenantRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_liara_DeleteTenantRequest(buffer_arg) {
  return eventsource_pb.DeleteTenantRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_liara_DeleteTenantResponse(arg) {
  if (!(arg instanceof eventsource_pb.DeleteTenantResponse)) {
    throw new Error('Expected argument of type liara.DeleteTenantResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_liara_DeleteTenantResponse(buffer_arg) {
  return eventsource_pb.DeleteTenantResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_liara_Event(arg) {
  if (!(arg instanceof eventsource_pb.Event)) {
    throw new Error('Expected argument of type liara.Event');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_liara_Event(buffer_arg) {
  return eventsource_pb.Event.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_liara_GetAfterGlobalVersionRequest(arg) {
  if (!(arg instanceof eventsource_pb.GetAfterGlobalVersionRequest)) {
    throw new Error('Expected argument of type liara.GetAfterGlobalVersionRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_liara_GetAfterGlobalVersionRequest(buffer_arg) {
  return eventsource_pb.GetAfterGlobalVersionRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_liara_GetByAggregateIDAndNameRequest(arg) {
  if (!(arg instanceof eventsource_pb.GetByAggregateIDAndNameRequest)) {
    throw new Error('Expected argument of type liara.GetByAggregateIDAndNameRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_liara_GetByAggregateIDAndNameRequest(buffer_arg) {
  return eventsource_pb.GetByAggregateIDAndNameRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_liara_GetByOutboxRequest(arg) {
  if (!(arg instanceof eventsource_pb.GetByOutboxRequest)) {
    throw new Error('Expected argument of type liara.GetByOutboxRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_liara_GetByOutboxRequest(buffer_arg) {
  return eventsource_pb.GetByOutboxRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_liara_GetOutboxRequest(arg) {
  if (!(arg instanceof eventsource_pb.GetOutboxRequest)) {
    throw new Error('Expected argument of type liara.GetOutboxRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_liara_GetOutboxRequest(buffer_arg) {
  return eventsource_pb.GetOutboxRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_liara_GetOutboxResponse(arg) {
  if (!(arg instanceof eventsource_pb.GetOutboxResponse)) {
    throw new Error('Expected argument of type liara.GetOutboxResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_liara_GetOutboxResponse(buffer_arg) {
  return eventsource_pb.GetOutboxResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_liara_GetRequest(arg) {
  if (!(arg instanceof eventsource_pb.GetRequest)) {
    throw new Error('Expected argument of type liara.GetRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_liara_GetRequest(buffer_arg) {
  return eventsource_pb.GetRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_liara_GetTenantRequest(arg) {
  if (!(arg instanceof eventsource_pb.GetTenantRequest)) {
    throw new Error('Expected argument of type liara.GetTenantRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_liara_GetTenantRequest(buffer_arg) {
  return eventsource_pb.GetTenantRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_liara_GetTenantResponse(arg) {
  if (!(arg instanceof eventsource_pb.GetTenantResponse)) {
    throw new Error('Expected argument of type liara.GetTenantResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_liara_GetTenantResponse(buffer_arg) {
  return eventsource_pb.GetTenantResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_liara_ListTenantsRequest(arg) {
  if (!(arg instanceof eventsource_pb.ListTenantsRequest)) {
    throw new Error('Expected argument of type liara.ListTenantsRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_liara_ListTenantsRequest(buffer_arg) {
  return eventsource_pb.ListTenantsRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_liara_RenameTenantRequest(arg) {
  if (!(arg instanceof eventsource_pb.RenameTenantRequest)) {
    throw new Error('Expected argument of type liara.RenameTenantRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_liara_RenameTenantRequest(buffer_arg) {
  return eventsource_pb.RenameTenantRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_liara_RenameTenantResponse(arg) {
  if (!(arg instanceof eventsource_pb.RenameTenantResponse)) {
    throw new Error('Expected argument of type liara.RenameTenantResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_liara_RenameTenantResponse(buffer_arg) {
  return eventsource_pb.RenameTenantResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_liara_Tenant(arg) {
  if (!(arg instanceof eventsource_pb.Tenant)) {
    throw new Error('Expected argument of type liara.Tenant');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_liara_Tenant(buffer_arg) {
  return eventsource_pb.Tenant.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_liara_TestIdempotencyRequest(arg) {
  if (!(arg instanceof eventsource_pb.TestIdempotencyRequest)) {
    throw new Error('Expected argument of type liara.TestIdempotencyRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_liara_TestIdempotencyRequest(buffer_arg) {
  return eventsource_pb.TestIdempotencyRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_liara_TestIdempotencyResponse(arg) {
  if (!(arg instanceof eventsource_pb.TestIdempotencyResponse)) {
    throw new Error('Expected argument of type liara.TestIdempotencyResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_liara_TestIdempotencyResponse(buffer_arg) {
  return eventsource_pb.TestIdempotencyResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_liara_UpdateOutboxPositionRequest(arg) {
  if (!(arg instanceof eventsource_pb.UpdateOutboxPositionRequest)) {
    throw new Error('Expected argument of type liara.UpdateOutboxPositionRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_liara_UpdateOutboxPositionRequest(buffer_arg) {
  return eventsource_pb.UpdateOutboxPositionRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_liara_UpdateOutboxPositionResponse(arg) {
  if (!(arg instanceof eventsource_pb.UpdateOutboxPositionResponse)) {
    throw new Error('Expected argument of type liara.UpdateOutboxPositionResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_liara_UpdateOutboxPositionResponse(buffer_arg) {
  return eventsource_pb.UpdateOutboxPositionResponse.deserializeBinary(new Uint8Array(buffer_arg));
}


var EventSourceServiceService = exports.EventSourceServiceService = {
  append: {
    path: '/liara.EventSourceService/Append',
    requestStream: false,
    responseStream: false,
    requestType: eventsource_pb.AppendRequest,
    responseType: eventsource_pb.AppendResponse,
    requestSerialize: serialize_liara_AppendRequest,
    requestDeserialize: deserialize_liara_AppendRequest,
    responseSerialize: serialize_liara_AppendResponse,
    responseDeserialize: deserialize_liara_AppendResponse,
  },
  testIdempotency: {
    path: '/liara.EventSourceService/TestIdempotency',
    requestStream: false,
    responseStream: false,
    requestType: eventsource_pb.TestIdempotencyRequest,
    responseType: eventsource_pb.TestIdempotencyResponse,
    requestSerialize: serialize_liara_TestIdempotencyRequest,
    requestDeserialize: deserialize_liara_TestIdempotencyRequest,
    responseSerialize: serialize_liara_TestIdempotencyResponse,
    responseDeserialize: deserialize_liara_TestIdempotencyResponse,
  },
  get: {
    path: '/liara.EventSourceService/Get',
    requestStream: false,
    responseStream: true,
    requestType: eventsource_pb.GetRequest,
    responseType: eventsource_pb.Event,
    requestSerialize: serialize_liara_GetRequest,
    requestDeserialize: deserialize_liara_GetRequest,
    responseSerialize: serialize_liara_Event,
    responseDeserialize: deserialize_liara_Event,
  },
  getByAggregateIDAndName: {
    path: '/liara.EventSourceService/GetByAggregateIDAndName',
    requestStream: false,
    responseStream: true,
    requestType: eventsource_pb.GetByAggregateIDAndNameRequest,
    responseType: eventsource_pb.Event,
    requestSerialize: serialize_liara_GetByAggregateIDAndNameRequest,
    requestDeserialize: deserialize_liara_GetByAggregateIDAndNameRequest,
    responseSerialize: serialize_liara_Event,
    responseDeserialize: deserialize_liara_Event,
  },
  getAfterGlobalVersion: {
    path: '/liara.EventSourceService/GetAfterGlobalVersion',
    requestStream: false,
    responseStream: true,
    requestType: eventsource_pb.GetAfterGlobalVersionRequest,
    responseType: eventsource_pb.Event,
    requestSerialize: serialize_liara_GetAfterGlobalVersionRequest,
    requestDeserialize: deserialize_liara_GetAfterGlobalVersionRequest,
    responseSerialize: serialize_liara_Event,
    responseDeserialize: deserialize_liara_Event,
  },
  getByOutbox: {
    path: '/liara.EventSourceService/GetByOutbox',
    requestStream: false,
    responseStream: true,
    requestType: eventsource_pb.GetByOutboxRequest,
    responseType: eventsource_pb.Event,
    requestSerialize: serialize_liara_GetByOutboxRequest,
    requestDeserialize: deserialize_liara_GetByOutboxRequest,
    responseSerialize: serialize_liara_Event,
    responseDeserialize: deserialize_liara_Event,
  },
  createOutbox: {
    path: '/liara.EventSourceService/CreateOutbox',
    requestStream: false,
    responseStream: false,
    requestType: eventsource_pb.CreateOutboxRequest,
    responseType: eventsource_pb.CreateOutboxResponse,
    requestSerialize: serialize_liara_CreateOutboxRequest,
    requestDeserialize: deserialize_liara_CreateOutboxRequest,
    responseSerialize: serialize_liara_CreateOutboxResponse,
    responseDeserialize: deserialize_liara_CreateOutboxResponse,
  },
  getOutbox: {
    path: '/liara.EventSourceService/GetOutbox',
    requestStream: false,
    responseStream: false,
    requestType: eventsource_pb.GetOutboxRequest,
    responseType: eventsource_pb.GetOutboxResponse,
    requestSerialize: serialize_liara_GetOutboxRequest,
    requestDeserialize: deserialize_liara_GetOutboxRequest,
    responseSerialize: serialize_liara_GetOutboxResponse,
    responseDeserialize: deserialize_liara_GetOutboxResponse,
  },
  updateOutboxPosition: {
    path: '/liara.EventSourceService/UpdateOutboxPosition',
    requestStream: false,
    responseStream: false,
    requestType: eventsource_pb.UpdateOutboxPositionRequest,
    responseType: eventsource_pb.UpdateOutboxPositionResponse,
    requestSerialize: serialize_liara_UpdateOutboxPositionRequest,
    requestDeserialize: deserialize_liara_UpdateOutboxPositionRequest,
    responseSerialize: serialize_liara_UpdateOutboxPositionResponse,
    responseDeserialize: deserialize_liara_UpdateOutboxPositionResponse,
  },
  createTenant: {
    path: '/liara.EventSourceService/CreateTenant',
    requestStream: false,
    responseStream: false,
    requestType: eventsource_pb.CreateTenantRequest,
    responseType: eventsource_pb.CreateTenantReponse,
    requestSerialize: serialize_liara_CreateTenantRequest,
    requestDeserialize: deserialize_liara_CreateTenantRequest,
    responseSerialize: serialize_liara_CreateTenantReponse,
    responseDeserialize: deserialize_liara_CreateTenantReponse,
  },
  deleteTenant: {
    path: '/liara.EventSourceService/DeleteTenant',
    requestStream: false,
    responseStream: false,
    requestType: eventsource_pb.DeleteTenantRequest,
    responseType: eventsource_pb.DeleteTenantResponse,
    requestSerialize: serialize_liara_DeleteTenantRequest,
    requestDeserialize: deserialize_liara_DeleteTenantRequest,
    responseSerialize: serialize_liara_DeleteTenantResponse,
    responseDeserialize: deserialize_liara_DeleteTenantResponse,
  },
  renameTenant: {
    path: '/liara.EventSourceService/RenameTenant',
    requestStream: false,
    responseStream: false,
    requestType: eventsource_pb.RenameTenantRequest,
    responseType: eventsource_pb.RenameTenantResponse,
    requestSerialize: serialize_liara_RenameTenantRequest,
    requestDeserialize: deserialize_liara_RenameTenantRequest,
    responseSerialize: serialize_liara_RenameTenantResponse,
    responseDeserialize: deserialize_liara_RenameTenantResponse,
  },
  getTenant: {
    path: '/liara.EventSourceService/GetTenant',
    requestStream: false,
    responseStream: false,
    requestType: eventsource_pb.GetTenantRequest,
    responseType: eventsource_pb.GetTenantResponse,
    requestSerialize: serialize_liara_GetTenantRequest,
    requestDeserialize: deserialize_liara_GetTenantRequest,
    responseSerialize: serialize_liara_GetTenantResponse,
    responseDeserialize: deserialize_liara_GetTenantResponse,
  },
  listTenants: {
    path: '/liara.EventSourceService/ListTenants',
    requestStream: false,
    responseStream: true,
    requestType: eventsource_pb.ListTenantsRequest,
    responseType: eventsource_pb.Tenant,
    requestSerialize: serialize_liara_ListTenantsRequest,
    requestDeserialize: deserialize_liara_ListTenantsRequest,
    responseSerialize: serialize_liara_Tenant,
    responseDeserialize: deserialize_liara_Tenant,
  },
};

exports.EventSourceServiceClient = grpc.makeGenericClientConstructor(EventSourceServiceService);
