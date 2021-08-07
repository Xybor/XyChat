# Notations
+ `%p`: `p` is an optional parameter.

# API
## v1
|Created version|Last update|Method|URL          |Query Parameters|Behavior             |
|----|------|------|-------------|----------------|---------------------|
|Alpha.0.0.1|Alpha.0.0.4|GET|`register`   |`%token, %role, username, password`|Register a user with a specific role|
|Alpha.0.0.1|Alpha.0.0.4|GET|`auth`|`username, password`|Get an authenticated token|
|Alpha.0.0.1|Alpha.0.0.4|GET|`profile`|`%token`|Get the profile|
|Alpha.0.0.4|Alpha.0.0.4|GET|`users/:id`|`%token`|Get user's profile by userid|
|Alpha.0.0.4|Alpha.0.0.4|PUT|`users/:id`|`%token, %age, %gender`|Update user's age, gender by userid|
|Alpha.0.0.4|Alpha.0.0.4|PUT|`users/:id/role`|`%token, role`|Update user's role by userid|
|Alpha.0.0.4|Alpha.0.0.4|PUT|`users/:id/password`|`%token, %oldpassword, newpassword`|Update user's password by userid|


# WS
## v1
|App version|Last update|Method|URL|Query Parameters|Behavior|
|----|----|--|-------------|----------------|---------------------|
|Alpha.0.0.3|Alpha.0.0.4|GET|`match`|`%token`|Return the room info or an error|