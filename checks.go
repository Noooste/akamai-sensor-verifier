package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"github.com/Noooste/go-utils"
	"github.com/fatih/color"
	"math"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func Ab(t []byte) int {
	var a = 0
	for e := 0; e < len(t); e++ {
		n := t[e]
		if n < 128 {
			a += int(n)
		}
	}
	return a
}

var AllChecks = []any{
	sensorValue,
	KeyOrder,
	"-100",
	[]any{
		deviceDataLength,
		userAgentHash,
		timestampDivided,
		randomCalculation,
		z1,
		"information",
		[]any{
			callPhantom,
			activeXObject,
			documentMode,
			webstore,
			onLine,
			opera,
			installTrigger,
			HTMLElement,
			RTCPeerConnection,
			mozInnerScreenY,
			vibrate,
			getBattery,
			forEach,
			FileReader,
		},
	},
	"-115",
	[]any{
		keVel,
		meVel,
		teVel,
		doeVel,
		dmeVel,
		totVel,
		d2,
		keCnt,
		meCnt,
		d2Divided,
		teCnt,
		peCnt,
		ta,
		pizteIndex,
		abckHash,
		u1U2,
		webdriverCheck,
	},
	"-124",
	[]any{
		powLength,
		powCheck,
	},
	"-70",
	[]any{
		FpValStrLength,
	},
	"-80",
	[]any{
		fpValStrCalculated,
	},
}

func Check(om utils.OrderedMap) bool {
	buf, ok := check(om, AllChecks, false, "")

	if !ok {
		fmt.Println("┌───────────────────────┐")
		fmt.Print("│ CHECK RESULT : ")
		fmt.Println(color.RedString("FAILED") + " │")
		fmt.Println("├───────────────────────┘")
	} else {
		fmt.Println("┌───────────────────────┐")
		fmt.Print("│   CHECK RESULT : ")
		fmt.Println(color.GreenString("OK") + "   │")
		fmt.Println("├───────────────────────┘")
	}

	fmt.Println(buf.String())

	return ok
}

func GetFunctionName(temp interface{}) string {
	strs := strings.Split(runtime.FuncForPC(reflect.ValueOf(temp).Pointer()).Name(), ".")
	return ToSnakeCase(strs[len(strs)-1])
}

func CheckAssert(information utils.OrderedMap, last bool, skip bool, tab string, fn func(utils.OrderedMap) (ok bool, expected, actual string)) (buf *bytes.Buffer, ok bool) {
	var expected, actual string
	if skip {
		buf = new(bytes.Buffer)
		buf.WriteString(tab)
		if !last {
			buf.WriteString("├─── ")
		} else {
			buf.WriteString("└─── ")
		}

		buf.WriteString(color.HiBlackString(GetFunctionName(fn)))
		buf.WriteString("\n")
		return
	}

	ok, expected, actual = fn(information)
	buf = new(bytes.Buffer)
	buf.WriteString(tab)
	if !last {
		buf.WriteString("├─── ")
	} else {
		buf.WriteString("└─── ")
	}

	if expected != actual {
		buf.WriteString(color.RedString("✗ ") + GetFunctionName(fn) + color.RedString(" (expected: %s, got: %s)", expected, actual))
	} else {
		buf.WriteString(color.GreenString("✓ ") + GetFunctionName(fn))
	}

	buf.WriteString("\n")
	return
}

func check(om utils.OrderedMap, list []any, skip bool, tab string) (buf *bytes.Buffer, ok bool) {
	ok = true
	buf = new(bytes.Buffer)
	listLength := len(list)
	for i := 0; i < listLength; i++ {
		k := list[i]

		switch k.(type) {
		case func(utils.OrderedMap) (ok bool, expected, actual string):
			b, g := CheckAssert(om, i == listLength-1, skip, tab, list[i].(func(utils.OrderedMap) (ok bool, expected, actual string)))
			ok = ok && g
			skip = skip || !g
			buf.WriteString(b.String())

		case string:
			var addTab string
			if i < listLength-2 {
				addTab = tab + "│  "
			} else {
				addTab = tab + "   "
			}

			buf.WriteString(tab)
			if i < listLength-2 {
				buf.WriteString("├──┬ ")
			} else {
				buf.WriteString("└──┬ ")
			}

			if skip {
				i++
				l := list[i].([]any)
				b, _ := check(om, l, skip, addTab)
				buf.WriteString(color.HiBlackString(k.(string)))
				buf.WriteString("\n")
				buf.WriteString(b.String())
				continue
			}

			sectionName := color.HiYellowString(k.(string))
			i++

			l := list[i].([]any)

			if b, g := check(om, l, skip, addTab); !g {
				ok = false
				buf.WriteString(color.RedString("✗ ") + sectionName)
				buf.WriteString("\n")
				buf.WriteString(b.String())

			} else {
				buf.WriteString(color.GreenString("✓ ") + sectionName)
				buf.WriteString("\n")
				buf.WriteString(b.String())
			}

			if i < listLength-2 {
				buf.WriteString(tab)
				buf.WriteString("│\n")
			}
		}
	}

	return
}

