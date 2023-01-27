package repository

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const userTable = `CREATE TABLE IF NOT EXISTS "users" (
	"id"				INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
	"username"			TEXT UNIQUE NOT NULL,
	"password"			TEXT NOT NULL,
	"email"				TEXT UNIQUE NOT NULL
);`

const postTable = `CREATE TABLE IF NOT EXISTS "posts" (
	"id"			INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
	"title"			TEXT NOT NULL,
	"author_id"		INTEGER NOT NULL,
	"author"		TEXT NOT NULL,
	"message"		TEXT NOT NULL,
	"likes" 		INTEGER DEFAULT 0,
	"dislikes" 	INTEGER DEFAULT 0,
	"category_id"	INTEGER NOT NULL,
  date DATETIME DEFAULT NULL,
	FOREIGN KEY(author_id) REFERENCES "users"(id), 
	FOREIGN KEY(category_id) REFERENCES "categories"(id) 
);`

const commentTable = `CREATE TABLE IF NOT EXISTS "comments" (
	"id"			INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
	"author_id"		INTEGER NOT NULL,
	"author"		TEXT NOT NULL,
	"likes" 		INTEGER DEFAULT 0,
	"dislikes" 		INTEGER DEFAULT 0,
	"post_id"		INTEGER NOT NULL,
	"message"		TEXT NOT NULL,
	"date"		DATETIME DEFAULT NULL,
	FOREIGN KEY(author_id) REFERENCES "users"(id),
	FOREIGN KEY(post_id) REFERENCES "posts"(id)
);`

const sessionTable = `CREATE TABLE IF NOT EXISTS "sessions" (
	"id"		INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
	"user_id"	INTEGER NOT NULL,
	"uuid"		TEXT NOT NULL,
	"created_at"	DATETIME DEFAULT NULL,
	"expires_at"	DATETIME DEFAULT NULL,
	FOREIGN KEY(user_id) REFERENCES "users"(id) ON DELETE CASCADE
);`

const categoryTable = `CREATE TABLE IF NOT EXISTS "categories" (
	"id"	INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE NOT NULL,
	"postid" INTEGER,
	"tag"	TEXT NOT NULL,
	"theme" TEXT NOT NULL DEFAULT ' ',
	"description"	TEXT NOT NULL DEFAULT ' '
);`

const likesTable = `CREATE TABLE IF NOT EXISTS "likes" (
	"user_id" INTEGER,
	"post_id" INTEGER DEFAULT NULL,
	"comment_id" INTEGER DEFAULT NULL
);`

const dislikesTable = `CREATE TABLE IF NOT EXISTS "dislikes" (
	"user_id" INTEGER,
	"post_id" INTEGER DEFAULT NULL,
	"comment_id" INTEGER DEFAULT NULL
);`

var tables = []string{userTable, postTable, commentTable, sessionTable, categoryTable, likesTable, dislikesTable}

func Init() (*sql.DB, error) {
	var err error

	db, err := sql.Open("sqlite3", "Forum.db")
	if err != nil {
		log.Println("❌ error | can't create DB")
		return nil, err
	}
	return db, nil
}

func CreateDatabase(db *sql.DB) error {
	for _, table := range tables {
		_, err := db.Exec(table)
		if err != nil {
			return err
		}
	}
	return nil
}
