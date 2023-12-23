/*
 * Vencord, a Discord client mod
 * Copyright (c) 2023 Vendicated and contributors
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

import { definePluginSettings } from "@api/Settings";
import definePlugin, { OptionType, PluginNative } from "@utils/types";

// KEEP THIS! Generates native code.
const Native = VencordNative.pluginHelpers.Core as PluginNative<typeof import("./native")>;

const settings = definePluginSettings({
	path: {
		type: OptionType.STRING,
		description: "Where to find Venjector",
		hidden: true,
		default: "/home/tizu/go/bin/venjector",
	},
	visualize: {
		type: OptionType.BOOLEAN,
		description: "Should injections be visualized?",
		default: false,
	},
});

export default definePlugin({
	name: "Venjector",
	description: "Adds Venjector stuff :3",
	authors: [
		{
			name: "tizu",
			id: 805510183736049684n,
		},
	],
	required: true,
	settings,

	patches: [],

	start() {
		if (!Native) return; // unused my ass
	}
});
