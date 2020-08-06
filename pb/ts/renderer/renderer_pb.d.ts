// package: renderer
// file: renderer.proto

/* tslint:disable */
/* eslint-disable */

import * as jspb from "google-protobuf";

export class RenderRequest extends jspb.Message { 
    getSrc(): string;
    setSrc(value: string): RenderRequest;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): RenderRequest.AsObject;
    static toObject(includeInstance: boolean, msg: RenderRequest): RenderRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: RenderRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): RenderRequest;
    static deserializeBinaryFromReader(message: RenderRequest, reader: jspb.BinaryReader): RenderRequest;
}

export namespace RenderRequest {
    export type AsObject = {
        src: string,
    }
}

export class RenderReply extends jspb.Message { 
    getHtml(): string;
    setHtml(value: string): RenderReply;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): RenderReply.AsObject;
    static toObject(includeInstance: boolean, msg: RenderReply): RenderReply.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: RenderReply, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): RenderReply;
    static deserializeBinaryFromReader(message: RenderReply, reader: jspb.BinaryReader): RenderReply;
}

export namespace RenderReply {
    export type AsObject = {
        html: string,
    }
}
