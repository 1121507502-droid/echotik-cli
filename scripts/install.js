#!/usr/bin/env node

const fs = require("fs");
const os = require("os");
const path = require("path");
const crypto = require("crypto");
const { execFileSync } = require("child_process");

const pkg = require("../package.json");

const REPO = "1121507502-droid/echotik-cli";
const NAME = "echotik";
const VERSION = pkg.version.replace(/-.*$/, "");

const PLATFORM_MAP = {
  darwin: "darwin",
  linux: "linux",
  win32: "windows",
};

const ARCH_MAP = {
  x64: "amd64",
  arm64: "arm64",
};

const platform = PLATFORM_MAP[process.platform];
const arch = ARCH_MAP[process.arch];
const isWindows = process.platform === "win32";
const ext = isWindows ? ".exe" : "";

if (!platform || !arch) {
  console.error(`Unsupported platform: ${process.platform}/${process.arch}`);
  process.exit(1);
}

const archiveExt = isWindows ? ".zip" : ".tar.gz";
const suffix = `${platform}-${arch}`;
const archiveName = `${NAME}-v${VERSION}-${suffix}${archiveExt}`;
const releaseURL = `https://github.com/${REPO}/releases/download/v${VERSION}/${archiveName}`;

const binDir = path.join(__dirname, "..", "bin");
const dest = path.join(binDir, NAME + ext);

function run(cmd, args, options = {}) {
  return execFileSync(cmd, args, {
    stdio: options.stdio || "pipe",
    encoding: options.encoding,
    timeout: options.timeout || 120000,
    env: process.env,
  });
}

function download(url, outputPath) {
  const args = [
    "--fail",
    "--location",
    "--silent",
    "--show-error",
    "--connect-timeout",
    "10",
    "--max-time",
    "120",
    "--max-redirs",
    "3",
    "--output",
    outputPath,
    url,
  ];
  run("curl", args, { stdio: ["ignore", "ignore", "pipe"] });
}

function extract(archivePath, tmpDir) {
  if (isWindows) {
    const ps = [
      "-NoProfile",
      "-ExecutionPolicy",
      "Bypass",
      "-Command",
      `Expand-Archive -LiteralPath '${archivePath.replace(/'/g, "''")}' -DestinationPath '${tmpDir.replace(/'/g, "''")}' -Force`,
    ];
    run("powershell.exe", ps, { stdio: "inherit" });
    return;
  }
  run("tar", ["-xzf", archivePath, "-C", tmpDir], { stdio: "inherit" });
}

function findExtractedBinary(tmpDir) {
  const expected = `${NAME}-${suffix}${ext}`;
  const direct = path.join(tmpDir, expected);
  if (fs.existsSync(direct)) return direct;

  const fallback = path.join(tmpDir, NAME + ext);
  if (fs.existsSync(fallback)) return fallback;

  const entries = fs.readdirSync(tmpDir);
  for (const entry of entries) {
    if (entry.startsWith(NAME) && entry.endsWith(ext)) {
      return path.join(tmpDir, entry);
    }
  }
  throw new Error(`Could not find extracted ${NAME} binary in ${tmpDir}`);
}

function sha256(filePath) {
  const hash = crypto.createHash("sha256");
  hash.update(fs.readFileSync(filePath));
  return hash.digest("hex");
}

function install() {
  fs.mkdirSync(binDir, { recursive: true });

  if (fs.existsSync(dest)) {
    try {
      run(dest, ["--version"], { stdio: "ignore", timeout: 10000 });
      return;
    } catch (_) {
      fs.rmSync(dest, { force: true });
    }
  }

  const tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), "echotik-cli-"));
  const archivePath = path.join(tmpDir, archiveName);

  try {
    console.log(`Downloading ${releaseURL}`);
    download(releaseURL, archivePath);
    console.log(`Downloaded ${archiveName} (${sha256(archivePath).slice(0, 12)}...)`);

    extract(archivePath, tmpDir);
    const extracted = findExtractedBinary(tmpDir);
    fs.copyFileSync(extracted, dest);
    fs.chmodSync(dest, 0o755);
    console.log(`${NAME} v${VERSION} installed successfully`);
  } catch (error) {
    console.error(
      `Failed to install ${NAME} v${VERSION} for ${process.platform}/${process.arch}.\n` +
        `Expected release asset: ${releaseURL}\n` +
        `Cause: ${error.message}\n\n` +
        `If this version was just tagged, wait for GitHub Actions to finish the release workflow.`
    );
    process.exit(1);
  } finally {
    fs.rmSync(tmpDir, { recursive: true, force: true });
  }
}

install();
