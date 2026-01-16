#!/usr/bin/env node

const https = require("https");
const fs = require("fs");
const path = require("path");
const { execSync } = require("child_process");
const zlib = require("zlib");

// Read version from package.json to keep in sync
const packageJson = require("./package.json");
const VERSION = packageJson.version;
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

  // Handle 32-bit Node.js on 64-bit systems (common Windows issue)
  if (!goArch && (arch === "ia32" || arch === "x86")) {
    console.error("");
    console.error("================================================================");
    console.error("  ERROR: 32-bit Node.js detected (architecture: " + arch + ")");
    console.error("================================================================");
    console.error("");
    console.error("  flowstate-cli requires 64-bit Node.js.");
    console.error("");
    console.error("  You're likely running 32-bit Node.js on a 64-bit system.");
    console.error("  This commonly happens when the wrong installer was downloaded.");
    console.error("");
    console.error("  To fix this:");
    console.error("    1. Uninstall your current Node.js");
    console.error("    2. Download 64-bit Node.js from: https://nodejs.org/");
    console.error("       - Windows: Choose 'Windows Installer (.msi)' for 64-bit");
    console.error("       - Look for 'x64' in the filename");
    console.error("    3. Install and retry: npm install -g flowstate-cli");
    console.error("");
    console.error("  Alternatively, download the binary directly from:");
    console.error("    https://github.com/" + REPO + "/releases/tag/v" + VERSION);
    console.error("");
    process.exit(1);
  }

  if (!goPlatform || !goArch) {
    console.error("");
    console.error("Unsupported platform: " + platform + "-" + arch);
    console.error("");
    console.error("flowstate-cli supports:");
    console.error("  - Windows (x64, arm64)");
    console.error("  - macOS (x64 Intel, arm64 Apple Silicon)");
    console.error("  - Linux (x64, arm64)");
    console.error("");
    console.error("You can request support for your platform at:");
    console.error("  https://github.com/" + REPO + "/issues");
    console.error("");
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

    console.log("");
    console.log("================================================================");
    console.log("  flowstate v" + VERSION + " installed successfully!");
    console.log("================================================================");
    console.log("");
    console.log("  To run flowstate:");
    console.log("    flowstate");
    console.log("");
    if (isWindows) {
      console.log("  If 'flowstate' is not recognized, try:");
      console.log("    npx flowstate");
      console.log("");
      console.log("  Or add npm global bin to your PATH:");
      console.log("    1. Run: npm config get prefix");
      console.log("    2. Add the returned path to your system PATH");
      console.log("");
    } else {
      console.log("  If 'flowstate' command is not found, try:");
      console.log("");
      console.log("    Option 1 - Use npx (always works):");
      console.log("      npx flowstate");
      console.log("");
      console.log("    Option 2 - Add npm bin to PATH (recommended):");
      console.log("      # For bash/zsh, add to ~/.bashrc or ~/.zshrc:");
      console.log("      export PATH=\"$(npm config get prefix)/bin:$PATH\"");
      console.log("");
      console.log("      # Then reload your shell:");
      console.log("      source ~/.bashrc  # or source ~/.zshrc");
      console.log("");
      console.log("    Option 3 - Run directly:");
      console.log("      $(npm config get prefix)/bin/flowstate");
      console.log("");
    }
  } catch (error) {
    console.error("Installation failed:", error.message);
    console.error("");
    console.error("Troubleshooting:");
    console.error("  1. Check your internet connection");
    console.error("  2. Verify the release exists at:");
    console.error(`     https://github.com/${REPO}/releases/tag/v${VERSION}`);
    console.error("");
    console.error("Manual installation:");
    console.error(`  - Download: flowstate-${platform}-${arch}.${platform === "windows" ? "zip" : "tar.gz"}`);
    console.error(`  - Extract to: ${binDir}`);
    if (!isWindows) {
      console.error(`  - Run: chmod +x ${binaryPath}`);
    }
    console.error("");
    console.error("Or install via Go:");
    console.error("  go install github.com/Jericoz-JC/flowState-CLI/cmd/flowState@latest");
    process.exit(1);
  }
}

main();

