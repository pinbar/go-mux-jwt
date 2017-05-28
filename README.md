## A simple example of Go (golang) api with jwt authentication

### tech stack
* **Go** - a programming language that is fast, uses minimal resources and support high concurrency
* **Gorilla Mux** - minimalistic request router and dispatcher
* **Gorilla Logging Handler** - middleware for http request/response logging
* **jwt-go** - a jwt library for Go
* **testify** - testing and assertion library
* **gin (optional)** - livereload utility for faster development turnaround in local

### pre-requisites
* Go is installed. To verify, run `go version`
* GOPATH is set (e.g. set to `~/go`)
* PATH includes `GOPATH/bin`

### getting started
* clone repo or download zip
* get the dependencies
    * `go get github.com/gorilla/mux`
    * `go get github.com/gorilla/handlers`
    * `go get github.com/dgrijalva/jwt-go`
* in the project directory, run `go build && ./go-mux-jwt`
* launch the browser and point to the baseurl `localhost:3001` (port can be changed in `main.go`)
* *optional:*
    * use **gin** to monitor for changes and automatically restart the application
    * if you don't have gin, `go get github.com/codegangsta/gin`
    * in the project directory run `gin` (no need to build or run executable, when you do this)

### running tests
* to run the tests, run `go test` in the project directory
* **test coverage:** 
    * to run tests and generate coverage report, run `go test -cover`. 
    * percentage covered is shown in the terminal upon execution of this command

### api and authentication scenarios
* access the unsecure api `GET /metacortex`
* all `/api/*` calls are secured with JWT authentication
* try accessing the secure api `GET /api/megacity` to see an auth error
* obtain a JWT token here: `/static/authenticate.html`
    * enter programName:programPassword (neo:keanu)
    * the response contains a JWT token for that program
* use the token when calling any secure api (`/api/*`):
    * set the `Authorization` request header and add the jwt token, like so:
    * `Authorization: Bearer \<token\>`
* `GET /api/megacity` can be accessed only with a valid token
