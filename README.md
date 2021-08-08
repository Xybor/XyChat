# XyChat
A chat application built with Golang, VueJS, and PostgreSQL.

# REQUIREMENTS
+ Go1.16.6
+ PostgreSQL

# INSTALLATION
+ Create a database in your PostgreSQL.
+ Create and modify `.env` file (with the same structure as [.env.example](.env.example)). The
other way is that create environment variables in your machine with the same names as [.env.example](.env.example).

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

Use values in .env file instead of environment variables
```shell
$ go run main.go -dotenv
```

Run the application
```shell
$ go run main.go -run
```

Then open your browser and access the web at `http://domain:port/ui/`. Example: `http://localhost:1999/ui`.

# DEVCONTAINER
The Visual Studio Code Remote - Containers extension lets you use a Docker container as a full-featured development environment. It allows you to open any folder inside (or mounted into) a container and take advantage of Visual Studio Code's full feature set.

# VUEJS BUILD

```
cd vue
```

Install node packages
```
npm install
```

For develop
```
npm run serve
```

For deploy
```
npm run build
``
