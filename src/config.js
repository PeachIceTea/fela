import fs from "node:fs"
import path from "node:path"
import toml from "toml"

import { resolvePath } from "./helper.js"

/**
 * Parse the "fela.config.toml" file and expose it to the rest of the
 * application.
 * @returns {Object} The parsed config.
 */
const loadConfig = () => {
	const base = toml.parse(
		fs.readFileSync(resolvePath(import.meta.url, "../fela.config.toml")),
	)
	base.database.path = resolveConfigPaths(base.database.path)
	return base
}

/**
 * Resolves paths in config
 * @param {string} relativePath Absolute path or path relative to the config file.
 * @returns {string} Absolute path.
 */
const resolveConfigPaths = relativePath => {
	if (path.isAbsolute(relativePath)) return relativePath
	return resolvePath(import.meta.url, "..", relativePath)
}

export default loadConfig()
