Fastabiz Mini RDBMS

A simple relational database management system implemented in Go.

This project is my submission for the Pesapal Junior Developer Challenge '26. It demonstrates building a working mini-RDBMS from scratch with a SQL-like interface and REPL.

Features

Table creation with primary keys

CRUD operations: INSERT, SELECT, UPDATE, DELETE

Basic indexing using primary keys for fast lookups

Simple joins between tables

Interactive REPL for executing SQL-like commands

String and integer support for columns

Lightweight in-memory storage, suitable for learning and demonstration

Example Usage
-- create table
CREATE TABLE users (id INT PRIMARY KEY, name TEXT);

-- insert data
INSERT INTO users (id, name) VALUES (1, 'John');

-- query table
SELECT * FROM users;
-- { id:1 name:John }

-- update data
UPDATE users SET name = 'Jane' WHERE id = 1;

-- delete data
DELETE FROM users WHERE id = 1;

How to Run
# clone the repo
git clone <your-repo-url>
cd fastabiz-mini-rdbms

# run the REPL
go run ./mini-db/main.go

Notes

This project is for demonstration and learning purposes, as part of a coding challenge.

The RDBMS is in-memory; no persistent storage yet.

Further enhancements could include additional data types, more SQL features, and persistent disk storage.

Acknowledgements

Implemented from scratch in Go.

Inspired by standard relational database concepts and personal project fastabiz.