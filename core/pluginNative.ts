/*
 * Vencord, a Discord client mod
 * Copyright (c) 2023 Vendicated and contributors
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

import { IpcMainInvokeEvent } from "electron";

import { spawn } from "child_process";

export async function run(
  _: IpcMainInvokeEvent,
  p: string,
  n: string,
  v: boolean
) {
  console.log(
    "Command to be executed:",
    [
      p,
      "--auto-choice=" + n,
      v ? "--visual=true" : "--visual=false",
      n != "-1" ? "--tipless=true" : "--tipless=false",
    ].join(" ")
  );

  const child = spawn(p, [
    "--auto-choice=" + n,
    v ? "--visual=true" : "--visual=false",
    n != "-1" ? "--tipless=true" : "--tipless=false",
  ]);

  child.stderr.on("data", (data) => {
    console.log(`stderr: ${data}`);
  });
  child.stdout.on("data", (data) => {
    console.log(`stdout: ${data}`);
  });

  const exitCode = await new Promise((resolve, reject) => {
    child.on("close", resolve);
  });
  console.log(`exit code: ${exitCode}`);

  return exitCode === 0;
}
