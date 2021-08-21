# ENVIRONMENT VARIABLE GUIDELINE
There are two ways of setting environment variables in xychat application, including system environment variable and `.env` file.
## System environment variable
You can set environment variables by following commands in cmd (Windows) or terminal (Linux):
```shell
# Windows
$ set NAME=value
# Linux
$ export NAME=value
```
To unset a variable, use:
```shell
# Windows
$ set NAME=
# Linux
$ unset NAME
```
To apply these changes permantly, please search on internet for further details.

## .env file
You can use a file to declare environment variable instead of setting on system directly. Create a file named `.env` and put enviroment variables as following:
```
NAME1=value1
NAME2=value2
``` 
See an example in [.env.example](./.env.example). When you run the application, please include `-dotenv` option in the command.


## List of environment variables
Note that all variables must be uppercase
|Name|Optional|Meaning|Example|
|----|--|-------|-------|
|XYCHAT|yes|Which state the application is running for, test, debug, or release. The default value is debug.|debug|
|DOMAIN||The mainly hosting domain|localhost|
|PORT||The current hosting port|8080|
|POSTGRES_HOST|yes|The hostname of postgreSQL database|localhost|
|POSTGRES_PORT|yes|The port of portgreSQL database|5432|
|POSTGRES_USER|yes|The user of postgreSQL database|postgres|
|POSTGRES_PASSWORD|yes|The password of current postgreSQL's user|password|
|POSTGRES_DBNAME|yes|The database name in your postgreSQL|xychat|
|DSN_NAME|yes|The reflected variable of another variable including DSN string. Note that if you don't provide POSTGRES_* variables, please declare this one.|DATABASE_URL
|CORS|yes|All website urls which this application allows to request cross-site. Seperate them by semicolons. If you don't declare this variable, all urls will be allowed.|http://your.first.domain; https://your.second.domain|
|TLS|yes|See [TLS.md](./TLS.md) for futher details|xychat|
