import fs from "fs"
import path from "path"

import docusaurusConfig from "../../docusaurus.config.js"

export default function parse({
	specPath,
	outputDir,
	versionsPath,
	...options
}) {
	const ROOT = process.cwd()

	const resolvedVersionsPath = versionsPath
		? path.resolve(ROOT, versionsPath)
		: path.resolve(ROOT, outputDir, "versions.json")

	const versionsArray = getVersions(resolvedVersionsPath)
	const versions = {}

	for (const entry of versionsArray) {
		if (!entry?.version) continue

		versions[entry.version] = {
			specPath: entry.specPath,
			outputDir: entry.outputDir,
			label: entry.label,
			baseUrl: entry.baseUrl,
		}
	}

	return {
		[resolvedVersionsPath]: {
			specPath,
			outputDir,
			label: "next",
			version: "next",
			baseUrl: `/${outputDir.replace(/\/$/, "")}`,
			versions,
			...options,
		},
	}
}

export async function getPluginOptions(pluginId) {
	try {
		for (const plugin of docusaurusConfig.plugins || []) {
			if (Array.isArray(plugin) && plugin[1]?.id === pluginId) {
				return plugin[1]
			}
		}

		console.error(`Plugin with id "${pluginId}" not found`)
		return null
	} catch (e) {
		console.error("Failed to load docusaurus.config.js", e)
		return null
	}
}

export function getVersions(versionsPath) {
	try {
		const raw = fs.readFileSync(versionsPath, "utf8")
		const parsed = JSON.parse(raw)

		return Array.isArray(parsed) ? parsed : []
	} catch {
		return []
	}
}

export async function getVersionsAsync(versionsPath) {
	try {
		const raw = await fs.promises.readFile(versionsPath, "utf8")
		const parsed = JSON.parse(raw)

		return Array.isArray(parsed) ? parsed : []
	} catch {
		return []
	}
}

export async function getConfig(pluginId) {
	if (!Array.isArray(docusaurusConfig.plugins)) {
		console.error("No plugins array found in docusaurus.config.js")
		return null
	}

	// Find plugin by id
	let pluginConfig

	for (const plugin of docusaurusConfig.plugins) {
		if (Array.isArray(plugin) && plugin[1]?.docsPluginId === pluginId) {
			pluginConfig = plugin[1]
			break
		}
	}

	if (!pluginConfig?.config) {
		console.error(`no config for "${pluginId}" found`)
		return null
	}

	let versionsPath = (Object.keys(pluginConfig.config) || [null])[0]

	if (!versionsPath) {
		console.error(`missing versionsPath for plugin "${pluginId}"`)
		return null
	}

	if (!pluginConfig.config[versionsPath]?.specPath) {
		console.error(`specPath not found for plugin "${pluginId}"`)
		return null
	}

	return { versionsPath, pluginConfig: pluginConfig?.config[versionsPath] }
}
