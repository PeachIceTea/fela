import { database } from "./connection.js"

export const getLibrary = (() => {
	const stmt = database.prepare(`
    SELECT 
        b.*, 
        l.status, 
        l.file, 
        l.progress
    FROM library l
    JOIN book b ON l.book = b.id
    WHERE l.user = @id
    `)
	return data => stmt.all(data)
})()

export const addAudiobookToLibrary = (() => {
	const stmt = database.prepare(`
    INSERT INTO library (user, book, file)
    VALUES (
        @userID, 
        @bookID, 
        (
            SELECT id 
            FROM file 
            WHERE 
                book = @bookID 
                AND position = 1
        ))
    `)
	return data => stmt.run(data)
})()

export const updateLibraryStatus = (() => {
	const stmt = database.prepare(`
    UPDATE library
    SET status = @status
    WHERE 
        user = @userID
        AND book = @bookID
    `)

	return ({ status, bookID, userID }) => stmt.run({ status, bookID, userID })
})()

export const updateLibraryProgress = (() => {
	const stmt = database.prepare(`
    UPDATE library
    SET 
        file = @fileID,
        progress = @progress
    WHERE 
        user = @userID
        AND book = @bookID
    `)

	return ({ fileID, progress, userID, bookID }) =>
		stmt.run({ fileID, progress, userID, bookID })
})()

export const deleteFromLibrary = (() => {
	const stmt = database.prepare(`
    DELETE FROM library
    WHERE 
        user = @userID
        AND book = @bookID
    `)
	return ({ userID, bookID }) => stmt.run({ userID, bookID })
})()
