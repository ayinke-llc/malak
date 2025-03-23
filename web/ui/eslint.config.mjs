import { defineConfig } from "eslint/config";
import globals from "globals";
import tseslint from "typescript-eslint";
import pluginReact from "eslint-plugin-react";
import reactHooks from "eslint-plugin-react-hooks";

const queryPlugins = {
  "react-query-keys": await import("eslint-plugin-react-query-keys"),
  "react-query-must-invalidate-queries": await import("eslint-plugin-react-query-must-invalidate-queries")
};

export default defineConfig([
  {
    files: ["**/*.{js,mjs,cjs,ts,jsx,tsx}"],
    languageOptions: {
      globals: {
        ...globals.browser
      },
      parser: tseslint.parser,
      parserOptions: {
        ecmaVersion: "latest",
        sourceType: "module",
        ecmaFeatures: {
          jsx: true
        }
      }
    },
    plugins: {
      "@typescript-eslint": tseslint.plugin,
      "react": pluginReact,
      "react-hooks": reactHooks,
      "react-query-keys": queryPlugins["react-query-keys"].default,
      "react-query-must-invalidate-queries": queryPlugins["react-query-must-invalidate-queries"].default
    },
    rules: {
      "no-unused-vars": "off",
      "@typescript-eslint/no-unused-vars": [
        "warn",
        {
          "argsIgnorePattern": "^_"
        }
      ],
      "react-hooks/rules-of-hooks": "error",
      "react-hooks/exhaustive-deps": "warn",
      "react-query-keys/no-plain-query-keys": "warn",
      "react-query-must-invalidate-queries/require-mutation-invalidation": "warn"
    }
  }
]);
