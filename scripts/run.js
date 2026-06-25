#!/usr/bin/env node

const { execFileSync } = require("child_process");
const path = require("path");

const ext = process.platform === "win32" ? ".exe" : "";
const bin = path.join(__dirname, "..", "bin", "echotik" + ext);

try {
  execFileSync(bin, process.argv.slice(2), { stdio: "inherit" });
} catch (error) {
  if (error.code === "ENOENT") {
    console.error("echotik binary is missing. Build it locally or publish release artifacts first.");
  }
  process.exit(error.status || 1);
}