func sensorValue(information utils.OrderedMap) (ok bool, expected, actual string) {
	t, ok := information.Map["raw"].([]uint8)

	if !ok {
		return false, "raw key found", "raw key not found"
	}

	sv := string(t)

	separator := strings.Split(sv[1:], ",2,")[0] + ","

	split := strings.Split(sv, separator)

	v := split[3]

	r := strconv.Itoa(24 ^ Ab([]byte(strings.Join(split[4:], separator))))

	if v != r {
		return false, v, r
	}

	return true, "", ""
}

var keyOrder = []string{"-100", "-105", "-108", "-101", "-110", "-117", "-109", "-102", "-111", "-114", "-103", "-106", "-115", "-112", "-119", "-122", "-123", "-124", "-126", "-127", "-128", "-131", "-132", "-133", "-70", "-80", "-90", "-116", "-129"}

func KeyOrder(information utils.OrderedMap) (ok bool, expected, actual string) {
	keys := information.Order[4:]

	for i, v := range keyOrder {
		if i > len(keys) {
			return false, v, ""
		}
		if v != keys[i] {
			return false, v, keys[i]
		}
	}
	return true, "", ""
}

func getSplitDeviceData(information utils.OrderedMap) (ok bool, expected, actual string, split []string) {
	dd, ok := information.Map["-100"]

	if !ok {
		return false, "DeviceData", "", nil
	}

	split1 := strings.Split(dd.(string), ",uaend,")
	split = strings.Split(split1[1], ",")

	s := make([]string, len(split)+2)

	s[0] = split1[0]
	s[1] = "uaend"

	for i, v := range split {
		s[i+2] = v
	}

	return true, "", "", s
}

func deviceDataLength(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	ok, expected, actual, split = getSplitDeviceData(information)

	if !ok {
		return
	}

	if len(split) != 39 && len(split) != 40 {
		return false, "39 or 40", strconv.Itoa(len(split))
	}

	return true, "", ""
}

func userAgentHash(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	ok, expected, actual, split = getSplitDeviceData(information)

	if !ok {
		return
	}

	ua := split[0]

	hash, err := strconv.Atoi(split[34])

	if err != nil {
		return false, "hash", err.Error()
	}

	value := Ab([]byte(ua))
	if value != hash {
		return false, strconv.Itoa(value), strconv.Itoa(hash)
	}

	return true, "", ""
}

func getStartTs(information utils.OrderedMap) (ok bool, expected, actual string, ts int) {
	var startTs any
	startTs, ok = information.Map["-115"]

	if !ok {
		return false, "StartTs", "", 0
	}

	startTs = strings.Split(startTs.(string), ",")[9]
	ts, err := strconv.Atoi(startTs.(string))

	if err != nil {
		return false, "StartTs", err.Error(), 0
	}

	return true, "", "", ts
}

func timestampDivided(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	ok, expected, actual, split = getSplitDeviceData(information)

	if !ok {
		return
	}

	ok, expected, actual, ts := getStartTs(information)

	if !ok {
		return
	}

	formatted := strconv.FormatFloat(float64(ts)/2, 'f', 1, 64)

	if strings.HasSuffix(formatted, ".0") {
		formatted = formatted[:len(formatted)-2]
	}

	if formatted != split[36] {
		return false, formatted, split[36]
	}

	return true, "", ""
}

