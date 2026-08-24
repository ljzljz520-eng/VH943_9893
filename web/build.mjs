import { mkdir, writeFile } from "node:fs/promises";

await mkdir("dist", { recursive: true });
await writeFile("dist/index.html", `<!doctype html>
<html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Campgear Staff Console</title></head>
<body><main><h1>Campgear Staff Console</h1><p>Use the Go staff API to manage inventory, maintenance, and rentals.</p></main></body></html>
`);

