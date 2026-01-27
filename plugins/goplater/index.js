import fs from "fs-extra"
import path from "path"

import { execFileSync } from "child_process"
import os from "os"

export default function goplater(
	file = {
		content,
		path,
	},
	options = {
		binary: "./node_modules/.bin/goplater",
	},
) {
	if (!options.binary) {
		throw new Error("`binary` option is required for remarkGoplater")
	}

	if (!file.content || !file.path) {
		return file.content
	}

	const tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), "goplater-"))
	const inputFile = file.path
	const outputFile = path.join(tmpDir, `output.${path.extname(file.path)}`)

	const args = ["template", inputFile, "-s", inputFile, "-o", outputFile]

	execFileSync(options.binary, args, {
		cwd: process.cwd(),
		encoding: "utf-8",
	})

	if (fs.pathExistsSync(outputFile)) {
		const processed = fs.readFileSync(outputFile, "utf-8")

		file.content = processed
	} else {
		console.warn(`Processed file not found: ${outputFile}, skipping`)
	}

	fs.removeSync(tmpDir)

	return file.content
}
