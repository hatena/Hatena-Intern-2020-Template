// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var renderer_pb = require('./renderer_pb.js');

function serialize_renderer_RenderReply(arg) {
  if (!(arg instanceof renderer_pb.RenderReply)) {
    throw new Error('Expected argument of type renderer.RenderReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_renderer_RenderReply(buffer_arg) {
  return renderer_pb.RenderReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_renderer_RenderRequest(arg) {
  if (!(arg instanceof renderer_pb.RenderRequest)) {
    throw new Error('Expected argument of type renderer.RenderRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_renderer_RenderRequest(buffer_arg) {
  return renderer_pb.RenderRequest.deserializeBinary(new Uint8Array(buffer_arg));
}


var RendererService = exports.RendererService = {
  render: {
    path: '/renderer.Renderer/Render',
    requestStream: false,
    responseStream: false,
    requestType: renderer_pb.RenderRequest,
    responseType: renderer_pb.RenderReply,
    requestSerialize: serialize_renderer_RenderRequest,
    requestDeserialize: deserialize_renderer_RenderRequest,
    responseSerialize: serialize_renderer_RenderReply,
    responseDeserialize: deserialize_renderer_RenderReply,
  },
};

exports.RendererClient = grpc.makeGenericClientConstructor(RendererService);
