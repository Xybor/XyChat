# Notations
+ `%p`: `p` is an optional parameter.

# API
## v1
|Method|URL          |Query Parameters|Behavior             |
|------|-------------|----------------|---------------------|
|GET   |`register`   |`%token, %role, username, password`|Register a user with a specific role|
|GET   |`auth`       |`username, password`|Get an authenticated token|
|GET   |`profile`    |`%token`|Get the profile|
|GET   |`users/:id`  |`%token`|Get user's profile by userid|
|PUT   |`users/:id`  |`%token, %age, %gender`|Update user's age, gender by userid|
|PUT   |`users/:id/role`|`%token, role`|Update user's role by userid|
|PUT   |`users/:id/password`|`%token, %oldpassword, newpassword`|Update user's password by userid|


# WS
## v1
|Method|URL          |Query Parameters|Behavior             |
|------|-------------|----------------|---------------------|
|GET   |`match`   |`%token`|Return the room info or an error|