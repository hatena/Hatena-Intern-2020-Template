// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('grpc');
var account_pb = require('./account_pb.js');

function serialize_account_SigninReply(arg) {
  if (!(arg instanceof account_pb.SigninReply)) {
    throw new Error('Expected argument of type account.SigninReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_account_SigninReply(buffer_arg) {
  return account_pb.SigninReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_account_SigninRequest(arg) {
  if (!(arg instanceof account_pb.SigninRequest)) {
    throw new Error('Expected argument of type account.SigninRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_account_SigninRequest(buffer_arg) {
  return account_pb.SigninRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_account_SignupReply(arg) {
  if (!(arg instanceof account_pb.SignupReply)) {
    throw new Error('Expected argument of type account.SignupReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_account_SignupReply(buffer_arg) {
  return account_pb.SignupReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_account_SignupRequest(arg) {
  if (!(arg instanceof account_pb.SignupRequest)) {
    throw new Error('Expected argument of type account.SignupRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_account_SignupRequest(buffer_arg) {
  return account_pb.SignupRequest.deserializeBinary(new Uint8Array(buffer_arg));
}


var AccountService = exports.AccountService = {
  signup: {
    path: '/account.Account/Signup',
    requestStream: false,
    responseStream: false,
    requestType: account_pb.SignupRequest,
    responseType: account_pb.SignupReply,
    requestSerialize: serialize_account_SignupRequest,
    requestDeserialize: deserialize_account_SignupRequest,
    responseSerialize: serialize_account_SignupReply,
    responseDeserialize: deserialize_account_SignupReply,
  },
  signin: {
    path: '/account.Account/Signin',
    requestStream: false,
    responseStream: false,
    requestType: account_pb.SigninRequest,
    responseType: account_pb.SigninReply,
    requestSerialize: serialize_account_SigninRequest,
    requestDeserialize: deserialize_account_SigninRequest,
    responseSerialize: serialize_account_SigninReply,
    responseDeserialize: deserialize_account_SigninReply,
  },
};

exports.AccountClient = grpc.makeGenericClientConstructor(AccountService);
