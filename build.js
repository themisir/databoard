require("esbuild").build({
  entryPoints: ["dashboard/index.tsx"],
  bundle: true,
  minify: true,
  sourcemap: true,
  outfile: "dashboard/public/bundle.js",
});
