import bodyParser from "body-parser"
import cookieParser from "cookie-parser"
import express from "express"
import morgan from "morgan"

import config from "./config.js"
import AudiobookController from "./controller/audiobook.js"
import FileController from "./controller/files.js"
import FSController from "./controller/fs.js"
import LibraryController from "./controller/library.js"
import UserController from "./controller/user.js"
import { database } from "./database/connection.js"
import { catchBodyParserErrors, sessionMiddleware } from "./middleware.js"

const app = express()

// Register logging middleware.
app.use(morgan("dev"))

// Register body parser.
app.use(bodyParser.json())
app.use(bodyParser.urlencoded({ extended: true }))
app.use(catchBodyParserErrors)

// Register cookie parser.
app.use(cookieParser())

// Register session middleware.
app.use(sessionMiddleware)

//Register controller.
app.use(FSController)
app.use(UserController)
app.use(AudiobookController)
app.use(FileController)
app.use(LibraryController)

// 404 handler.
app.use("*", (request, response) => {
	response.status(404).json({
		error: "Requested route could not be found.",
	})
})

// Start server.
const server = app.listen(config.server.port, config.server.host, () => {
	console.log(
		`Server listening at http://${config.server.host}:${config.server.port}.`,
	)
})

/**
 * Makes the server shutdown gracefully. Rather than exiting immediately the
 * close functions will block while sqlite and express are still processing.
 */
const shutdown = () => {
	server.close(() => {
		database.close()
		process.exit(0)
	})
}
process.on("exit", shutdown)
process.on("SIGTERM", shutdown)
process.on("SIGINT", shutdown)
