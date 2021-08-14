# Notations
+ `%p`: `p` is an optional parameter.

# API
## v1
|Created version|Last update|Method|URL          |Query Parameters|Behavior             |
|----|------|------|-------------|----------------|---------------------|
|Alpha.0.0.1|Alpha.0.0.4|POST|`register`   |`%role, username, password`|Register a user with a specific role|
|Alpha.0.0.1|Alpha.0.0.4|POST|`auth`|`username, password`|Get an authenticated cookie|
|Alpha.0.0.1|Alpha.0.0.4|GET|`profile`||Get the profile|
|Alpha.0.0.4|Alpha.0.0.4|GET|`users/:id`||Get user's profile by userid|
|Alpha.0.0.4|Alpha.0.0.4|PUT|`users/:id`|`%age, %gender`|Update user's age, gender by userid|
|Alpha.0.0.4|Alpha.0.0.4|PUT|`users/:id/role`|`role`|Update user's role by userid|
|Alpha.0.0.4|Alpha.0.0.4|PUT|`users/:id/password`|`%oldpassword, newpassword`|Update user's password by userid|


# WS
## v1
|Created version|Last update|URL|Query Parameters|Behavior|
|----|----|-------------|----------------|---------------------|
|Alpha.0.0.3|Alpha.0.0.4|`match`|`%xytok`|Return the room info or an error.|
|Alpha.0.0.5|Alpha.0.0.5|`chat`|`%xytok`|A channel for exchanging chat messages.|