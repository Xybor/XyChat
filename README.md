# XyChat
A chat application built with Golang, VueJS, and PostgreSQL.

# REQUIREMENTS
+ Go1.16.6
+ PostgreSQL

# INSTALLATION
+ Create a database in your PostgreSQL.
+ Set environment variables with the same name as [.env.example](./.env.example) file.

# USAGE
Clone the repository
```shell
$ git clone https://github.com/Xybor/XyChat.git
$ cd XyChat
```

Run directly
```shell
Xychat $ go run main.go -h
```
or build and run
```shell
Xychat $ go build main.go
Xychat $ main -h
```

## Some examples
Reset the database
```shell
$ go run main.go -reset
# or
$ main -reset
```

Create an admin account
```shell
$ go run main.go -admin root:p@ss
# or
$ main -admin root:p@ss
```

Use values in a .env file instead of environment variables
```shell
$ go run main.go -dotenv
# or
$ main -dotenv
```

Run the application
```shell
$ go run main.go -run
# or
$ main -run
```
You can combine all above options together.  
Now open your browser and access the web at `http://domain:port/ui/`. Example: [http://localhost:1999/ui](http://localhost:1999/ui).

# DEPLOY ON HEROKU
Clone the repository
```shell
$ git clone https://github.com/Xybor/XyChat.git
$ cd XyChat
```

Create a Heroku app
```shell
(Xychat) $ heroku create
```

Create a `Procfile` file, then add the `web` command as the following text (you can edit the command by adding more options as aforementioned)
```py
web: ./bin/xychat -run
```
Don't use `-dotenv` option, using [.env]() file isn't recommended on Heroku. Then commit your change
```shell
(Xychat) $ git add Procfile
(Xychat) $ git commit -m "Add the Procfile" 
```  
Setup PostgreSQL on Heroku
```shell
(Xychat) $ heroku addons:create heroku-postgresql:hobby-dev
(Xychat) $ heroku config:set DSN_NAME=DATABASE_URL
```

Setup domain and port
```shell
(Xychat) $ heroku config:set DOMAIN=xxxx.heroku.com PORT=1999
```

Push your repo to heroku
```shell
(Xychat) $ git push heroku main:main
```

If lucky, there is no error when you are installing, you can open the app now:
```shell
(Xychat) $ heroku open
```

# DOCKER

