# Fela

# Install

1. Install [Go](https://golang.org/), [Node.js](https://nodejs.org/en/),
   [Yarn](https://yarnpkg.com/),
   [MariaDB](https://mariadb.com/),
   [ffmpeg](https://ffmpeg.org/) and
   [Dbmate](https://github.com/amacneil/dbmate).
2. Create a new database user and database and adjust the `.env` file.
3. Run `dbmate up` to run the migrations.

# Run

1. Run `run.sh` which will start the server on port 80.

# Login

When running migrations an admin user account is created, which can be used to
login for the first time and create additional accounts. Remember to change the
password!

```
Initial login

Name: 		admin
Password:	password
```
