# HOW TO DEPLOY LOCALLY

## REQUIREMENTS

- Go1.16.6
- PostgreSQL
- Git

## Step 1. Create a database in your PostgreSQL

## Step 2. Setup environment variables
For further details, see [EVAR.md](./EVAR.md)
Method 1: Set environment variables with the same name as [.env.example](./.env.example) file. Example:

```shell
# Example on Windows
$ set DOMAIN=localhost
$ set PORT=1999
```

Method 2: Create a file named .env with the same structure as [.env.example](./.env.example) file.

## Step 3. Clone the repository

```shell
$ git clone https://github.com/Xybor/XyChat.git
$ cd XyChat
```

## Step 4. Build (optional)

```shell
Xychat$ go build
```

## Step 5. Run

Start the application by running the main.exe (the output file from go build)

```shell
Xychat$ main -run
```

If you didn't build the source code, you could run the application from `go` command

```shell
Xychat$ go run main.go -run
```

Use `-dotenv` option if you use .env file instead of setting environment variables

```shell
Xychat$ main -run -dotenv
# or
Xychat$ go run main.go -run -dotenv
```

For further details in command options, run:

```shell
Xychat$ main -h
# or
Xychat$ go run main.go -h
```

## Step 6. Access the application
See the [FEATURES.md](./FEATURES.md) for all APIs of the application.

## Some option examples

Reset the database

```shell
Xychat$ main -reset
```

Create an admin account with the username and password is stored in ADMIN
environment variable.  The form of that variable must be `<username>:<password>`

```shell
Xychat$ main -admin ADMIN
```

Use the .env file instead of environment variables

```shell
Xychat$ main -dotenv
```

Only run the application

```shell
Xychat$ main -run
```

You could combine all above options together.

# HOW TO DEPLOY ON HEROKU

## REQUIREMENTS

- Git
- Heroku

## Step 1. Clone the repository

```shell
$ git clone https://github.com/Xybor/XyChat.git
$ cd XyChat
```

## Step 2. Create a Heroku app

```shell
Xychat$ heroku create
```

## Step 3. Create Procfile

Create a file named `Procfile`, then add the `web` command as the following text (you can edit the command by adding more options as aforementioned)

```
web: ./bin/xychat -run
```

Don't use `-dotenv` option, using [.env]() file isn't recommended on Heroku. Then commit your change

```shell
Xychat$ git add Procfile
Xychat$ git commit -m "Add the Procfile"
```

## Step 4. Setup PostgreSQL

```shell
Xychat$ heroku addons:create heroku-postgresql:hobby-dev
```

## Step 5. Setup config var

```shell
Xychat$ heroku config:set DSN_NAME=DATABASE_URL
Xychat$ heroku config:set DOMAIN=XXX.herokuapp.com PORT=1999
```

## Step 6. Push to and deploy on Heroku

```shell
Xychat$ git push heroku main:main
```

## Step 7. Access the application

If lucky, there is no error when you are installing, you can access the application
now by the domain on herokuapp.com. See all APIs in [FEATURES.md](./FEATURES.md)