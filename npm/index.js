#!/usr/bin/env node

import { execFileSync } from "child_process";
import os from "os";
import { fileURLToPath } from "url";

function getPlatform() {
    const platform = os.platform();

    switch (platform) {
        case "win32":
            return "windows";
        case "darwin":
        case "linux":
            return platform;
        default:
            console.error(`Unsupported platform: ${platform}`);
            process.exit(1);
    }
}

function getArch() {
    let arch = os.arch();

    switch (arch) {
        case 'arm64':
            return 'arm64_v8.0';
        case 'x64':
            return 'amd64_v1';
        default:
            console.error(`Unsupported architecture: ${arch}`);
            process.exit(1);
    }
}

(function() {
    const platform = getPlatform();
    const arch = getArch();

    const pathToBinary = fileURLToPath(new URL(`./zeabur_${platform}_${arch}/zeabur${platform === "windows" ? ".exe" : ""}`, import.meta.url));
    const args = process.argv.slice(2);

    // This wrapper must be transparent: the binary already wrote its own
    // stdout/stderr (stdio: "inherit"), so on failure we only mirror its exit
    // status. Letting execFileSync throw would print a Node stack trace with
    // `stdout: null, stderr: null` and replace the real exit code with 1.
    try {
        execFileSync(pathToBinary, args, { stdio: "inherit" });
    } catch (e) {
        if (typeof e.status === "number") {
            process.exit(e.status);
        }
        if (e.signal) {
            // Shell convention for a process killed by a signal.
            process.exit(128 + (os.constants.signals[e.signal] ?? 0));
        }
        console.error(e.message ?? String(e));
        process.exit(1);
    }
})()
