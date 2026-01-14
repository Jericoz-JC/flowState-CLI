#!/usr/bin/env node

const https = require("https");
const fs = require("fs");
const path = require("path");
const { execSync } = require("child_process");
const zlib = require("zlib");

const VERSION = "0.1.1";
const REPO = "Jericoz-JC/flowState-CLI";

// Map Node.js platform/arch to GoReleaser naming
function getPlatformInfo() {
  const platform = process.platform;
  const arch = process.arch;

  const platformMap = {
    darwin: "darwin",
    linux: "linux",
    win32: "windows",
  };

  const archMap = {
    x64: "amd64",
    arm64: "arm64",
  };

  const goPlatform = platformMap[platform];
  const goArch = archMap[arch];

  if (!goPlatform || !goArch) {
    console.error(`Unsupported platform: ${platform}-${arch}`);
    process.exit(1);
  }

  return { platform: goPlatform, arch: goArch, isWindows: platform === "win32" };
}

function getDownloadUrl(platform, arch) {
  const ext = platform === "windows" ? "zip" : "tar.gz";
  return `https://github.com/${REPO}/releases/download/v${VERSION}/flowstate-${platform}-${arch}.${ext}`;
}

function download(url) {
  return new Promise((resolve, reject) => {
    const request = (url) => {
      https.get(url, (response) => {
        if (response.statusCode === 302 || response.statusCode === 301) {
          // Follow redirect
          request(response.headers.location);
          return;
        }

        if (response.statusCode !== 200) {
          reject(new Error(`Failed to download: ${response.statusCode} ${url}`));
          return;
        }

        const chunks = [];
        response.on("data", (chunk) => chunks.push(chunk));
        response.on("end", () => resolve(Buffer.concat(chunks)));
        response.on("error", reject);
      }).on("error", reject);
    };

    request(url);
  });
}

async function extractTarGz(buffer, destDir) {
  const gunzip = zlib.createGunzip();
  const { Readable } = require("stream");

  // Write to temp file and use tar command (simpler than implementing tar in JS)
  const tempFile = path.join(destDir, "temp.tar.gz");
  fs.writeFileSync(tempFile, buffer);

  try {
    execSync(`tar -xzf "${tempFile}" -C "${destDir}"`, { stdio: "pipe" });
  } finally {
    fs.unlinkSync(tempFile);
  }
}

async function extractZip(buffer, destDir) {
  // Write to temp file and use PowerShell to extract (Windows)
  const tempFile = path.join(destDir, "temp.zip");
  fs.writeFileSync(tempFile, buffer);

  try {
    execSync(
      `powershell -command "Expand-Archive -Path '${tempFile}' -DestinationPath '${destDir}' -Force"`,
      { stdio: "pipe" }
    );
  } finally {
    fs.unlinkSync(tempFile);
  }
}

async function main() {
  const { platform, arch, isWindows } = getPlatformInfo();
  const url = getDownloadUrl(platform, arch);
  const binDir = path.join(__dirname, "bin");
  const binaryName = isWindows ? "flowstate.exe" : "flowstate";
  const binaryPath = path.join(binDir, binaryName);

  console.log(`Downloading flowstate v${VERSION} for ${platform}-${arch}...`);
  console.log(`URL: ${url}`);

  try {
    const buffer = await download(url);

    // Ensure bin directory exists
    if (!fs.existsSync(binDir)) {
      fs.mkdirSync(binDir, { recursive: true });
    }

    // Extract based on platform
    if (isWindows) {
      await extractZip(buffer, binDir);
    } else {
      await extractTarGz(buffer, binDir);
    }

    // Make binary executable on Unix
    if (!isWindows) {
      fs.chmodSync(binaryPath, 0o755);
    }

    console.log(`flowstate installed successfully to ${binaryPath}`);
  } catch (error) {
    console.error("Installation failed:", error.message);
    console.error("");
    console.error("You can manually download from:");
    console.error(`  https://github.com/${REPO}/releases/tag/v${VERSION}`);
    process.exit(1);
  }
}

main();

