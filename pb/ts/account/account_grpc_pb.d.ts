// package: account
// file: account.proto

/* tslint:disable */
/* eslint-disable */

import * as grpc from "grpc";
import * as account_pb from "./account_pb";

interface IAccountService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
    signup: IAccountService_ISignup;
    signin: IAccountService_ISignin;
}

interface IAccountService_ISignup extends grpc.MethodDefinition<account_pb.SignupRequest, account_pb.SignupReply> {
    path: string; // "/account.Account/Signup"
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<account_pb.SignupRequest>;
    requestDeserialize: grpc.deserialize<account_pb.SignupRequest>;
    responseSerialize: grpc.serialize<account_pb.SignupReply>;
    responseDeserialize: grpc.deserialize<account_pb.SignupReply>;
}
interface IAccountService_ISignin extends grpc.MethodDefinition<account_pb.SigninRequest, account_pb.SigninReply> {
    path: string; // "/account.Account/Signin"
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<account_pb.SigninRequest>;
    requestDeserialize: grpc.deserialize<account_pb.SigninRequest>;
    responseSerialize: grpc.serialize<account_pb.SigninReply>;
    responseDeserialize: grpc.deserialize<account_pb.SigninReply>;
}

export const AccountService: IAccountService;

export interface IAccountServer {
    signup: grpc.handleUnaryCall<account_pb.SignupRequest, account_pb.SignupReply>;
    signin: grpc.handleUnaryCall<account_pb.SigninRequest, account_pb.SigninReply>;
}

export interface IAccountClient {
    signup(request: account_pb.SignupRequest, callback: (error: grpc.ServiceError | null, response: account_pb.SignupReply) => void): grpc.ClientUnaryCall;
    signup(request: account_pb.SignupRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: account_pb.SignupReply) => void): grpc.ClientUnaryCall;
    signup(request: account_pb.SignupRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: account_pb.SignupReply) => void): grpc.ClientUnaryCall;
    signin(request: account_pb.SigninRequest, callback: (error: grpc.ServiceError | null, response: account_pb.SigninReply) => void): grpc.ClientUnaryCall;
    signin(request: account_pb.SigninRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: account_pb.SigninReply) => void): grpc.ClientUnaryCall;
    signin(request: account_pb.SigninRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: account_pb.SigninReply) => void): grpc.ClientUnaryCall;
}

export class AccountClient extends grpc.Client implements IAccountClient {
    constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
    public signup(request: account_pb.SignupRequest, callback: (error: grpc.ServiceError | null, response: account_pb.SignupReply) => void): grpc.ClientUnaryCall;
    public signup(request: account_pb.SignupRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: account_pb.SignupReply) => void): grpc.ClientUnaryCall;
    public signup(request: account_pb.SignupRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: account_pb.SignupReply) => void): grpc.ClientUnaryCall;
    public signin(request: account_pb.SigninRequest, callback: (error: grpc.ServiceError | null, response: account_pb.SigninReply) => void): grpc.ClientUnaryCall;
    public signin(request: account_pb.SigninRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: account_pb.SigninReply) => void): grpc.ClientUnaryCall;
    public signin(request: account_pb.SigninRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: account_pb.SigninReply) => void): grpc.ClientUnaryCall;
}
