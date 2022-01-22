import express from "express"
import { exec } from "node:child_process"
import { constants as fsConstants } from "node:fs"
import fs from "node:fs/promises"

import {
	createBook,
	getAudiobook,
	getAudiobooks,
	insertFile,
} from "../database/book.js"
import { database } from "../database/connection.js"
import { requireAdmin, requireAuthentication } from "../middleware.js"

const router = express.Router()

router.get("/audiobook", requireAuthentication, (request, response) => {
	const books = getAudiobooks()
	response.json({ books })
})

router.get("/audiobook/:bookID", requireAuthentication, (request, response) => {
	const bookID = Number.parseInt(request.params.bookID)
	if (!bookID) {
		return response.status(400).json({
			error: "Invalid audiobook id.",
		})
	}

	const book = getAudiobook({ bookID: bookID })
	if (!book) {
		return response.status(404).json({
			error: "No audiobook with that id.",
		})
	}

	return response.json({ book })
})

// POST /audiobook
// Create an audiobook entry in the database.
router.post("/audiobook", requireAdmin, async (request, response) => {
	const { files } = request.body
	if (!files || !Array.isArray(files) || !isArrayOfStrings(files)) {
		return response.status(400).json({
			error: "Invalid or missing data.",
		})
	}

	// Test if given files exists and we have read access.
	try {
		const queue = []
		for (const filePath of files) {
			// fs.access throws if we don't have access and resolves without
			// value if we do. I just want to document my dislike of that API.
			queue.push(fs.access(filePath, fsConstants.R_OK))
		}
		await Promise.all(queue)
	} catch (error) {
		console.error(error)
		return response.status(400).json({
			error: "Could not access file.",
		})
	}

	// Create file objects for the database. This also ensures that all files are actually audio
	// files.
	let fileObjects
	try {
		const fileObjectsQueue = []
		for (const [index, filePath] of files.entries()) {
			fileObjectsQueue.push(createFileObject(filePath, index + 1))
		}
		fileObjects = await Promise.all(fileObjectsQueue)
	} catch (error) {
		console.error(error)
		return response.status(500).json({
			error: "Could not parse files.",
		})
	}

	// Create book, we need it to insert the files.
	const bookInfo = await getInfo(files[0])

	// Insert data into the database.
	let bookID = 0
	try {
		database.transaction(() => {
			// Create new book in the database.
			bookID = createBook(bookInfo)

			// Try to extract a cover from the first file.
			extractAlbumArt(files[0], bookID)

			// Insert files into the
			for (const fileObject of fileObjects) {
				fileObject.book = bookID
				insertFile(fileObject)
			}
		})()
	} catch (error) {
		if (error.message.includes("UNIQUE")) {
			return response.status(409).json({
				error: "One or more files already are part of an audiobook.",
			})
		}

		return response.status(500).json({
			error: "Could not create book.",
		})
	}

	return response.json({ bookID: bookID })
})

/**
 * Create a object describing the file, to be inserted into the database.
 * @param {string} filePath Path to the file.
 * @param {number} position Position of the file in the audiobook.
 * @returns {{path: string, duration: number, position: number}} The object.
 */
const createFileObject = async (filePath, position) => ({
	path: filePath,
	duration: await getDuration(filePath),
	position,
})

/**
 * Uses ffprobe to extract the duration of the audiostream of an audiofile.
 * @param {string} filePath Path of an audiofile.
 * @returns {number} Duration of the audiostream.
 */
const getDuration = async filePath => {
	const streams = JSON.parse(
		await runCommand(
			`ffprobe -v quiet -print_format json -show_streams -i "${filePath}"`,
		),
	).streams
	for (const stream of streams) {
		if (stream.codec_type === "audio") {
			const duration = Number.parseFloat(stream.duration)
			if (duration) {
				return duration
			} else {
				throw "Could not extract non-zero duration from audio stream."
			}
		}
	}
	throw "Could not find audio stream."
}

/**
 * Tries to extract audiobook information from the given file. This isn't an
 * exact science, but better than nothing.
 * @param {string} filePath Path to the file.
 * @returns {{title: string, author: string}} The extracted audiobook info.
 */
const getInfo = async filePath => {
	const tags = JSON.parse(
		await runCommand(
			`ffprobe -v quiet -print_format json -show_format -i "${filePath}"`,
		),
	).format?.tags

	let author = ""
	let title = ""
	if (!tags) return { author, title }

	// In general the album should be the title of the book. If the book is a
	// single file, there are cases where the title is the title of the book.
	if (tags.album) {
		title = tags.album
	} else if (tags.title) {
		title = tags.title
	}

	// There really is only one place to look for the author.
	if (tags.artist) {
		author = tags.artist
	}

	return { author, title }
}

const extractAlbumArt = async (filePath, bookID) => {
	try {
		await runCommand(
			`ffmpeg -i "${filePath}" -an -vcodec copy cover/${bookID}.jpg`,
		)
		// eslint-disable-next-line no-empty
	} catch {}
}

/**
 * Runs a commandline command.
 * @param {string} command Command to run.
 * @returns {string} The stdout of the command. Throws with error or stderr.
 */
const runCommand = command =>
	new Promise((resolve, reject) =>
		exec(command, (error, stdout, stderr) => {
			if (error) {
				return reject(error)
			}

			if (stderr) {
				return reject(stderr)
			}

			resolve(stdout)
		}),
	)

/**
 * Tests if an array contains only strings.
 * @param {any[]} array The array to be tested.
 * @returns {boolean} Whether array contains only strings.
 */
const isArrayOfStrings = array => {
	for (const entry of array) {
		if (typeof entry !== "string") {
			return false
		}
	}
	return true
}
export default router
