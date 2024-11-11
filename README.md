# Weather Cloud Run

### Running version

There is a running version of this code on this link: https://gopos-cloudrun-747408531609.us-central1.run.app
Simply add a `zipcode` parameter with a zipcode (Brazilian CEP), like so: `?zipcode=01001000`

### Running tests

To run the tests locally, you can either run with Docker Compose:

```sh
docker compose -f docker-compose-test.yaml up --build
```

Or use the native command:

```sh
WEATHER_API_KEY=testkey go test ./...
```