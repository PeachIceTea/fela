import { getSession, updateSession } from "./database/session.js"
import { setSessionCookie } from "./helper.js"

/**
 * Middleware that extracts the session id from the session cookie, gets the
 * session info from the database and puts it into response.locals.session.
 * @param {import("express").Request} request
 * @param {import("express").Response} response
 * @param {import("express").NextFunction} next
 */
export const sessionMiddleware = (request, response, next) => {
	const { session: sessionID } = request.cookies

	// Simply do nothing if we don't have a session.
	if (!sessionID) return next()

	// Get session info from the database.
	const sessionRow = getSession({ id: sessionID })

	// Instruct client to discard the session cookie if the session can no
	// longer be found in database.
	if (!sessionRow) {
		response.clearCookie("session")
		return next()
	}

	// Refresh the last_updated column in the session. As long as a session is
	// actively used we don't want to destroy it.
	updateSession({ id: sessionID })

	// Update session cookies maxAge.
	setSessionCookie(response, sessionID)

	// Make session information available to other parts of the application.
	response.locals.session = sessionRow
	next()
}

/**
 * Middleware ensuring that a route is only accessed with a valid session.
 * @param {import("express").Request} request
 * @param {import("express").Response} response
 * @param {import("express").NextFunction} next
 */
export const requireAuthentication = (request, response, next) => {
	if (!response.locals.session) {
		return response.status(401).json({
			error: "Authentication is required to access this path.",
		})
	}
	next()
}

/**
 * Middleware ensuring that a route is only accessed by an admin.
 * @param {import("express").Request} request
 * @param {import("express").Response} response
 * @param {import("express").NextFunction} next
 */
export const requireAdmin = (request, response, next) => {
	if (response.locals.session?.role !== "admin") {
		return response.status(403).json({
			error: "You do not have the required role to access this path.",
		})
	}
	next()
}

/**
 * This middleware catches errors thrown by body-parser. Body-parser not handle
 * various exceptions itself, for example if you send maleformed JSON it will
 * call JSON.parse on it and simply have the expception bubble up. Body-parser
 * does not offer a way to register a error handler, because why would it
 * sending stack traces over the wire is totally expected behavior.
 *
 * https://github.com/expressjs/body-parser/issues/122
 */
export const catchBodyParserErrors = (() => {
	const bodyParserErrors = new Set([
		"encoding.unsupported",
		"entity.parse.failed",
		"entity.verify.failed",
		"request.aborted",
		"request.size.invalid",
		"stream.encoding.set",
		"parameters.too.many",
		"charset.unsupported",
		"encoding.unsupported",
		"entity.too.large",
	])

	return (error, request, response, next) => {
		if (error.expose) {
			/* eslint-disable unicorn/prefer-ternary */
			if (error && bodyParserErrors.has(error.type)) {
				return response.json({
					error: `Could not request body: ${error.message}.`,
				})
			} else {
				return next(error)
			}
		} else {
			return response.json({
				error: `Could not parse request body.`,
			})
		}
	}
})()