func randomCalculation(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	ok, expected, actual, split = getSplitDeviceData(information)

	if !ok {
		return
	}

	value, err := strconv.ParseFloat(split[35][:11], 64)

	if err != nil {
		return false, "value", err.Error()
	}

	value2, err := strconv.Atoi(split[35][11:])
	if err != nil {
		return false, "value", err.Error()
	}

	if value2 != int(1e3*value/2) {
		return false, strconv.Itoa(int(1e3 * value / 2)), strconv.Itoa(value2)
	}

	return true, "", ""
}

func z1(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	ok, expected, actual, split = getSplitDeviceData(information)

	if !ok {
		return
	}

	var startTime int
	ok, expected, actual, startTime = getStartTs(information)

	z1 := int(math.Floor(float64(startTime / 4064256)))

	z1dd, err := strconv.Atoi(split[10])

	if err != nil {
		return false, "value", err.Error()
	}

	if z1 != z1dd {
		return false, strconv.Itoa(z1), strconv.Itoa(z1dd)
	}

	return true, "", ""
}

func callPhantom(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	ok, expected, actual, split = getSplitDeviceData(information)

	if !ok {
		return
	}

	if split[20] != "cpen:0" {
		return false, "cpen:0", split[20]
	}

	return true, "", ""
}

func activeXObject(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	ok, expected, actual, split = getSplitDeviceData(information)

	if !ok {
		return
	}

	name := strings.Split(split[21], ":")[0]
	if name != "i1" {
		return false, "i1", name
	}

	return true, "", ""
}

func documentMode(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	ok, expected, actual, split = getSplitDeviceData(information)

	if !ok {
		return
	}

	name := strings.Split(split[22], ":")[0]
	if name != "dm" {
		return false, "dm", name
	}

	return true, "", ""
}

func webstore(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	ok, expected, actual, split = getSplitDeviceData(information)

	if !ok {
		return
	}

	name := strings.Split(split[23], ":")[0]
	if name != "cwen" {
		return false, "cwen", name
	}

	return true, "", ""
}

func onLine(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	ok, expected, actual, split = getSplitDeviceData(information)

	if !ok {
		return
	}

	name := strings.Split(split[24], ":")[0]
	if name != "non" {
		return false, "non", name
	}

	return true, "", ""
}

func opera(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	ok, expected, actual, split = getSplitDeviceData(information)

	if !ok {
		return
	}

	name := strings.Split(split[25], ":")[0]
	if name != "opc" {
		return false, "opc", name
	}

	return true, "", ""
}

func installTrigger(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	ok, expected, actual, split = getSplitDeviceData(information)

	if !ok {
		return
	}

	name := strings.Split(split[26], ":")[0]
	if name != "fc" {
		return false, "fc", name
	}

	return true, "", ""
}

func HTMLElement(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	ok, expected, actual, split = getSplitDeviceData(information)

	if !ok {
		return
	}

	name := strings.Split(split[27], ":")[0]
	if name != "sc" {
		return false, "sc", name
	}

	return true, "", ""
}

func RTCPeerConnection(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	ok, expected, actual, split = getSplitDeviceData(information)

	if !ok {
		return
	}

	name := strings.Split(split[28], ":")[0]
	if name != "wrc" {
		return false, "wrc", name
	}

	return true, "", ""
}

func mozInnerScreenY(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	ok, expected, actual, split = getSplitDeviceData(information)

	if !ok {
		return
	}

	name := strings.Split(split[29], ":")[0]
	if name != "isc" {
		return false, "isc", name
	}

	return true, "", ""
}

func vibrate(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	ok, expected, actual, split = getSplitDeviceData(information)

	if !ok {
		return
	}

	name := strings.Split(split[30], ":")[0]
	if name != "vib" {
		return false, "vib", name
	}

	return true, "", ""
}

func getBattery(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	ok, expected, actual, split = getSplitDeviceData(information)

	if !ok {
		return
	}

	name := strings.Split(split[31], ":")[0]
	if name != "bat" {
		return false, "bat", name
	}

	return true, "", ""
}

func forEach(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	ok, expected, actual, split = getSplitDeviceData(information)

	if !ok {
		return
	}

	name := strings.Split(split[32], ":")[0]
	if name != "x11" {
		return false, "x11", name
	}

	return true, "", ""
}

func FileReader(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	ok, expected, actual, split = getSplitDeviceData(information)

	if !ok {
		return
	}

	name := strings.Split(split[33], ":")[0]
	if name != "x12" {
		return false, "x12", name
	}

	return true, "", ""
}

