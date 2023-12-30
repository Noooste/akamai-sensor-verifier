# Akamai Sensor Verifier

This project aims verify your sensor data sent to akamai.

## BUILD AND RUN
1. Install Golang on https://go.dev/doc/install
2. Clone the project : `git clone git@github.com:Noooste/akamai-sensor-verifier.git`
3. Install the dependencies : `go mod download`
4. Build the project by doing `go build -o akamai .`
5. Run with `./akamai`

## USE
Paste your raw request payload, e.g. `{"sensor_data": "..."}` 
