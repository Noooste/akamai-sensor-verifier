package main

import (
	"encoding/json"
	"fmt"
	"github.com/Noooste/go-utils"
	"github.com/fatih/color"
	"github.com/mattn/go-tty"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var terminalWidth int

func init() {
	var err error
	terminalWidth, _, err = terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		log.Println("Error getting terminal size:", err)
		return
	}
}

func main() {
	fmt.Print("Enter payload: ")

	t, err := tty.Open()
	if err != nil {
		fmt.Println("Error opening terminal:", err)
		return
	}
	defer t.Close()

	payload, err := t.ReadString()
	if err != nil {
		fmt.Println("Error reading from terminal:", err)
		return
	}

	// Remove the newline character from the payload
	payload = strings.TrimSpace(payload)

	fmt.Println(strings.Repeat("─", terminalWidth))

	r := decryptMain(payload)

	sensorData := r.Map["raw"].([]uint8)

	separator := strings.Split(string(sensorData[1:]), ",2,")[0] + ","

	split := strings.Split(string(sensorData), separator)[2:]

	r.Map["sensor_data"] = split

	fmt.Println("sensor_data :")
	fmt.Println("[")
	for i := 0; i < len(split); i++ {
		fmt.Printf("  %s,\n", color.GreenString("'%s'", split[i]))
	}
	fmt.Println("]")

	fmt.Println(strings.Repeat("─", terminalWidth))

	if Check(r) {
		DisplayInformation(r)
	}
}

type sensorDataStruct struct {
	SensorData string `json:"sensor_data"`
}

func decryptMain(payload string) (result utils.OrderedMap) {
	var sensorData, prefix []byte
	var separator string
	var key1, key2 int

	var sensor sensorDataStruct
	// parse json
	if err := json.Unmarshal([]byte(payload), &sensor); err != nil {
		log.Fatal(err)
	}

	// extract keys
	keys := strings.Split(sensor.SensorData, ";")
	key1, _ = strconv.Atoi(keys[2])
	key2, _ = strconv.Atoi(keys[3])

	// Parse prefix
	sensorData = []byte(sensor.SensorData)
	encrypted := sensorData

	re := regexp.MustCompile(`(\d+;\d+;\d+;\d+;[\d,]+;)`)
	if re.Match(sensorData) {
		prefix = re.Find(sensorData)
		sensorData = sensorData[len(prefix):]
	}

	// Obfuscated strings to plaintext
	sensorData = []byte(decrypt(string(sensorData), uint32(key1)))

	sensorData = []byte(decryptInner(string(sensorData), uint32(key2)))
	raw := sensorData

	separator = strings.Split(string(sensorData[1:]), ",2,")[0] + ","

	split := strings.Split(string(sensorData), separator)[2:]

	result.Order = make([]string, 4+len(split[2:])/2)
	result.Map = make(map[string]any)

	result.Order[0] = "key"
	result.Map["key"] = split[0]

	result.Order[1] = "sensor_value"
	result.Map["sensor_value"] = split[1]

	result.Order[2] = "raw"
	result.Map["raw"] = raw

	result.Order[3] = "encrypted"
	result.Map["encrypted"] = encrypted

	orderIndex := 4
	split = split[2:]
	for i := 0; i < len(split); i += 2 {
		result.Order[orderIndex+i/2] = split[i]
		result.Map[split[i]] = split[i+1]
	}

	return result
}
