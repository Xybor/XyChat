# Note
+ `%p`: `p` is an optional parameter.
+ All POST, PUT, DELETE method APIs get paramters from body under json format.
+ All GET method APIs get parameters from url query.
+ Web socket handshake always gets parameters from url query.
+ Web socket communication always uses data under json format.

# Response format
This is the common format for both API and WS responses
```py
{
    "data": <data>,
    "meta": {
        "errno": <error number>,
        "error": <error message>
    }
}
```
In the `meta` field, if `errno` equals to zero, this is a response of the successful request and there is no `error` field. Otherwise, please show the `error` field for the client.

The `data` field contains received parameters if the request is successful.


# API
## v1
|Created version|Last update|Method|URL|Sent parameters|Received parameters|Behavior|
|:----:|:------:|:------:|-------------|----------------|----|---------------------|
|Alpha.0.1|Alpha.0.6|POST|`register`|`username` `password` `%role`||Register a user with a specific role|
|Alpha.0.1|Alpha.0.6|POST|`auth`|`username` `password`|`id` `username` `role` `%age` `%gender`|Get an authenticated cookie|
|Alpha.0.1|Alpha.0.4|GET|`profile`||`id` `username` `role` `%age` `%gender`|Get the profile|
|Alpha.0.4|Alpha.0.4|GET|`users/:id`||`id` `username` `role` `%age` `%gender`|Get user's profile by userid|
|Alpha.0.4|Alpha.0.6|PUT|`users/:id`|`%age` `%gender`||Update user's age, gender by userid|
|Alpha.0.4|Alpha.0.6|PUT|`users/:id/role`|`role`||Update user's role by userid|
|Alpha.0.4|Alpha.0.6|PUT|`users/:id/password`|`new_password` `%old_password`||Update user's password by userid|


# WS
## v1
|Created version|Last update|URL|Handshake parameters|Sent parameters|Received parameter|Behavior|
|:----:|:----:|-------------|----------------|---------------------|---|---|
|Alpha.0.3|Alpha.0.4|`match`|`%xytok`||`roomid`|Request to push the current user to the matching queue and wait for a matching room|
|Alpha.0.5|Alpha.0.5|`chat`|`%xytok`|`roomid` `message`|`userid` `roomid` `message` `time`|A channel for sending and receiving messages from all rooms of the current user|
