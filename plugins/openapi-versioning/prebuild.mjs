import fs from "fs/promises"
import path from "path"
import { promisify } from "util"
import { exec } from "child_process"

import { getConfig, getVersionsAsync, getPluginOptions } from "./index.mjs"

const execAsync = promisify(exec)

function replacePrefixRecursive(items, replacePrefix) {
	if (!Array.isArray(items)) return items
	return items.map((item) => {
		const newItem = { ...item }
		if (replacePrefix && typeof newItem.id === "string") {
			newItem.id = newItem.id.replace(replacePrefix, "")
		}
		if (newItem.items)
			newItem.items = replacePrefixRecursive(newItem.items, replacePrefix)
		if (newItem.link?.id) {
			newItem.link.id = newItem.link.id.replace(replacePrefix, "")
		}
		return newItem
	})
}

export async function mergeSidebarFromPlugin(
	pluginId,
	sidebarTsDir,
	replacePrefix,
) {
	const ROOT = process.cwd()

	const configModule = await import(path.join(ROOT, "docusaurus.config.js"))
	const docusaurusConfig = configModule.default

	const navbarItems = docusaurusConfig.themeConfig?.navbar?.items || []

	// find first sidebarId for pluginId
	let sidebarId = null
	for (const item of navbarItems) {
		if (
			item?.type === "docSidebar" &&
			item?.docsPluginId === pluginId &&
			item?.sidebarId
		) {
			sidebarId = item.sidebarId
			break
		}
	}

	if (!sidebarId) {
		console.error(`No docSidebar navbar item found for pluginId "${pluginId}"`)
		return {}
	}

	// load sidebar.ts
	const sidebarTsPath = path.join(ROOT, sidebarTsDir, "sidebar.ts")
	let sidebarTs = {}
	try {
		const module = await import(sidebarTsPath)
		sidebarTs = module.default || module
	} catch {
		console.warn(
			`sidebar.ts not found in ${sidebarTsDir}, returning empty object`,
		)
		return {}
	}

	const content = replacePrefixRecursive(sidebarTs, replacePrefix)

	return {
		[sidebarId]: content,
	}
}

export async function overwriteSidebarJson(
	pluginId,
	sidebarTsDir,
	outPath,
	replacePrefix,
) {
	const merged = await mergeSidebarFromPlugin(
		pluginId,
		sidebarTsDir,
		replacePrefix,
	)

	await fs.mkdir(path.dirname(outPath), { recursive: true })
	await fs.writeFile(outPath, JSON.stringify(merged, null, 2), "utf8")
	console.log(`✔ Overwrote JSON sidebar: ${path.basename(outPath)}`)
}

export async function overwriteSidebarJs(
	pluginId,
	sidebarTsDir,
	outPath,
	replacePrefix,
) {
	const merged = await mergeSidebarFromPlugin(
		pluginId,
		sidebarTsDir,
		replacePrefix,
	)

	const code =
		`// @ts-check\n\n` +
		`/**\n` +
		` @type {import('@docusaurus/plugin-content-docs').SidebarsConfig}\n` +
		` */\n` +
		`const sidebars = ${JSON.stringify(merged, null, 2)};\n\n` +
		`export default sidebars;\n`

	await fs.mkdir(path.dirname(outPath), { recursive: true })
	await fs.writeFile(outPath, code, "utf8")
	console.log(`✔ Overwrote JS sidebar: ${path.basename(outPath)}`)
}

async function prebuildPlugin(pluginId) {
	const ROOT = process.cwd()

	const { versionsPath, pluginConfig } = await getConfig(pluginId)

	if (!pluginConfig) {
		process.exit(1)
	}

	const outputDir = pluginConfig.outputDir

	let versionsArray = await getVersionsAsync(versionsPath)

	// ensure outputDir exists
	await fs.mkdir(outputDir, { recursive: true })

	// generate all versioned API docs
	for (const v of versionsArray) {
		console.log(`\nGenerating versioned API docs: ${v.version}...`)
		const dir = path.resolve(ROOT, v.outputDir)
		try {
			await fs.rm(dir, { recursive: true, force: true })
			await fs.mkdir(dir, { recursive: true })
			await fs.writeFile(path.join(dir, ".process"), "")
			await execAsync(
				`npm run docusaurus gen-api-docs:version ${versionsPath}:${v.version}`,
			)
			await overwriteSidebarJson(
				pluginId,
				v.outputDir,
				path.resolve(
					ROOT,
					`${pluginId}_versioned_sidebars/version-${v.version}-sidebars.json`,
				),
				`version-${v.version}/`,
			)
			console.log(`✔ Version ${v.version} docs generated`)
		} catch (err) {
			console.error(`Error generating version ${v.version}:`, err.stderr || err)
		}

		await fs.rm(path.join(dir, ".process"))
	}

	if (versionsArray.length != 0) {
		// save versions.json
		try {
			await fs.writeFile(versionsPath, JSON.stringify(versionsArray))
		} catch (err) {
			console.error(
				`Error saving versions file '${versionsPath}':`,
				err.stderr || err,
			)
		}
	}

	// generate NEXT API docs
	console.log(`\nGenerating NEXT API docs for ${pluginId}...`)
	try {
		await fs.rm(outputDir, { recursive: true, force: true })
		await fs.mkdir(outputDir, { recursive: true })
		await execAsync(`npm run docusaurus gen-api-docs ${versionsPath}`)
		const { sidebarPath } = await getPluginOptions(pluginId)
		await overwriteSidebarJs(
			pluginId,
			outputDir,
			path.resolve(ROOT, sidebarPath),
		)
		console.log("✔ NEXT docs generated")
	} catch (err) {
		console.error("Error generating NEXT docs:", err.stderr || err)
	}

	if (versionsArray.length != 0) {
		// save versions.json
		try {
			await fs.writeFile(versionsPath, JSON.stringify(versionsArray))
		} catch (err) {
			console.error(
				`Error saving versions file '${versionsPath}':`,
				err.stderr || err,
			)
		}
	}

	console.log(`\n✅ Prebuild complete for ${pluginId}`)
}

const pluginId = process.argv[2]

if (!pluginId) {
	console.error("Usage: npm run prebuild-api <pluginId>")
	process.exit(1)
}

prebuildPlugin(pluginId)
