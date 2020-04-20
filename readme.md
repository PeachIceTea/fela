# Fela

# Install

1. Install [Go](https://golang.org/), [Node.js](https://nodejs.org/en/),
[MariaDB](https://mariadb.com/) and
[dbmate](https://github.com/amacneil/dbmate).
2. Create a new database user and database and adjust the `.env` file.
3. Run `dbmate up` to run the migrations.
4. Run `npm install` inside `./client` to install the required node modules.

# Run

1. To serve the javascript client, run `npm run serve` in the `./client` folder.
2. To run the Go server, run `go run .` in the root directory

# Login

When running migrations an admin user account is created, which can be used to
login for the first time and create additional account. Remember to change the
password!

```
Initial login

Name: 		admin
Password:	password
```
