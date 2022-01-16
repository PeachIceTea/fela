import argon2 from "argon2"
import express from "express"

import { createSession } from "../database/session.js"
import { getUserForLogin, insertUser } from "../database/user.js"
import { hashPassword, setSessionCookie } from "../helper.js"
import { requireAdmin } from "../middleware.js"

const router = express.Router()

router.post("/users/login", async (request, response) => {
	let { name, password } = request.body
	name = name.trim()
	if (!name || !password) {
		return response.json({
			error: "Username or password missing.",
		})
	}

	const userRow = getUserForLogin({ name })
	if (!userRow) {
		return response.json({
			error: "Could not find user.",
		})
	}

	if (!(await argon2.verify(userRow.password, password))) {
		return response.json({
			error: "Passwords do not match.",
		})
	}

	const sessionID = createSession(userRow.id)
	setSessionCookie(response, sessionID)
	return response.json({ msg: "Succesfully logged in." })
})

router.post("/users/create", requireAdmin, async (request, response) => {
	let { name, password, role } = request.body
	if (typeof name !== "string" || typeof password !== "string") {
		return response.status(400).json({
			error: "Missing or invalid username or password.",
		})
	}

	name = name.trim().toLocaleLowerCase()
	if (!name || !password) {
		return response.status(400).json({
			error: "Username or password missing.",
		})
	}

	password = await hashPassword(password)

	try {
		insertUser({
			name,
			password,
			role: role ? role : "default",
		})
	} catch (error) {
		if (error.message.includes("UNIQUE")) {
			return response
				.status(409)
				.json({ error: "User with that username already exists." })
		}
		return response
			.status(500)
			.json({ error: "Could not write to database." })
	}

	return response.json({ msg: "User was created." })
})

export default router
