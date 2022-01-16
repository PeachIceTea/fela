import { formatISO, sub } from "date-fns"
import crypto from "node:crypto"
import { scheduleJob } from "node-schedule"

import config from "../config.js"
import { database } from "./connection.js"

/**
 * Get a single session.
 */
export const getSession = (() => {
	const stmt = database.prepare(`
    SELECT u.id, u.name, u.role
    FROM session s
	JOIN user u ON s.user = u.id
    WHERE s.id = @id
    `)
	return data => stmt.get(data)
})()

/**
 * Insert a session.
 */
export const insertSession = (() => {
	const stmt = database.prepare(`
    INSERT INTO session (id, user)
    VALUES (@id, @user)
    `)
	return data => stmt.run(data)
})()

/**
 * Creates a new session for a user.
 * @param {number} userID ID of the user the session belongs to.
 * @returns Session ID.
 */
export const createSession = userID => {
	const sessionID = crypto.randomBytes(33).toString("base64url")
	insertSession({ id: sessionID, user: userID })
	return sessionID
}

/**
 * Update a sessions last_used.
 */
export const updateSession = (() => {
	const stmt = database.prepare(`
    UPDATE session
    SET last_used = CURRENT_TIMESTAMP
    WHERE id = @id
    `)
	return data => stmt.run(data)
})()

/**
 * Delete a session. Equivalent to logout.
 */
export const deleteSession = (() => {
	const stmt = database.prepare(`
    DELETE FROM session
    WHERE id = @id
    `)
	return data => stmt.run(data)
})()

/**
 * Delete expired sessions that are still in the database.
 */
const deleteExpiredSessions = (() => {
	const stmt = database.prepare(`
    DELETE FROM session
    WHERE last_used < @expireDate 
    `)

	return () => {
		stmt.run({
			expireDate: formatISO(
				sub(new Date(), {
					days: config.server.session_length,
				}),
			),
		})
	}
})()

scheduleJob("delete expired sessions", "0 0 * * *", deleteExpiredSessions)
