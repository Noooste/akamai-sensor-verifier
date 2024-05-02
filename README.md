# Akamai Sensor Verifier

The purpose of this project is to verify the sensor data sent to Akamai.

## BUILD AND RUN
1. Install Golang on https://go.dev/doc/install (`^1.20`)
2. Clone the project : `git clone git@github.com:Noooste/akamai-sensor-verifier.git`
3. Install the dependencies : `go mod download`
4. Build the project by doing `go build -o akamai .`
5. Run with `./akamai`

## USE
1. Open your browser and go to a protected website, e.g. [nike](https://www.nike.com/).
2. Open the "dev" tab and look for requests that have a URL like `https://www.nike.com/Bpk0/YiEp/lV1db/mt8kw/7SEhmLNf1aaO1Y/YyZNWVcPAQ/LwITYH/BgN2U`, if you want to get the url easily, you can simply click anywhere on the page to see that it's sent.
3. Copy the request body (`{"sensor_data": "..."}`), paste it to the program: you'll see what Akamai sends to its servers.
