# benzinga-backend-challenge
Benzinga Backend Challenge, a simple webhook receiver and forwarder.

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
- Logger (`"github.com/sirupsen/logrus"`)
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

### Arguments
- Can be set through os env or overriding through flag args passed while running the build
- Flags
```
batch-interval duration
        Batch Interval (default 10s)
  -batch-size int
        Batch Size (default 10)
  -http string
        HTTP  (default ":8080")
  -log-file string
        log file
  -post-endpoint string
        Post Endpoint
```
OR
- Env
```
WEBHOOK_POST_ENDPOINT=<valid string url>
env=WEBHOOK_BATCH_SIZE=<valid int>
WEBHOOK_BATCH_INTERVAL=<valid golang time.Duration>
```

### Algorithm / Implementation
- A buffered channel is init at program startup, if `batch size > 0`, `buffered channel = 2*batch_size` or `100` by default.
- Buffered channel will make sure, /log is not blocking deserializing while our forwarder consumer is processing.
- At same, time HTTP handlers / receivers are regsitered. (Routing is handled via gorrila/mux).
- `/log` recieves the json, deserile in thedefined strcut and pushes to a buffered channel.
- We start a `forwarder consumer` in background and it iterates over buffer channel and accumulates msgs in a simple in-memory array.
- `forwarder consumer` itself regsiter a background process which starts a forever blocking select channel.
- select channel has three cases, either a batch_size is full or batch_interval is reached or accumualte the msg in in-memory array.
- When batch(size/interval) is reached, it forwards the accumualted msgs to webhook endpoint as per the requirement
- if there is error in "post call" after three retries, it sends the error to err channel and then the programs exits as per the requirement



### Scope of Improvements
- Start `forwarder` consumer in waitgroup so it multiple consumers could be started in background in case of high throughput.
- Add more unit tests 🙈
- 
## github-action

- On every push to github, it runs linter / test and docker build
- #TODO upload the docker build to GCR (but it requires a paid gcr account)

## Build

- `make build` puts the binary executable in `$root/build` folder.
- it's good to run `make lint test` before the make build command to ensure lint and test passes.
- make `docker-build` generates the docker image in `$root/build` folder.

## Lint
- `make lint` runs the golang-ci lint (installs it if not present) and runs the linter as defined in `root/.golangci.yml`.

## Test
- `make test` runs all the test in main package and generates the coverage report.

## cmd
- `main` http service is inside the `cmd/benzinga-backend-challenge`

## Postman Code for easy accesibilty
- Log Payload
```
curl -X POST \
  http://localhost:8080/log \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{
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
}'
````
- Health Check
```
curl -X GET \
  http://localhost:8080/healthz \
  -H 'cache-control: no-cache'
```

