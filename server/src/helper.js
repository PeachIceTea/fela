import argon2 from "argon2"
import { milliseconds } from "date-fns"
import path from "node:path"

import config from "./config.js"

/**
 * Resolves the path relative to "filepath".
 * @param {string} filepath Path of the file as returned by "import.meta.url".
 * @param  {string[]} paths Any additional paths.
 * @returns {string} Absolute path.
 */
export const resolvePath = (filepath, ...paths) =>
	path.join(new URL(filepath).pathname, "..", ...paths)

/**
 * Sets the session cookie.
 * @param {import("express").Response} response
 * @param {string} sessionID
 */
export const setSessionCookie = (response, sessionID) => {
	response.cookie("session", sessionID, {
		httpOnly: true,
		maxAge: milliseconds({ days: config.server.session_length }),
	})
}

/**
 * Hashes a password using argon2.
 * @param {string} password Plaintext password to hash.
 * @returns {string} Hashed password.
 */
export const hashPassword = password =>
	argon2.hash(password, {
		type: argon2.argon2id,
		memoryCost: config.argon2.memory << 10,
		timeCost: config.argon2.iterations,
		parallelism: config.argon2.parallelism,
		saltLength: config.argon2.salt_length,
		hashLength: config.argon2.hash_length,
	})

const libraryStatus = new Set(["listening", "finished"])
/**
 * Checks if a given string is a valid library status.
 * @param {string} status String to check,
 * @returns {boolean} Boolean indicating if the given string is a valid status.
 */
export const isValidLibraryStatus = status => libraryStatus.has(status)
