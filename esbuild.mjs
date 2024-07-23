import * as esbuild from 'esbuild'
import postCssPlugin from 'esbuild-style-plugin';
import tailwind from 'tailwindcss';
import autoprefixer from 'autoprefixer';
import { writeFileSync } from "node:fs";

const result = await esbuild.build({
    entryPoints: [
        'assets/main.ts',
        'assets/nutrition.ts',
    ],
    bundle: true,
    minify: true,
    outdir: 'dist/assets',
    metafile: true,
    entryNames: '[dir]/[name]-[hash]',
    sourcemap: true,
    plugins: [
        postCssPlugin({
            postcss: {
                plugins: [tailwind, autoprefixer],
            },
        }),
    ],
    assetNames: '[dir]/[name]-[hash]'
});

const outputFiles = Object.entries(result.metafile.outputs)
    .map(([filename, details]) => {
        return filename.slice(5);
    });

const encoder = new TextEncoder();

writeFileSync("dist/manifest.json", encoder.encode(JSON.stringify({ outputFiles })));