func FpValStrLength(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	inf, ok := information.Map["-70"]
	if !ok {
		return false, "true", "false"
	}

	split = strings.Split(inf.(string), ";")

	if len(split) != 14 {
		return false, "14", strconv.Itoa(len(split))
	}

	return true, "", ""
}

func fpValStrCalculated(information utils.OrderedMap) (ok bool, expected, actual string) {
	inf, ok := information.Map["-70"]
	if !ok {
		return false, "true", "false"
	}

	hash := Ab([]byte(inf.(string)))

	hashFound, ok := information.Map["-80"]

	if !ok {
		return false, "true", "false"
	}

	val, err := strconv.Atoi(hashFound.(string))

	if err != nil {
		return false, "true", "false"
	}

	if hash != val {
		return false, strconv.Itoa(hash), strconv.Itoa(val)
	}

	return true, "", ""
}

func getAbck(information utils.OrderedMap) (ok bool, expected, actual, abck string) {
	inf, ok := information.Map["-115"]
	if !ok {
		return false, "true", "false", ""
	}

	split := strings.Split(inf.(string), ",")

	if len(split) < 20 {
		return false, "length > 20", strconv.Itoa(len(split)), ""
	}

	abck = split[20]

	return
}

func splitPizte(information utils.OrderedMap) (ok bool, expected, actual string, split []string) {
	inf, ok := information.Map["-115"]
	if !ok {
		return false, "true", "false", nil
	}

	split = strings.Split(inf.(string), ",")

	return true, "", "", split
}

func abckHash(information utils.OrderedMap) (ok bool, expected, actual string) {
	ok, expected, actual, abck := getAbck(information)

	if !ok {
		return
	}

	hash := Ab([]byte(abck))

	inf, ok := information.Map["-115"]
	if !ok {
		return false, "true", "false"
	}

	split := strings.Split(inf.(string), ",")

	if len(split) < 20 {
		return false, "length > 20", strconv.Itoa(len(split))
	}

	hashFound := split[21]

	if !ok {
		return false, "true", "false"
	}

	val, err := strconv.Atoi(hashFound)

	if err != nil {
		return false, "true", "false"
	}

	if hash != val {
		return false, strconv.Itoa(hash), strconv.Itoa(val)
	}

	return true, "", ""
}

func pizteIndex(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	ok, expected, actual, split = splitPizte(information)

	if !ok {
		return
	}

	if split[25] != "PiZtE" {
		return false, "PiZtE", split[25]
	}

	return true, "", ""
}

func CalDis(list []float64) int {
	a := list[0] - list[1]
	e := list[2] - list[3]
	n := list[4] - list[5]
	return int(math.Sqrt(a*a + e*e + n*n))
}

func JrsReversed(rd, t int64) (o int) {
	var (
		a     = float64(rd)
		b     = float64(t)
		e     = strconv.FormatFloat(a*b, 'f', -1, 64)
		n     = 0
		oList []float64
		m     = len(e) >= 18
	)

	for len(oList) < 6 {
		value, err := strconv.Atoi(e[n : n+2])
		if err != nil {
			panic(err)
		}
		oList = append(oList, float64(value))
		if m {
			n += 3
		} else {
			n += 2
		}
	}
	return CalDis(oList)
}

func u1U2(information utils.OrderedMap) (ok bool, expected, actual string) {
	var startTime int
	ok, expected, actual, startTime = getStartTs(information)

	var split []string
	ok, expected, actual, split = splitPizte(information)

	if !ok {
		return
	}

	u1Found, err := strconv.ParseInt(split[26], 10, 64)

	if err != nil {
		return false, "true", "false"
	}

	var u2 = JrsReversed(u1Found, int64(startTime))

	u2Found := split[27]

	if u2Found != strconv.Itoa(u2) {
		return false, strconv.Itoa(u2), u2Found
	}

	return true, "", ""
}

func splitMouseData(information utils.OrderedMap) (ok bool, expected, actual string, split []string) {
	inf, ok := information.Map["-110"]
	if !ok {
		return false, "true", "false", nil
	}

	split = strings.Split(inf.(string), ";")

	if split[len(split)-1] != "" {
		return false, "last element is ;", split[len(split)-1], nil
	}

	return true, "", "", split[:len(split)-1]
}

