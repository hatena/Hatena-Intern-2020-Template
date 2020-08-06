// package: account
// file: account.proto

/* tslint:disable */
/* eslint-disable */

import * as jspb from "google-protobuf";

export class SignupRequest extends jspb.Message { 
    getName(): string;
    setName(value: string): SignupRequest;

    getPassword(): string;
    setPassword(value: string): SignupRequest;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SignupRequest.AsObject;
    static toObject(includeInstance: boolean, msg: SignupRequest): SignupRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SignupRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SignupRequest;
    static deserializeBinaryFromReader(message: SignupRequest, reader: jspb.BinaryReader): SignupRequest;
}

export namespace SignupRequest {
    export type AsObject = {
        name: string,
        password: string,
    }
}

export class SignupReply extends jspb.Message { 
    getToken(): string;
    setToken(value: string): SignupReply;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SignupReply.AsObject;
    static toObject(includeInstance: boolean, msg: SignupReply): SignupReply.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SignupReply, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SignupReply;
    static deserializeBinaryFromReader(message: SignupReply, reader: jspb.BinaryReader): SignupReply;
}

export namespace SignupReply {
    export type AsObject = {
        token: string,
    }
}

export class SigninRequest extends jspb.Message { 
    getName(): string;
    setName(value: string): SigninRequest;

    getPassword(): string;
    setPassword(value: string): SigninRequest;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SigninRequest.AsObject;
    static toObject(includeInstance: boolean, msg: SigninRequest): SigninRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SigninRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SigninRequest;
    static deserializeBinaryFromReader(message: SigninRequest, reader: jspb.BinaryReader): SigninRequest;
}

export namespace SigninRequest {
    export type AsObject = {
        name: string,
        password: string,
    }
}

export class SigninReply extends jspb.Message { 
    getToken(): string;
    setToken(value: string): SigninReply;


    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SigninReply.AsObject;
    static toObject(includeInstance: boolean, msg: SigninReply): SigninReply.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SigninReply, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SigninReply;
    static deserializeBinaryFromReader(message: SigninReply, reader: jspb.BinaryReader): SigninReply;
}

export namespace SigninReply {
    export type AsObject = {
        token: string,
    }
}
