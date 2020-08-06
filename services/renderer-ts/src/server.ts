import health from "grpc-health-check";
import winston from "winston";
import { RenderRequest, RenderReply } from "../pb/renderer/renderer_pb";
import { IRendererServer } from "../pb/renderer/renderer_grpc_pb";
import { UnaryHandler, unaryHandler, createLoggingMiddleware } from "./server-utils";
import { render } from "./renderer";

export interface Server extends IRendererServer {
  healthServer: health.IHealthServer;
}

/**
 * サーバー (IRendererServer) を作成する
 */
export function createServer(logger: winston.Logger): Server {
  const healthServer = new health.Implementation({
    "": health.messages.HealthCheckResponse.ServingStatus.SERVING,
    "renderer.Renderer": health.messages.HealthCheckResponse.ServingStatus.SERVING,
  });
  const loggingMiddleware = createLoggingMiddleware(logger);
  return {
    healthServer,
    render: unaryHandler(handleRender, loggingMiddleware("renderer.Renderer/render")),
  };
}

/**
 * renderer.Renderer/Render に対するハンドラ
 */
export const handleRender: UnaryHandler<RenderRequest, RenderReply> = async (req) => {
  const src = req.getSrc();
  const html = await render(src);
  const reply = new RenderReply();
  reply.setHtml(html);
  return reply;
};
