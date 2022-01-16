import { database } from "./connection.js"

/**
 * Get a single user by username. Includes user password for auth.
 */
export const getUserForLogin = (() => {
	const stmt = database.prepare(`
    SELECT 
        id, 
        name, 
        password
    FROM user
    WHERE name = @name
    `)
	return ({ name }) => stmt.get({ name })
})()

/**
 * Insert a user.
 */
export const insertUser = (() => {
	const stmt = database.prepare(`
    INSERT INTO user (name, password, role)
    VALUES (@name, @password, @role)
    `)
	return ({ name, password, role }) => stmt.run({ name, password, role })
})()
