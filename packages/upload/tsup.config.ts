import { defineConfig } from "tsup";

export default defineConfig({
  // One entry per subpath export. Split rather than bundled into one file so a
  // web app never pulls the Expo backend into its bundle, and vice versa.
  entry: {
    index: "src/index.ts",
    web: "src/web.ts",
    expo: "src/expo.ts",
    react: "src/react.ts",
    ui: "src/ui.tsx",
    transport: "src/transport.ts",
  },
  format: ["esm", "cjs"],
  dts: true,
  clean: true,
  sourcemap: true,
  treeshake: true,
  // React is a peer, and expo-image-manipulator is loaded through a runtime
  // import so it must never be bundled or resolved at build time.
  external: ["react", "expo-image-manipulator"],
});