func meVel(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	ok, expected, actual, split = splitMouseData(information)

	if !ok {
		return
	}

	var vel int

	for _, v := range split {
		s := strings.Split(v, ",")
		for _, g := range s[:5] {
			j, _ := strconv.Atoi(g)
			vel += j
		}
	}

	ok, expected, actual, split = splitPizte(information)

	if !ok {
		return
	}

	v, _ := strconv.Atoi(split[1])

	if vel+32 != v {
		return false, strconv.Itoa(vel + 32), split[1]
	}

	return true, "", ""
}

func meCnt(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	ok, expected, actual, split = splitMouseData(information)
	if len(split) == 0 {
		return true, "", ""
	}

	if !ok {
		return
	}

	var moveCnt int
	var clickCnt int
	var previousValue = -1

	for _, v := range split {
		s := strings.Split(v, ",")
		if s[1] == "1" {
			moveCnt++
		} else {
			clickCnt++
		}
		now, _ := strconv.Atoi(s[0])
		if now < 100 && now != previousValue+1 {
			return false, "now == previousValue + 1", fmt.Sprintf("now: %d, previousValue: %d", now, previousValue)
		}
		previousValue = now
	}

	if moveCnt > 100 {
		return false, "movement counter (0 -> 99)", strconv.Itoa(moveCnt)
	}

	if clickCnt > 75 {
		return false, "click counter (0 -> 74)", strconv.Itoa(clickCnt)
	}

	cnt, _ := strconv.Atoi(strings.Split(split[len(split)-1], ",")[0])

	ok, expected, actual, split = splitPizte(information)

	if !ok {
		return
	}

	if (moveCnt < 100 || clickCnt < 75) && split[13] != strconv.Itoa(cnt+1) {
		return false, strconv.Itoa(cnt + 1), split[13]
	}

	return true, "", ""
}

func splitKeyboardData(information utils.OrderedMap) (ok bool, expected, actual string, split []string) {
	inf, ok := information.Map["-108"]
	if !ok {
		return false, "true", "false", nil
	}

	split = strings.Split(inf.(string), ";")

	if split[len(split)-1] != "" {
		return false, "last element is ;", split[len(split)-1], nil
	}

	return true, "", "", split[:len(split)-1]
}

func keVel(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	ok, expected, actual, split = splitKeyboardData(information)

	if !ok {
		return
	}

	var vel int

	for _, v := range split {
		s := strings.Split(v, ",")
		for _, g := range s {
			j, _ := strconv.Atoi(g)
			vel += j
		}
	}

	ok, expected, actual, split = splitPizte(information)

	if !ok {
		return
	}

	v, _ := strconv.Atoi(split[0])

	if vel+1 != v {
		return false, strconv.Itoa(vel + 1), split[1]
	}

	return true, "", ""
}

func keCnt(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	ok, expected, actual, split = splitKeyboardData(information)
	if len(split) == 0 {
		return true, "", ""
	}

	if !ok {
		return
	}

	var keCntVal int
	var previousValue = -1

	for _, v := range split {
		keCntVal++
		s := strings.Split(v, ",")
		now, _ := strconv.Atoi(s[0])
		if now < 150 && now != previousValue+1 {
			return false, "key events are in order", "not in order"
		}
		previousValue = now
	}

	if keCntVal > 150 {
		return false, "keyboard counter < 150", strconv.Itoa(keCntVal)
	}

	cnt, _ := strconv.Atoi(strings.Split(split[len(split)-1], ",")[0])

	ok, expected, actual, split = splitPizte(information)

	if !ok {
		return
	}

	if cnt >= 150 && split[13] != strconv.Itoa(cnt+1) {
		return false, strconv.Itoa(cnt + 1), split[13]
	}

	return true, "", ""
}

func splitTouchData(information utils.OrderedMap) (ok bool, expected, actual string, split []string) {
	inf, ok := information.Map["-117"]
	if !ok {
		return false, "true", "false", nil
	}

	split = strings.Split(inf.(string), ";")

	if split[len(split)-1] != "" {
		return false, "last element is ;", split[len(split)-1], nil
	}

	return true, "", "", split[:len(split)-1]
}

