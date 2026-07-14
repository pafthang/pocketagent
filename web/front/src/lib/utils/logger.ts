// Browser console logger.

type LogFn = (...args: unknown[]) => void;

interface Logger {
  trace: LogFn;
  debug: LogFn;
  info: LogFn;
  warn: LogFn;
  error: LogFn;
}

function log(level: "trace" | "debug" | "info" | "warn" | "error", args: unknown[]) {
  const message = args
    .map((a) => (a instanceof Error ? `${a.message}\n${a.stack}` : String(a)))
    .join(" ");
  // eslint-disable-next-line no-console
  console[level === "trace" ? "debug" : level](message);
}

export const logger: Logger = {
  trace: (...args) => log("trace", args),
  debug: (...args) => log("debug", args),
  info: (...args) => log("info", args),
  warn: (...args) => log("warn", args),
  error: (...args) => log("error", args),
};