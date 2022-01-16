import express from "express"
import fs from "node:fs/promises"
import path from "node:path"

import { requireAdmin } from "../middleware.js"

const router = express.Router()

/**
 * GET /fs
 *
 * Lists the content of a given directory.
 */
router.get("/fs", requireAdmin, async (request, response) => {
	let { path: queryPath } = request.query
	queryPath = queryPath
		? path.resolve(queryPath)
		: "/mnt/c/users/jonas/nextcloud/ebooks and audiobooks"

	const directoryListing = []
	try {
		for (const entry of await fs.readdir(queryPath, {
			withFileTypes: true,
		})) {
			directoryListing.push({
				name: entry.name,
				path: path.join(queryPath, entry.name),
				isDirectory: entry.isDirectory(),
			})
		}
	} catch {
		return response.json({ error: "Could not read directory." })
	}
	directoryListing.sort((a, b) => {
		if (a.isDirectory && !b.isDirectory) {
			return -1
		} else if (b.isDirectory && !a.isDirectory) {
			return 1
		} else {
			return a.name.localeCompare(b.name)
		}
	})

	response.json({
		directory: {
			path: queryPath,
			entries: [...createDotDirectories(queryPath), ...directoryListing],
		},
	})
})

/**
 * Returns directory listing for the . and .. directories for a given path.
 * @param {string} queryPath Path of the . directory.
 * @returns Directory listing of the . and .. directories.
 */
const createDotDirectories = queryPath => {
	return [
		{ name: ".", path: queryPath, isDirectory: true },
		{ name: "..", path: path.join(queryPath, ".."), isDirectory: true },
	]
}

export default router
