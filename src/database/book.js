import { database } from "./connection.js"

export const getAudiobooks = (() => {
	const stmt = database.prepare(`
    SELECT *
    FROM book
    `)

	return () => stmt.all()
})()

export const getAudiobook = (() => {
	const bookStmt = database.prepare(`
    SELECT *
    FROM book
    WHERE id = @bookID
    `)
	const fileStmt = database.prepare(`
    SELECT
        id,
        duration,
        position
    FROM file
    WHERE book = @bookID
    `)

	return ({ bookID }) => {
		const book = bookStmt.get({ bookID })
		if (book) {
			book.files = fileStmt.all({ bookID })
		}
		return book
	}
})()

export const createBook = (() => {
	const stmt = database.prepare(`
    INSERT INTO book(title, author)
    VALUES (@title, @author)
    `)
	return ({ title, author }) => stmt.run({ title, author }).lastInsertRowid
})()

export const insertFile = (() => {
	const stmt = database.prepare(`
    INSERT INTO file (path, duration, position, book)
    VALUES (@path, @duration, @position, @book)
    `)
	return ({ path, duration, position, book }) =>
		stmt.run({ path, duration, position, book }).lastInsertRowid
})()

export const getFilePath = (() => {
	const stmt = database.prepare(`
    SELECT path
    FROM file
    WHERE id = @id
    `)
	return ({ fileID }) => stmt.get({ fileID }).path
})()
