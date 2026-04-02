import fs from "fs/promises"
import path from "path"
import { promisify } from "util"
import { exec } from "child_process"

import chokidar from "chokidar"

import consts from "./consts.mjs"

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
	const sidebarTsPath = path.join(sidebarTsDir, "sidebar.ts")
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

	await fs.mkdir(path.dirname(outPath), { recursive: true, cleanup: false })
	await fs.writeFile(outPath, code, "utf8")
	console.log(`✔ Overwrote JS sidebar: ${path.basename(outPath)}`)
}

async function prebuildPlugin(pluginId, { onlyWatch = false } = {}) {
	const ROOT = process.cwd()

	const { versionsPath, pluginConfig } = await getConfig(pluginId)

	if (!pluginConfig) {
		process.exit(1)
	}

	const specPath = pluginConfig.specPath
	const files = []

	let versionsArray = await getVersionsAsync(versionsPath)

	const generateDir = pluginConfig.outputDir
	const outputDir = path.relative(consts.GENERATED_PREFIX, generateDir)

	// ensure outputDir exists
	await fs.mkdir(outputDir, { recursive: true })

	for (const v of versionsArray) {
		files.push(path.resolve(ROOT, v.specPath))

		const genDir = path.resolve(consts.GENERATED_PREFIX, v.outputDir)
		const dir = path.resolve(ROOT, v.outputDir)

		try {
			await fs.rm(genDir, { recursive: true, force: true })
			await fs.mkdir(genDir, { recursive: true })

			await fs.mkdir(dir, { recursive: true })
			const processFile = path.join(dir, ".process")
			await fs.writeFile(processFile, "")
		} catch {}
	}

	if (!onlyWatch) {
		// generate all versioned API docs
		for (const v of versionsArray) {
			console.log(`\nGenerating versioned API docs: ${v.version}...`)

			const genDir = path.resolve(consts.GENERATED_PREFIX, v.outputDir)
			const dir = path.resolve(ROOT, v.outputDir)

			const processFile = path.join(dir, ".process")

			try {
				await fs.rm(genDir, { recursive: true, force: true })
				await fs.mkdir(genDir, { recursive: true })

				await execAsync(
					`npm run docusaurus gen-api-docs:version ${versionsPath}:${v.version}`,
				)

				const segments = consts.GENERATED_PREFIX.split(path.sep).filter(Boolean)
				const idPath = segments.slice(1).join(path.sep)

				await overwriteSidebarJson(
					pluginId,
					genDir,
					path.resolve(
						ROOT,
						`${pluginId}_versioned_sidebars/version-${v.version}-sidebars.json`,
					),
					`${idPath}/${pluginId}_versioned_docs/version-${v.version}/`,
				)

				await fs.cp(genDir, dir, {
					errorOnExist: false,
					force: true,
					recursive: true,
				})

				await fs.rm(processFile, { force: true })

				console.log(`✔ Version ${v.version} docs generated`)
			} catch (err) {
				console.error(
					`Error generating version ${v.version}:`,
					err.stderr || err,
				)
			}
		}
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
	files.push(path.resolve(ROOT, specPath))

	if (!onlyWatch) {
		console.log(`\nGenerating NEXT API docs for ${pluginId}...`)
		try {
			await fs.rm(generateDir, { recursive: true, force: true })
			await fs.mkdir(generateDir, { recursive: true })

			await fs.mkdir(outputDir, { recursive: true })

			await execAsync(`npm run docusaurus gen-api-docs ${versionsPath}`)

			const { sidebarPath } = await getPluginOptions(pluginId)

			await overwriteSidebarJs(
				pluginId,
				generateDir,
				path.resolve(ROOT, sidebarPath),
			)

			await fs.cp(generateDir, outputDir, {
				errorOnExist: false,
				force: true,
				recursive: true,
			})

			console.log("✔ NEXT docs generated")
		} catch (err) {
			console.error("Error generating NEXT docs:", err.stderr || err)
		}
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

	console.log(`\nPrebuild complete for ${pluginId}`)

	if (onlyWatch && files.length != 0) {
		const watcher = chokidar.watch(files)

		watcher.on("change", async () => {
			await prebuildPlugin(pluginId)
		})
	}

	console.log(`\nDone.`)
}

const pluginId = process.argv[2]

if (!pluginId) {
	console.error("Usage: npm run prebuild-api <pluginId>")
	process.exit(1)
}

let onlyWatch = false

if (process.argv.length > 2) {
	const flags = process.argv.slice(3)

	onlyWatch = flags.includes("--only-watch")
}

prebuildPlugin(pluginId, { onlyWatch })
