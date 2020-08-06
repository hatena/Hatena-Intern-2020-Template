"use strict";

module.exports = {
  plugins: ["prettier"],
  overrides: [
    {
      files: ["src/**/*.ts"],
      extends: [
        "eslint:recommended",
        "plugin:@typescript-eslint/recommended",
        "prettier",
        "prettier/@typescript-eslint",
        "plugin:node/recommended-module",
        "plugin:eslint-comments/recommended",
      ],
      parser: "@typescript-eslint/parser",
      parserOptions: {
        ecmaVersion: 2019,
        sourceType: "module",
        project: "./tsconfig.build.json",
      },
      env: {
        es6: true,
        node: true,
      },
      rules: {
        "prettier/prettier": "error",
        "node/no-unsupported-features/es-syntax": "off",
        "node/no-missing-import": [
          "error",
          {
            tryExtensions: [".js", ".ts"],
          },
        ],
      },
    },
    {
      files: ["src/**/*.spec.ts"],
      extends: ["plugin:jest/recommended", "plugin:jest-formatting/recommended"],
      parserOptions: {
        project: "./tsconfig.test.json",
      },
      env: {
        "jest/globals": true,
      },
    },
    {
      files: ["*.js"],
      extends: [
        "eslint:recommended",
        "prettier",
        "plugin:node/recommended-script",
        "plugin:eslint-comments/recommended",
      ],
      parserOptions: {
        ecmaVersion: 2019,
        sourceType: "script",
      },
      env: {
        es6: true,
        node: true,
      },
      rules: {
        "prettier/prettier": "error",
      },
    },
  ],
};
