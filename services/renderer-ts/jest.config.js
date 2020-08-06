"use strict";

module.exports = {
  roots: ["./src"],
  testMatch: ["**/*.spec.ts"],
  reporters: ["default"],
  collectCoverage: false, // カバレッジが欲しいときは --coverage で有効化する
  collectCoverageFrom: ["./src/**/*.ts", "!./src/**/*.spec.ts"],
  coverageDirectory: "coverage",
  globals: {
    "ts-jest": {
      tsConfig: "./tsconfig.test.json",
    },
  },
  transform: {
    "\\.ts$": "ts-jest",
  },
};
