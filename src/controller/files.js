import express from "express"

import { getFilePath } from "../database/book.js"
import { resolvePath } from "../helper.js"
import { requireAuthentication } from "../middleware.js"

const router = express.Router()

// Server covers.
router.use(
	"/file/cover",
	requireAuthentication,
	express.static(resolvePath(import.meta.url, "../../cover"), {}),
)

// Serve audio file.
router.get(
	"/file/audio/:fileID",
	requireAuthentication,
	(request, response) => {
		let fileID = Number.parseInt(request.params.fileID)
		if (!fileID) {
			return response.status(400).json({
				error: "Invalid file id.",
			})
		}

		const filePath = getFilePath({ fileID })
		if (!filePath) {
			return response.status(404).json({
				error: "No file with that id.",
			})
		}

		response.sendFile(filePath)
	},
)

export default router
