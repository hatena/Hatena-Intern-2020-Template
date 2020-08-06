import grpc from "grpc";
import winston from "winston";

export type UnaryContext = Map<string, unknown>;

export type UnaryHandler<Request, Reply> = (req: Request, ctx: UnaryContext) => Promise<Reply>;

export type UnaryMiddleware = <Request, Reply>(
  next: UnaryHandler<Request, Reply>
) => UnaryHandler<Request, Reply>;

const callContextKey = "call";

export function unaryHandler<Request, Reply>(
  handler: UnaryHandler<Request, Reply>,
  ...middlewares: readonly UnaryMiddleware[]
): grpc.handleUnaryCall<Request, Reply> {
  return (call, callback) => {
    const ctx = new Map<string, unknown>();
    ctx.set(callContextKey, call);
    const h = applyUnaryMiddlewares(handler, middlewares);
    h(call.request, ctx).then(
      (reply) => {
        callback(null, reply);
      },
      (err) => {
        callback(err, null);
      }
    );
  };
}

function applyUnaryMiddlewares<Request, Reply>(
  handler: UnaryHandler<Request, Reply>,
  middlewares: readonly UnaryMiddleware[]
) {
  let h: UnaryHandler<Request, Reply> = handler;
  for (let i = middlewares.length - 1; i >= 0; i--) {
    h = middlewares[i](h);
  }
  return h;
}

export function createLoggingMiddleware(logger: winston.Logger): (rpc: string) => UnaryMiddleware {
  return (rpc) => <Request, Reply>(next: UnaryHandler<Request, Reply>) => async (
    req: Request,
    ctx: UnaryContext
  ) => {
    const start = Date.now();
    let code: grpc.status | undefined = undefined;
    let error: unknown = undefined;
    return next(req, ctx)
      .then(
        (reply) => {
          code = grpc.status.OK;
          return reply;
        },
        (err) => {
          code = err && err.code ? (err as grpc.ServiceError).code : grpc.status.UNKNOWN;
          error = err;
          throw err;
        }
      )
      .finally(() => {
        const end = Date.now();
        const level = codeToLevel(code ?? grpc.status.UNKNOWN);
        const call = ctx.get(callContextKey) as grpc.ServerUnaryCall<Request>;
        logger.log(level, "rpc finished", {
          ...call.metadata.getMap(),
          peer: call.getPeer(),
          rpc,
          start,
          end,
          reqtime_ms: end - start,
          code,
          ...(error ? { error: String(error) } : undefined),
        });
      });
  };
}

const infoCodes = [
  grpc.status.OK,
  grpc.status.CANCELLED,
  grpc.status.INVALID_ARGUMENT,
  grpc.status.NOT_FOUND,
  grpc.status.ALREADY_EXISTS,
  grpc.status.UNAUTHENTICATED,
];
const warnCodes = [
  grpc.status.DEADLINE_EXCEEDED,
  grpc.status.PERMISSION_DENIED,
  grpc.status.RESOURCE_EXHAUSTED,
  grpc.status.FAILED_PRECONDITION,
  grpc.status.ABORTED,
  grpc.status.OUT_OF_RANGE,
  grpc.status.UNAVAILABLE,
];
const errorCodes = [
  grpc.status.UNKNOWN,
  grpc.status.UNIMPLEMENTED,
  grpc.status.INTERNAL,
  grpc.status.DATA_LOSS,
];

function codeToLevel(code: grpc.status): string {
  if (infoCodes.includes(code)) {
    return "info";
  } else if (warnCodes.includes(code)) {
    return "warn";
  } else if (errorCodes.includes(code)) {
    return "error";
  } else {
    return "error";
  }
}
