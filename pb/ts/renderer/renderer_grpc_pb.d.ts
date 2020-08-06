// package: renderer
// file: renderer.proto

/* tslint:disable */
/* eslint-disable */

import * as grpc from "grpc";
import * as renderer_pb from "./renderer_pb";

interface IRendererService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
    render: IRendererService_IRender;
}

interface IRendererService_IRender extends grpc.MethodDefinition<renderer_pb.RenderRequest, renderer_pb.RenderReply> {
    path: string; // "/renderer.Renderer/Render"
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<renderer_pb.RenderRequest>;
    requestDeserialize: grpc.deserialize<renderer_pb.RenderRequest>;
    responseSerialize: grpc.serialize<renderer_pb.RenderReply>;
    responseDeserialize: grpc.deserialize<renderer_pb.RenderReply>;
}

export const RendererService: IRendererService;

export interface IRendererServer {
    render: grpc.handleUnaryCall<renderer_pb.RenderRequest, renderer_pb.RenderReply>;
}

export interface IRendererClient {
    render(request: renderer_pb.RenderRequest, callback: (error: grpc.ServiceError | null, response: renderer_pb.RenderReply) => void): grpc.ClientUnaryCall;
    render(request: renderer_pb.RenderRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: renderer_pb.RenderReply) => void): grpc.ClientUnaryCall;
    render(request: renderer_pb.RenderRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: renderer_pb.RenderReply) => void): grpc.ClientUnaryCall;
}

export class RendererClient extends grpc.Client implements IRendererClient {
    constructor(address: string, credentials: grpc.ChannelCredentials, options?: object);
    public render(request: renderer_pb.RenderRequest, callback: (error: grpc.ServiceError | null, response: renderer_pb.RenderReply) => void): grpc.ClientUnaryCall;
    public render(request: renderer_pb.RenderRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: renderer_pb.RenderReply) => void): grpc.ClientUnaryCall;
    public render(request: renderer_pb.RenderRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: renderer_pb.RenderReply) => void): grpc.ClientUnaryCall;
}
