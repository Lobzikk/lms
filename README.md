## Encryption

This project doesn't store users' passwords via the plain text, but stores the encrypted (hs256 signing method) sum of key and password in `users` table - that way even in the case of `users` table leak their data will be completely safed, and no one could ever sign in using their data!

## Autosave

The server autosaves all the data it has every 10 seconds on PostgreSQL DB. That's it.

## How to run the server?

1. Run `setup.bat` and follow instructions
2. Run `go mod tidy`
3. Use .sql files in `sql_stuff/sql_run'n'go` to create needed tables in your PostgreSQL database 
4. Run `start.bat` and enjoy!

## How to use the server? (GET request instructions)

* "/signup" - Sign up a new user
  * password - strong password
  * key - second part of the encrypted hash sum
  * username - user name
* "/signin" - Sign in, parameters are the same as when signing up
* "/" - Get all the expressions assigned to your account
* "/add" - Add new expression
  * exp - needed expression
* "/signout" - Sign out

---

To get more information about the project - take a look at scheme.jpg!