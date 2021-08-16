# Change logs

# Alpha.0.5
## Alpha.0.5.1
- Add chat feature with websocket at `/ws/v1/chat`.
- Allow CORS by using environment variable.
- Split README.md into many files.

## Alpha.0.5.2
- Use TLS environment variable to dynamically enable HTTPS method.
- Use environment variable to set admin user.
- Use XYCHAT environment variable to determine that the application is running for what (test, debug or release).

## Alpha.0.5.3
- Return user information in auth API.
- Redesign the error management.
- Remove ApplyAPIHeader.


# Alpha.0.4
- Modify general machanism in services.
- Add some GET and PUT APIs in `/api/v1/users/:id/*`.
- Add FEATURES.md for listing all APIs and WSes in the application.
- Add CHECKLIST.md for listing all must-do tasks before committing.
- Add command line options to build the application easily.
- Add many comments.
- Change the loading environment variable method.
- Use the correct exiting way by using `fatal` or `panic`.
- Now it is possible to deploy on Heroku.
- Add vendor.

# Alpha.0.3

- Add matching feature with websocket at `/ws/v1/match`

# Alpha.0.2

- Add docker container.

# Alpha.0.1

- Add three APIs: `/api/v1/[auth|register|profile]`.
