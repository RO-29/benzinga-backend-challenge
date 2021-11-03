# benzinga-backend-challenge
Benzinga Backend Challenge, A webhhok receiver and forwarder.

### The application should be a basic webhook receiver that has two endpoints.
1. GET /healthz - should return HTTP 200-OK and “OK” as a string in the body
2. POST /log - should accept JSON payload

```
{
    "user_id": 1,
    "total": 1.65,
    "title": "delectus aut autem",
    "meta": {
        "logins": [{
            "time": "2020-08-08T01:52:50Z",
            "ip": "0.0.0.0"
        }],
        "phone_numbers": {
            "home": "555-1212",
            "mobile": "123-5555"
        }
    },
    "completed": false
}
```

### Requirements 

- The application should have three configurable values that should be read from environment variables: 
   - batch size
   - batch interval  
   - post endpoint
- The application should deserialize the JSON payload received at /log endpoint into a struct.
- Retain them all in-memory.
- When batch size is reached OR the batch interval has passed forward the collected records as an array to the post endpoint.
- Clear the in-memory cache of objects.
- Logger 
   - log an initialization message on startup 
   - On each HTTP request 
   - Each time it sends a batch  
     - log the batch size.
     - result status code.
     - duration of the POST request to the external endpoint. 
#### Addtional Requirements 
- If the POST fails 
   - Retry 3 times, waiting 2 seconds before each retry. 
   - After 3 failures log this failure and exit. 
- Testing for the Post endpoint output will be done against a service such as http://requestbin.net.

### Algorithm / Implementation

### Scope of Improvements

## github-action

- On every push to github, it runs linter / test and docker build
- #TODO upload the docker build to GCR (but it requires a paid gcr account)

## Build

- `docker-compose up -d` will run the http path(s) on localhost:8080
- `make build` puts the binary executable in `$root/build` folder.
- it's good to run `make lint test` before the make build command to ensure lint and test passes.
- make `docker-build` generates the docker image in `$root/build` folder.

## Lint
- `make lint` runs the golang-ci lint (installs it if not present) and runs the linter as defined in `root/.golangci.yml`.

## Test
- `make test` runs all the test in main as well as helper packages and generates the coverage report.

## cmd
- `main` http service is inside the `cmd/benzinga-backend-challenge`

## Postman Code for easy accesibilty
