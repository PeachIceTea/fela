import express from "express"

import {
	addAudiobookToLibrary,
	deleteFromLibrary,
	getLibrary,
	updateLibraryProgress,
	updateLibraryStatus,
} from "../database/library.js"
import { isValidLibraryStatus } from "../helper.js"
import { requireAuthentication } from "../middleware.js"

const router = express.Router()

router.get("/library", requireAuthentication, (request, response) => {
	const userID = response.locals.session.id
	const library = getLibrary({ id: userID })
	response.json({ library })
})

router.post("/library/:bookID", requireAuthentication, (request, response) => {
	let { bookID } = request.params
	bookID = Number.parseInt(bookID)
	if (!bookID) {
		return response.status(400).json({
			error: "Invalid audiobook id.",
		})
	}
	try {
		addAudiobookToLibrary({ userID: response.locals.session.id, bookID })
	} catch (error) {
		console.error(error)
		// TODO: Catch UNIQUE errors and non existant bookIDs.
		return response.status(500).json({
			error: "Could not add audiobook to your library.",
		})
	}
	return response.json({ success: true })
})

router.put(
	"/library/:bookID/status",
	requireAuthentication,
	(request, response) => {
		let { bookID } = request.params
		bookID = Number.parseInt(bookID)
		if (!bookID) {
			return response.status(400).json({
				error: "Invalid audiobook id.",
			})
		}
		let { status } = request.body
		if (!status || !isValidLibraryStatus(status)) {
			return response.status(400).json({
				error: "Invalid status.",
			})
		}

		updateLibraryStatus({
			status,
			userID: response.locals.session.id,
			bookID,
		})
		return response.json({ success: true })
	},
)

router.put(
	"/library/:bookID/progress",
	requireAuthentication,
	(request, response) => {
		let { bookID, fileID, progress } = request.body
		bookID = Number.parseInt(bookID)
		fileID = Number.parseInt(fileID)
		if (!bookID || !fileID) {
			return response.status(400).json({
				error: "Invalid id.",
			})
		}
		progress = Number.parseFloat(progress)
		if (Number.isNaN(progress)) {
			return response.status(400).json({
				error: "Invalid progress.",
			})
		}
		updateLibraryProgress({
			fileID,
			progress,
			userID: response.locals.session.id,
			bookID,
		})
		return response.json({ success: true })
	},
)

router.delete(
	"/library/:bookID",
	requireAuthentication,
	(request, response) => {
		let { bookID } = request.body
		bookID = Number.parseInt(bookID)
		if (!bookID) {
			return response.status(400).json({
				error: "Invalid audiobook id.",
			})
		}
		deleteFromLibrary({ userID: response.locals.session.id, bookID })
		return response.json({ success: true })
	},
)

export default router
