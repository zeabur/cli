#!/usr/bin/env node

import { execSync } from "child_process";
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

    execSync(`${pathToBinary} ${args.join(" ")}`, { stdio: "inherit" });
})()
