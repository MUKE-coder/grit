import { defineConfig } from "vitest/config";

export default defineConfig({
  test: {
    // jsdom for the component test. The logic tests do not need a DOM but are
    // cheap enough to run under one.
    environment: "jsdom",
    globals: false,
  },
  esbuild: { jsx: "automatic" },
  resolve: {
    /**
     * Without this, React hooks throw "Cannot read properties of null" the
     * moment a component renders.
     *
     * Under pnpm's strict layout @testing-library/react brings its own nested
     * react-dom, which then talks to a different React instance than the one
     * beside it. Two copies of React in one tree means the dispatcher is null,
     * and the error it produces points at useId rather than at the real cause.
     */
    dedupe: ["react", "react-dom"],
  },
});
