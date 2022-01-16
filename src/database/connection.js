import sqlite from "better-sqlite3"
import fs from "node:fs"

import config from "../config.js"
import { resolvePath } from "../helper.js"

// Open database file.
export const database = sqlite(config.database.path)

// Enable WAL mode and foreign keys.
database.pragma("journal_mode = WAL")
database.pragma("foreign_keys = ON")

// Ensure schema.
database.exec(
	fs.readFileSync(resolvePath(import.meta.url, "../schema.sql"), "utf-8"),
)
