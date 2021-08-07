# Change logs
# Alpha.0.0.4
+ Modify general machanism in services.
+ Move controllers/debug.go to helpers/context/context.go.
+ Add some GET and PUT APIs in `/api/v1/users/:id/*`.
+ Add FEATURES.md for listing all APIs and WSes in the application.
+ Add CHECKLIST.md for listing all must-do tasks before committing.
+ Add command line options to build the application easily.
+ Add admin user seeder in command line.
+ Add many comments.
+ Change the loading environment variable method.
+ Use a correct exiting way by using `fatal` or `panic`.
+ Now it is possible to deploy on Heroku.
+ Add vendor.
+ Fix bugs of `Alpha.0.0.3`.

# Alpha.0.0.3
+ Add frontend with VueJS
+ Add matching feature with websocket at `/ws/v1/match`

# Alpha.0.0.2
+ Add docker container.

# Alpha.0.0.1
+ Add three APIs: `/api/v1/[auth|register|profile]`.
