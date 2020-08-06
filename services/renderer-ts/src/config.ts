export type Config = Readonly<{
  mode: string;
  grpcPort: number;
  gracefulStopTimeout: number;
}>;

/**
 * 環境変数から設定を読み込む
 */
export function loadConfig(): Config {
  const config = {
    mode: "production",
    grpcPort: 50051,
    gracefulStopTimeout: 10 * 1000,
  };

  // mode
  if (process.env.MODE) {
    config.mode = process.env.MODE;
  }

  // grpcPort
  if (process.env.GRPC_PORT) {
    const grpcPort = parseInt(process.env.GRPC_PORT, 10);
    if (Number.isNaN(grpcPort)) {
      throw new Error(`GRPC_PORT is invalid: ${process.env.GRPC_PORT}`);
    }
    config.grpcPort = grpcPort;
  }

  // gracefulStopTimeout
  if (process.env.GRACEFUL_STOP_TIMEOUT) {
    const gracefulStopTimeout = parseInt(process.env.GRACEFUL_STOP_TIMEOUT, 10);
    if (Number.isNaN(gracefulStopTimeout)) {
      throw new Error(`GRACEFUL_STOP_TIMEOUT is invalid: ${process.env.GRACEFUL_STOP_TIMEOUT}`);
    }
    config.gracefulStopTimeout = gracefulStopTimeout;
  }

  return config;
}