func teVel(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	ok, expected, actual, split = splitTouchData(information)

	if !ok {
		return
	}

	var vel int

	for _, v := range split {
		s := strings.Split(v, ",")
		for _, g := range s {
			j, _ := strconv.Atoi(g)
			vel += j
		}
	}

	ok, expected, actual, split = splitPizte(information)

	if !ok {
		return
	}

	v, _ := strconv.Atoi(split[2])

	if vel+32 != v {
		return false, strconv.Itoa(vel + 32), split[2]
	}

	return true, "", ""
}

func teCnt(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	ok, expected, actual, split = splitTouchData(information)
	if len(split) == 0 {
		return true, "", ""
	}

	if !ok {
		return
	}

	var TeCntVal int
	var previousValue = -1

	for _, v := range split {
		TeCntVal++
		s := strings.Split(v, ",")
		now, _ := strconv.Atoi(s[0])
		if now < 25 && now != previousValue+1 {
			return false, "key events are in order", "not in order"
		}
		previousValue = now
	}

	if TeCntVal > 150 {
		return false, "touch counter < 150", strconv.Itoa(TeCntVal)
	}

	cnt, _ := strconv.Atoi(strings.Split(split[len(split)-1], ",")[0])

	ok, expected, actual, split = splitPizte(information)

	if !ok {
		return
	}

	if cnt >= 25 && split[15] != strconv.Itoa(cnt+1) && split[16] != strconv.Itoa(cnt+1) {
		return false, strconv.Itoa(cnt + 1), split[15]
	}

	return true, "", ""
}

func doeVel(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	info := information.Map["-109"].(string)

	split = strings.Split(info, ";")
	split = split[:len(split)-1]

	var vel int

	for _, v := range split {
		s := strings.Split(v, ",")
		j, _ := strconv.Atoi(s[1])
		vel += j
		j, _ = strconv.Atoi(s[0])
		vel += j
	}

	ok, expected, actual, split = splitPizte(information)

	if !ok {
		return
	}

	if split[4] != strconv.Itoa(vel) {
		return false, strconv.Itoa(vel), split[4]
	}

	return true, "", ""
}

func dmeVel(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	info := information.Map["-111"].(string)
	split = strings.Split(info, ";")
	split = split[:len(split)-1]

	var vel int

	for _, v := range split {
		s := strings.Split(v, ",")
		j, _ := strconv.Atoi(s[1])
		vel += j
		j, _ = strconv.Atoi(s[0])
		vel += j
	}

	ok, expected, actual, split = splitPizte(information)

	if !ok {
		return
	}

	if split[3] != strconv.Itoa(vel) {
		return false, strconv.Itoa(vel), split[3]
	}

	return true, "", ""
}

func totVel(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	info := information.Map["-115"].(string)

	split = strings.Split(info, ",")[:6]

	var vel int
	v, _ := strconv.Atoi(split[0])
	vel += v - 1
	v, _ = strconv.Atoi(split[1])
	vel += v - 32
	v, _ = strconv.Atoi(split[2])
	vel += v - 32

	for _, f := range split[3:] {
		v, _ = strconv.Atoi(f)
		vel += v
	}

	ok, expected, actual, split = splitPizte(information)

	if !ok {
		return
	}

	v, _ = strconv.Atoi(split[6])

	if vel != v {
		return false, strconv.Itoa(vel), split[6]
	}

	return true, "", ""
}

func d2(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	ok, expected, actual, split = splitPizte(information)

	if !ok {
		return
	}

	var startTime int
	ok, expected, actual, startTime = getStartTs(information)

	z := int(math.Floor(float64(startTime / 4064256)))

	if split[11] != strconv.Itoa(z/23) {
		return false, strconv.Itoa(z / 23), split[11]
	}

	return true, "", ""
}

func d2Divided(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	ok, expected, actual, split = splitPizte(information)

	if !ok {
		return
	}

	var startTime int
	ok, expected, actual, startTime = getStartTs(information)

	z := int(math.Floor(float64(startTime / 4064256)))

	if split[14] != strconv.Itoa((z/23)/6) {
		return false, strconv.Itoa((z / 23) / 6), split[14]
	}

	return true, "", ""
}

func splitPointerData(information utils.OrderedMap) (ok bool, expected, actual string, split []string) {
	info := information.Map["-114"].(string)

	split = strings.Split(info, ";")
	split = split[:len(split)-1]

	if len(split) == 0 {
		return false, "", "", split
	}

	return true, "", "", split
}

