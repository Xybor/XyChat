# XyChat
A API-base chat application at server-side

# REQUIREMENTS
+ Go1.16.6
+ PostgreSQL

# INSTALLATION
+ Modify `.env` file with your database credentials ('postgres_*').
+ Create a database with the same name in your PostgreSQL.

# USAGE
Run the command in terminal or command line
```
$ go run .
```
Then open your browser, request the API by using following URLs:
+ `register`: `http://localhost:1999/api/v1/register/?username=<USN>&password=<PWD>`
+ `authenticate`: `http://localhost:1999/api/v1/auth/?username=<USN>&password=<PWD>`
+ `profile`: `http://localhost:1999/api/v1/profile/?token=<TOKEN>`

# DOCKER

