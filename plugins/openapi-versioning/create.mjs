import fs from "fs/promises"
import path from "path"
import { getConfig, getPluginOptions } from "./index.mjs"

const ROOT = process.cwd()

function parseArg() {
	const args = process.argv
	if (args.length != 4) {
		console.error("Usage: node ./create.mjs <pluginId> <version>")
		process.exit(1)
	}

	return { pluginId: args[2], version: args[3] }
}

async function duplicateSpec(pluginConfig, version) {
	const specPath = pluginConfig.specPath
	const absSpecPath = path.join(ROOT, specPath)

	const specDir = path.dirname(absSpecPath)
	const ext = path.extname(specPath)
	const baseName = path.basename(specPath, ext)

	const versionedFilename = `${baseName}-${version}${ext}`
	const versionedPath = path.join(specDir, versionedFilename)

	await fs.copyFile(absSpecPath, versionedPath)
	console.log(`✔ Created ${versionedFilename}`)

	return {
		versionedFilename,
		relativeSpecPath: path
			.join(path.dirname(specPath), versionedFilename)
			.replace(/\\/g, "/"),
	}
}

async function updateVersionsJson(
	versionsPath,
	pluginId,
	version,
	relativeSpecPath,
) {
	let versions = []

	try {
		const raw = await fs.readFile(versionsPath, "utf8")
		const parsed = JSON.parse(raw)
		if (Array.isArray(parsed)) {
			versions = parsed
		}
	} catch {
		// file may not exist
	}

	if (!versions.some((v) => v.version === version)) {
		versions.push({
			version,
			specPath: relativeSpecPath,
			outputDir: `${pluginId}_versioned_docs/version-${version}`,
			label: version,
			baseUrl: `/${pluginId}_versioned_docs/version-${version}`,
		})

		await fs.writeFile(versionsPath, JSON.stringify(versions, null, 2))
		console.log("✔ Updated versions.json")
	}
}

async function updateApiVersions(version) {
	const apiVersionsPath = path.join(ROOT, "api_versions.json")

	let apiVersions = []

	try {
		const raw = await fs.readFile(apiVersionsPath, "utf8")
		const parsed = JSON.parse(raw)
		if (Array.isArray(parsed)) {
			apiVersions = parsed
		}
	} catch {
		// file may not exist
	}

	if (!apiVersions.includes(version)) {
		apiVersions.push(version)
		await fs.writeFile(apiVersionsPath, JSON.stringify(apiVersions, null, 2))
		console.log(`✔ Added ${version} to ${path.basename(apiVersionsPath)}`)
	}
}

async function main() {
	const { pluginId, version } = parseArg()

	const { versionsPath, pluginConfig } = await getConfig(pluginId)
	if (!pluginConfig) process.exit(1)

	const { versionedFilename, relativeSpecPath } = await duplicateSpec(
		pluginConfig,
		version,
	)

	await updateVersionsJson(versionsPath, pluginId, version, relativeSpecPath)

	await updateApiVersions(version)

	console.log("\nDone.")
}

main()