func peCnt(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	ok, expected, actual, split = splitPointerData(information)
	if !ok {
		ok, expected, actual, split = splitMouseData(information)
		if !ok {
			return
		}
	}

	var cnt int

	for _, v := range split {
		s := strings.Split(v, ",")
		if s[1] == "4" || s[1] == "3" {
			cnt++
		}
	}

	ok, expected, actual, split = splitPizte(information)
	if !ok {
		return
	}

	if cnt < 51 && split[16] != strconv.Itoa(cnt) && split[15] != strconv.Itoa(cnt) {
		return false, strconv.Itoa(cnt), split[16]
	}

	return true, "", ""
}

func ta(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split2 []string
	ok, expected, actual, split2 = splitMouseData(information)

	if !ok {
		return
	}

	var split3 []string
	ok, expected, actual, split3 = splitKeyboardData(information)
	if !ok {
		return
	}

	var split4 []string
	ok, expected, actual, split4 = splitTouchData(information)
	if !ok {
		return
	}

	var split5 []string
	ok, expected, actual, split5 = splitPointerData(information)

	var t int

	for _, v := range split2 {
		s := strings.Split(v, ",")
		f, _ := strconv.Atoi(s[2])
		t += f
	}

	for _, v := range split3 {
		s := strings.Split(v, ",")
		f, _ := strconv.Atoi(s[2])
		t += f
	}

	for _, v := range split4 {
		s := strings.Split(v, ",")
		f, _ := strconv.Atoi(s[2])
		t += f
	}

	for _, v := range split5 {
		s := strings.Split(v, ",")
		f, _ := strconv.Atoi(s[2])
		t += f
	}

	info := information.Map["-111"].(string)
	split := strings.Split(info, ";")
	split = split[:len(split)-1]

	var dmadoa int
	for _, v := range split {
		s := strings.Split(v, ",")
		j, _ := strconv.Atoi(s[1])
		t += j
		dmadoa += j
	}

	info = information.Map["-109"].(string)
	split = strings.Split(info, ";")
	split = split[:len(split)-1]

	for _, v := range split {
		s := strings.Split(v, ",")
		j, _ := strconv.Atoi(s[1])
		t += j
		dmadoa += j
	}

	ok, expected, actual, split = splitPizte(information)

	if !ok {
		return
	}

	if split[18] != strconv.Itoa(t) && split[18] != strconv.Itoa(t-dmadoa) {
		return false, strconv.Itoa(t), split[18]
	}

	return true, "", ""
}

func webdriverCheck(information utils.OrderedMap) (ok bool, expected, actual string) {
	var split []string
	ok, expected, actual, split = splitPizte(information)

	if !ok {
		return
	}

	if split[28] == "0;0" {
		if split[29] != "0" {
			return false, "0;0", split[29]
		}

		if split[30] != "0" {
			return false, "0", split[30]
		}
	} else {
		if split[28] != "0" {
			return false, "0", split[28]
		}

		if split[29] != "0" {
			return false, "0", split[29]
		}

		if split[30] != "0" {
			return false, "0", split[30]
		}
	}

	return true, "", ""
}

func powLength(information utils.OrderedMap) (ok bool, expected, actual string) {
	challenge := information.Map["-124"].(string)
	if challenge == "" {
		return true, "", ""
	}

	split := strings.Split(challenge, ";")

	if len(split) != 5 {
		return false, "5", strconv.Itoa(len(split))
	}

	return true, "", ""
}

func bdm(t [32]byte, a int) int {
	e := int(t[0])
	for n := 1; n < 32; n++ {
		e = (e<<8 + int(t[n])) % a
	}
	return e
}

func powCheck(information utils.OrderedMap) (ok bool, expected, actual string) {
	challenge := information.Map["-124"].(string)
	if challenge == "" {
		return true, "", ""
	}

	split := strings.Split(challenge, ";")
	split2 := strings.Split(split[3], ",")

	for i, el := range strings.Split(split[0], ",") {
		if el == "" {
			return true, "", ""
		}
		val, _ := strconv.Atoi(split2[4])

		hash := sha256.Sum256([]byte(split2[3] + strconv.Itoa(val+i) + el))
		r := bdm(hash, val+i)
		if r != 0 {
			return false, "0", strconv.Itoa(bdm(hash, val+i))
		}
	}

	return true, "", ""
}
