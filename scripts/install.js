#!/usr/bin/env node

const fs = require("fs");
const path = require("path");

const binDir = path.join(__dirname, "..", "bin");
const ext = process.platform === "win32" ? ".exe" : "";
const bin = path.join(binDir, "echotik" + ext);

fs.mkdirSync(binDir, { recursive: true });

if (fs.existsSync(bin)) {
  process.exit(0);
}

console.warn(
  "No prebuilt echotik binary was bundled with this package.\n" +
    "For local development, run: go build -o bin/echotik .\n" +
    "For release, upload platform binaries and extend scripts/install.js to download them."
);
