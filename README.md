# Fastabiz Mini RDBMS

**A lightweight in-memory relational database management system implemented in Go.**  

This project is my submission for the **Pesapal Junior Developer Challenge ’26**.  
It demonstrates building a functional mini-RDBMS from scratch with a **SQL-like interface** and an **interactive REPL**.

---

## Table of Contents

- [Features](#features)  
- [Example Usage](#example-usage)  
  - [Creating Tables](#creating-tables)  
  - [CRUD Operations](#crud-operations)  
  - [Join Example](#join-example)  
- [Getting Started](#getting-started)  
- [Notes](#notes)  
- [Future Improvements](#future-improvements)  
- [Acknowledgements](#acknowledgements)  

---

## Features

- **Create tables** with primary keys  
- **CRUD operations**: `INSERT`, `SELECT`, `UPDATE`, `DELETE`  
- **Basic indexing** for fast primary key lookups  
- **Simple joins** between tables  
- **Interactive REPL** for executing SQL-like commands  
- Supports **string** and **integer** column types  
- **In-memory storage** — lightweight and easy to experiment with  

---

## Example Usage

### Creating Tables

```sql
-- Create a users table
CREATE TABLE users (id INT PRIMARY KEY, name TEXT);

-- Insert data
INSERT INTO users (id, name) VALUES (1, 'John');

-- Query table
SELECT id, name FROM users;
-- Output: { id: 1, name: John }

-- Update data
UPDATE users SET name = 'Jane' WHERE id = 1;

-- Delete data
DELETE FROM users WHERE id = 1;

-- Join users with orders
SELECT users.id, users.name, orders.id AS order_id
FROM users
JOIN orders ON users.id = orders.user_id;
```

## Getting Started
Follow these steps to clone the repository and run the REPL:

**Clone the Repository**
```bash
bash/zsh

# via SSH
git clone git@github.com:JohnKibet/Fastabiz-mini-rdbms.git
cd Fastabiz-mini-rdbms

# or via HTTPS
git clone https://github.com/JohnKibet/Fastabiz-mini-rdbms.git
cd Fastabiz-mini-rdbms

Run the REPL 
go run ./mini-db/main.go
```

## Notes
- This project is for demonstration and learning purposes as part of a coding challenge.
- The RDBMS runs entirely in-memory; there is no persistent storage yet.

## Future Improvements
- Support for additional data types
- More advanced SQL-like features
- Disk-backed persistence

## Acknowledgements
- Fully implemented from scratch in Go
- Inspired by standard relational database concepts and my personal project Fastabiz