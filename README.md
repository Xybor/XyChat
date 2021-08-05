# XyChat
A chat application built with Golang, VueJS, and PostgreSQL.

# REQUIREMENTS
+ Go1.16.6
+ PostgreSQL

# INSTALLATION
+ Create a database in your PostgreSQL.
+ Create and modify `.env` file (with the same structure as [.env.example](.env.example)).

# USAGE
Run the following command in terminal for more detail
```shell
$ go run main.go -h
```

## Some examples
Reset the database
```shell
$ go run main.go -reset
```

Create an admin account
```shell
$ go run main.go -admin root:p@ss
```

Run the application
```shell
$ go run main.go -run
```

Then open your browser and access the web at `http://domain:port/ui/`. Example: `http://localhost:1999/ui`.

# DOCKER

