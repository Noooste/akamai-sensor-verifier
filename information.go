package main

import (
	"bytes"
	"fmt"
	"github.com/Noooste/go-utils"
	"github.com/fatih/color"
	"github.com/leekchan/accounting"
	"math"
	"strconv"
	"strings"
	"time"
)

var AllInformation = []any{
	"sensor_format",
	[]any{
		sensorSeparator,
		sensorBuildTime,
		scriptInitTime,
		sensorIdTime,
		"shuffling",
		[]any{
			shufflingKey,
			shufflingTime,
		},
		"encryption",
		[]any{
			encryptionKey,
			encryptionTime,
		},
	},
	"sensor_content",
	[]any{
		url,
		"device",
		[]any{
			userAgent,
			rcfp,
			rValue,
			fpValStr,
			"screen",
			[]any{
				screenWidth,
				screenHeight,
				availableWidth,
				availableHeight,
				innerWidth,
				innerHeight,
				outerWidth,
			},
		},
		"-131",
		[]any{
			"js_heap",
			[]any{
				JSHeapLimit,
				JSHeapTotal,
				JSHeapUsed,
			},
		},
		"locale",
		[]any{
			timezoneOffset,
			lang,
			langLen,
			langHash,
		},
		"timestamps",
		[]any{
			startTimestampLocalTime,
			dtBpd,
		},
		"activity",
		[]any{
			"mact",
			[]any{
				"delta_time",
				[]any{
					dtMact,
					dtMactList,
				},
			},
			"tact",
			[]any{
				"delta_time",
				[]any{
					dtTact,
					dtTactList,
				},
				"delta_position",
				[]any{
					dPosTactList,
				},
				"delta_time_delta_position",
				[]any{
					dtdPosTactList,
				},
				"velocity",
				[]any{
					tactVelocity,
				},
				"ratio counter movement",
				[]any{
					dtdPosCntTactList,
				},
			},
		},
	},
}

func DisplayInformation(information utils.OrderedMap) {
	buf := display(information, AllInformation, "")

	fmt.Println("┌───────────────────────┐")
	fmt.Println("│      INFORMATION      │")
	fmt.Println("├───────────────────────┘")

	fmt.Println(buf.String())

}

func getTerminalSize() (width int, height int) {
	width = 100
	height = 100
	return
}
func displaySpecificInformation(information utils.OrderedMap, last bool, tab string, fn func(om utils.OrderedMap) (buf *bytes.Buffer)) (buf *bytes.Buffer) {
	// Call the provided function and get the output buffer

	buf = new(bytes.Buffer)

	// Add proper indentation and tree structure
	buf.WriteString(tab)
	if !last {
		buf.WriteString("├─── ")
	} else {
		buf.WriteString("└─── ")
	}

	// Add function name and set the color
	fnName := GetFunctionName(fn)

	addLen := len(fnName) + len(": ")
	avail := terminalWidth - len(tab) - addLen

	buf.WriteString(fnName + ": \u001B[36m")
	information.Map["x-add-indent"] = addLen
	information.Map["x-available-width"] = avail
	information.Map["x-tab"] = tab

	b := fn(information)

	buf.WriteString(b.String())
	buf.WriteString("\n")

	// Reset color and return buffer
	buf.WriteString("\u001B[0m")
	return
}

func display(om utils.OrderedMap, list []any, tab string) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	listLength := len(list)
	for i := 0; i < listLength; i++ {
		k := list[i]

		switch k.(type) {
		case func(utils.OrderedMap) (buf *bytes.Buffer):
			buf.WriteString(displaySpecificInformation(om, i == listLength-1, tab, k.(func(utils.OrderedMap) (buf *bytes.Buffer))).String())

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

			buf.WriteString(color.HiYellowString(k.(string)))
			buf.WriteString("\n")
			i++
			l := list[i].([]any)
			buf.WriteString(display(om, l, addTab).String())

			if i < listLength-2 {
				buf.WriteString(tab)
				buf.WriteString("│\n")
			}
		}
	}

	return
}

func url(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	buf.WriteString(information.Map["-112"].(string))
	return
}

func userAgent(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	_, _, _, split := getSplitDeviceData(information)
	buf.WriteString(split[0])
	return
}

func startTimestampLocalTime(information utils.OrderedMap) (buf *bytes.Buffer) {
	_, _, _, ts := getStartTs(information) //ts is in millisecond
	buf = new(bytes.Buffer)
	timeInThisZone := time.UnixMilli(int64(ts)).In(time.Local)
	buf.WriteString(timeInThisZone.Format(time.DateTime))
	return
}

func dtMact(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	_, _, _, split := splitMouseData(information) //ts is in millisecond

	if len(split) == 0 {
		buf.WriteString("no mouse movement")
		return
	}

	var averageDt int
	var nbDt int
	var minDt, maxDt int

	lastT := 0

	for i := 0; i < len(split); i++ {
		s := strings.Split(split[i], ",")
		t, _ := strconv.Atoi(s[2])

		dt := t - lastT
		lastT = t

		if dt < 30 && dt > 0 {
			averageDt += dt
			nbDt++
		}

		if minDt == 0 || (dt > 0 && dt < minDt) {
			minDt = dt
		}

		if dt < 30 && dt > maxDt {
			maxDt = dt
		}

	}

	if nbDt > 0 {
		averageDt = averageDt / nbDt
	}

	buf.WriteString(fmt.Sprintf("average: %d ms, min: %d ms, max: %d ms", averageDt, minDt, maxDt))
	return
}

func dtMactList(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	_, _, _, split := splitMouseData(information) //ts is in millisecond

	if len(split) == 0 {
		buf.WriteString("no mouse movement")
		return
	}

	lastT := 0

	w := information.Map["x-available-width"].(int)
	indent := information.Map["x-add-indent"].(int)
	tab := information.Map["x-tab"].(string)

	currentWidthLeft := w - indent - len(tab) - 8

	//inform about the color code
	buf.WriteString(color.WhiteString(""))
	buf.WriteString("what's type : move, \u001B[32mclick\u001B[0m, \033[36mdown\u001B[0m, \033[33mup\033[0m | ")
	for i := 0; i < len(split); i++ {
		s := strings.Split(split[i], ",")
		t, _ := strconv.Atoi(s[2])

		dt := t - lastT
		lastT = t

		var inf string
		switch s[1][0] {
		case '1':
			inf = fmt.Sprintf("%d, ", dt)
		case '2':
			inf = fmt.Sprintf("\u001B[32m%d\u001B[0m, ", dt)
		case '3':
			inf = fmt.Sprintf("\033[36m%d\033[0m, ", dt)
		case '4':
			inf = fmt.Sprintf("\033[33m%d\033[0m, ", dt)
		default:
			inf = fmt.Sprintf("%d, ", dt)
		}

		lengthWithoutColors := len(fmt.Sprintf("%d, ", dt))
		if lengthWithoutColors > currentWidthLeft {
			buf.WriteString("\n")
			buf.WriteString(tab)
			buf.WriteString(strings.Repeat(" ", indent))
			currentWidthLeft = w
		}

		buf.WriteString(inf)
	}

	buf.Truncate(buf.Len() - 2)

	return
}

func dtTact(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	_, _, _, split := splitTouchData(information) //ts is in millisecond

	if len(split) == 0 {
		buf.WriteString("no touch activity")
		return
	}

	var averageDt int
	var nbDt int
	var minDt, maxDt int

	lastT := 0

	for i := 0; i < len(split); i++ {
		s := strings.Split(split[i], ",")
		t, _ := strconv.Atoi(s[2])

		dt := t - lastT
		lastT = t

		if dt < 30 && dt > 0 {
			averageDt += dt
			nbDt++
		}

		if minDt == 0 || (dt > 0 && dt < minDt) {
			minDt = dt
		}

		if dt < 30 && dt > maxDt {
			maxDt = dt
		}

	}

	if nbDt > 0 {
		averageDt = averageDt / nbDt
	}

	buf.WriteString(fmt.Sprintf("average: %d ms, min: %d ms, max: %d ms", averageDt, minDt, maxDt))
	return
}

func dtTactList(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	_, _, _, split := splitTouchData(information) //ts is in millisecond

	if len(split) == 0 {
		buf.WriteString("no touch activity")
		return
	}

	lastT := 0

	w := information.Map["x-available-width"].(int)
	indent := information.Map["x-add-indent"].(int)
	tab := information.Map["x-tab"].(string)

	currentWidthLeft := w - indent - len(tab) - 8

	//inform about the color code
	buf.WriteString(color.WhiteString(""))
	buf.WriteString("what's type : move, \033[36mdown\u001B[0m, \033[33mup\033[0m | ")
	for i := 0; i < len(split); i++ {
		s := strings.Split(split[i], ",")
		t, _ := strconv.Atoi(s[2])

		dt := t - lastT
		lastT = t

		var inf string
		switch s[1][0] {
		case '2':
			inf = fmt.Sprintf("\u001B[36m%d\u001B[0m, ", dt)
		case '3':
			inf = fmt.Sprintf("\033[33m%d\033[0m, ", dt)
		default:
			inf = fmt.Sprintf("%d, ", dt)
		}

		lengthWithoutColors := len(fmt.Sprintf("%d, ", dt))
		if lengthWithoutColors > currentWidthLeft {
			buf.WriteString("\n")
			buf.WriteString(tab)
			buf.WriteString(strings.Repeat(" ", indent))
			currentWidthLeft = w
		}

		buf.WriteString(inf)
		currentWidthLeft -= lengthWithoutColors
	}

	buf.Truncate(buf.Len() - 2)
	return
}

func dPosTactList(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	_, _, _, split := splitTouchData(information) //ts is in millisecond

	if len(split) == 0 {
		buf.WriteString("no touch activity")
		return
	}

	var (
		lastX, lastY int
	)
	w := information.Map["x-available-width"].(int)
	indent := information.Map["x-add-indent"].(int)
	tab := information.Map["x-tab"].(string)

	currentWidthLeft := w - indent - len(tab) - 8

	//inform about the color code
	buf.WriteString(color.WhiteString(""))
	buf.WriteString("what's type : move, \033[36mdown\u001B[0m, \033[33mup\033[0m | ")
	for i := 0; i < len(split); i++ {
		s := strings.Split(split[i], ",")

		x, _ := strconv.Atoi(s[3])
		y, _ := strconv.Atoi(s[4])

		dX := x - lastX
		dY := y - lastY

		var inf string
		switch s[1][0] {
		case '2':
			inf = fmt.Sprintf("\033[36m(%d, %d)\033[0m, ", dX, dY)
		case '3':
			inf = fmt.Sprintf("\033[33m(%d, %d)\033[0m, ", dX, dY)
		default:
			inf = fmt.Sprintf("(%d, %d), ", dX, dY)
		}

		lengthWithoutColors := len(fmt.Sprintf("(%d, %d), ", dX, dY))
		if lengthWithoutColors > currentWidthLeft {
			buf.WriteString("\n")
			buf.WriteString(tab)
			buf.WriteString(strings.Repeat(" ", indent))
			currentWidthLeft = w
		}
		buf.WriteString(inf)
		currentWidthLeft -= lengthWithoutColors

		lastX = x
		lastY = y
	}

	buf.Truncate(buf.Len() - 2)
	return
}

func dtdPosTactList(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	_, _, _, split := splitTouchData(information) //ts is in millisecond

	if len(split) == 0 {
		buf.WriteString("no touch activity")
		return
	}

	var (
		lastX, lastY int
		lastT        int
	)

	w := information.Map["x-available-width"].(int)
	indent := information.Map["x-add-indent"].(int)
	tab := information.Map["x-tab"].(string)

	currentWidthLeft := w - indent - len(tab) - 4

	//inform about the color code
	buf.WriteString(color.WhiteString(""))
	buf.WriteString("what's type : move, \033[36mdown\u001B[0m, \033[33mup\033[0m | ")

	for i := 0; i < len(split); i++ {
		s := strings.Split(split[i], ",")

		x, _ := strconv.Atoi(s[3])
		y, _ := strconv.Atoi(s[4])
		t, _ := strconv.Atoi(s[2])

		dX := x - lastX
		dY := y - lastY
		dt := t - lastT

		var inf string
		switch s[1][0] {
		case '1':
			inf = fmt.Sprintf("(%d: %d, %d), ", dt, dX, dY)
		case '2':
			inf = fmt.Sprintf("\033[36m(%d: %d, %d)\033[0m, ", dt, dX, dY)
		case '3':
			inf = fmt.Sprintf("\033[33m(%d: %d, %d)\033[0m, ", dt, dX, dY)
		default:
			inf = fmt.Sprintf("(%d: %d, %d), ", dt, dX, dY)
		}

		lengthWithoutColors := len(fmt.Sprintf("(%d: %d, %d), ", dt, dX, dY))
		if lengthWithoutColors > currentWidthLeft {
			buf.WriteString("\n")
			buf.WriteString(tab)
			buf.WriteString(strings.Repeat(" ", indent))
			currentWidthLeft = w
		}
		buf.WriteString(inf)
		currentWidthLeft -= lengthWithoutColors

		lastX = x
		lastY = y
		lastT = t
	}

	buf.Truncate(buf.Len() - 2)
	return
}

func tactVelocity(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	_, _, _, split := splitTouchData(information) //ts is in millisecond

	if len(split) == 0 {
		buf.WriteString("no touch activity")
		return
	}

	var (
		lastX, lastY int
		lastT        int
	)

	w := information.Map["x-available-width"].(int)
	indent := information.Map["x-add-indent"].(int)
	tab := information.Map["x-tab"].(string)

	currentWidthLeft := w - indent - len(tab) - 8

	//inform about the color code
	buf.WriteString(color.WhiteString(""))
	buf.WriteString("what's type : move, \033[36mdown\u001B[0m, \033[33mup\033[0m | ")

	for i := 0; i < len(split); i++ {
		s := strings.Split(split[i], ",")

		x, _ := strconv.Atoi(s[3])
		y, _ := strconv.Atoi(s[4])
		t, _ := strconv.Atoi(s[2])

		dX := x - lastX
		dY := y - lastY
		dt := t - lastT

		calc := strconv.FormatFloat(math.Sqrt(float64(dX*dX+dY*dY))/float64(dt), 'f', 2, 64)
		var inf string
		switch s[1][0] {
		case '2':
			inf = fmt.Sprintf("\033[36m%s\u001B[0m, ", calc)
		case '3':
			inf = fmt.Sprintf("\033[33m%s\u001B[0m, ", calc)
		default:
			inf = fmt.Sprintf("%s, ", calc)
		}

		lengthWithoutColors := len(fmt.Sprintf("%s, ", calc))
		if lengthWithoutColors > currentWidthLeft {
			buf.WriteString("\n")
			buf.WriteString(tab)
			buf.WriteString(strings.Repeat(" ", indent))
			currentWidthLeft = w
		}
		buf.WriteString(inf)
		currentWidthLeft -= lengthWithoutColors

		lastX = x
		lastY = y
		lastT = t
	}

	buf.Truncate(buf.Len() - 2)
	return
}

func dtdPosCntTactList(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	_, _, _, split := splitTouchData(information) //ts is in millisecond

	if len(split) == 0 {
		buf.WriteString("no touch activity")
		return
	}

	var (
		lastCnt, lastT int
	)

	w := information.Map["x-available-width"].(int)
	indent := information.Map["x-add-indent"].(int)
	tab := information.Map["x-tab"].(string)

	currentWidthLeft := w - indent - len(tab) - 8

	//inform about the color code
	buf.WriteString(color.WhiteString(""))
	buf.WriteString("what's type : move, \033[36mdown\u001B[0m, \033[33mup\033[0m | ")

	for i := 0; i < len(split); i++ {
		s := strings.Split(split[i], ",")

		cnt, _ := strconv.Atoi(s[0])
		t, _ := strconv.Atoi(s[2])

		dcnt := cnt - lastCnt
		dt := t - lastT

		calc := strconv.FormatFloat(math.Sqrt(float64(dcnt))/float64(dt), 'f', 2, 64)
		var inf string
		switch s[1][0] {
		case '2':
			inf = fmt.Sprintf("\033[36m%s\u001B[0m, ", calc)
		case '3':
			inf = fmt.Sprintf("\033[33m%s\u001B[0m, ", calc)
		default:
			inf = fmt.Sprintf("%s, ", calc)
		}

		lengthWithoutColors := len(fmt.Sprintf("%s, ", calc))
		if lengthWithoutColors > currentWidthLeft {
			buf.WriteString("\n")
			buf.WriteString(tab)
			buf.WriteString(strings.Repeat(" ", indent))
			currentWidthLeft = w
		}
		buf.WriteString(inf)
		currentWidthLeft -= lengthWithoutColors

		lastCnt = cnt
		lastT = t
	}

	buf.Truncate(buf.Len() - 2)
	return
}

func sensorSeparator(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	separator := strings.Split(string(information.Map["raw"].([]uint8))[1:], ",2,")[0] + ","
	buf.WriteString(separator)
	return
}

func encryptionKey(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	keys := strings.Split(string(information.Map["encrypted"].([]uint8)), ";")
	buf.WriteString(keys[1])
	return
}

func encryptionTime(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	keys := strings.Split(string(information.Map["encrypted"].([]uint8)), ";")
	buf.WriteString(strings.Split(keys[3], ",")[4])
	return
}

func shufflingKey(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	keys := strings.Split(string(information.Map["encrypted"].([]uint8)), ";")
	buf.WriteString(keys[2])
	return
}

func shufflingTime(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	keys := strings.Split(string(information.Map["encrypted"].([]uint8)), ";")
	buf.WriteString(strings.Split(keys[3], ",")[3])
	return
}

func rValue(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	var split []string
	_, _, _, split = splitPizte(information)
	buf.WriteString(split[23])
	return
}

func rcfp(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	var split []string
	_, _, _, split = splitPizte(information)
	buf.WriteString(split[22])
	return
}

func fpValStr(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	buf.WriteString(strings.Split(information.Map["-70"].(string), ";")[0])
	return
}

func sensorIdTime(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	keys := strings.Split(string(information.Map["encrypted"].([]uint8)), ";")
	buf.WriteString(strings.Split(keys[3], ",")[2])
	return
}

func scriptInitTime(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	keys := strings.Split(string(information.Map["encrypted"].([]uint8)), ";")
	buf.WriteString(strings.Split(keys[3], ",")[1])
	return
}

func sensorBuildTime(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	keys := strings.Split(string(information.Map["encrypted"].([]uint8)), ";")
	buf.WriteString(strings.Split(keys[3], ",")[0])
	return
}

func dtBpd(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	var split []string
	_, _, _, split = splitPizte(information)
	dt1, _ := strconv.Atoi(split[7])
	dt2, _ := strconv.Atoi(split[17])
	buf.WriteString(fmt.Sprintf("dt 1 : %d ms, dt 2 : %d ms, difference : %d ms", dt1, dt2, dt2-dt1))
	return
}

func screenWidth(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	var split []string
	_, _, _, split = getSplitDeviceData(information)
	buf.WriteString(split[12])
	return
}

func screenHeight(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	var split []string
	_, _, _, split = getSplitDeviceData(information)
	buf.WriteString(split[13])
	return
}

func availableWidth(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	var split []string
	_, _, _, split = getSplitDeviceData(information)
	buf.WriteString(split[14])
	return
}

func availableHeight(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	var split []string
	_, _, _, split = getSplitDeviceData(information)
	buf.WriteString(split[15])
	return
}

func innerWidth(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	var split []string
	_, _, _, split = getSplitDeviceData(information)
	buf.WriteString(split[16])
	return
}

func innerHeight(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	var split []string
	_, _, _, split = getSplitDeviceData(information)
	buf.WriteString(split[17])
	return
}

func outerWidth(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	var split []string
	_, _, _, split = getSplitDeviceData(information)
	buf.WriteString(split[18])
	return
}

var ac accounting.Accounting

func formatNumber(number string) string {
	n, _ := strconv.Atoi(number)
	return ac.FormatMoney(n)
}

func JSHeapLimit(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	var split []string
	split = strings.Split(information.Map["-131"].(string), ",")
	buf.WriteString(formatNumber(split[0]))
	return
}

func JSHeapTotal(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	var split []string
	split = strings.Split(information.Map["-131"].(string), ",")
	v1, _ := strconv.Atoi(split[0])
	v2, _ := strconv.Atoi(split[1])
	buf.WriteString(formatNumber(split[1]) + " (limit-total = " + formatNumber(strconv.Itoa(v1-v2)) + ")")
	return
}

func JSHeapUsed(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	var split []string
	split = strings.Split(information.Map["-131"].(string), ",")
	v1, _ := strconv.Atoi(split[1])
	v2, _ := strconv.Atoi(split[2])
	buf.WriteString(formatNumber(split[2]) + " (total-used = " + formatNumber(strconv.Itoa(v1-v2)) + ")")
	return
}

func timezoneOffset(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	var split []string
	split = strings.Split(information.Map["-70"].(string), ";")
	tzRaw, _ := strconv.Atoi(split[7])
	//tzRaw is in minutes
	//write in the buffer the time zone egs : UTC+08:00
	buf.WriteString("UTC")
	if tzRaw < 0 {
		buf.WriteString("-")
	} else {
		buf.WriteString("+")
	}
	tzRaw = int(math.Abs(float64(tzRaw)))
	tzHour := tzRaw / 60
	tzMin := tzRaw % 60
	buf.WriteString(fmt.Sprintf("%02d:%02d", tzHour, tzMin))
	return
}

func lang(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	var _, _, _, split = getSplitDeviceData(information)
	buf.WriteString(split[4])
	return
}

func langLen(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	var split []string
	split = strings.Split(information.Map["-131"].(string), ",")
	if len(split) > 4 {
		buf.WriteString(split[4])
		return
	}
	buf.WriteString("<unknown>")
	return
}

func langHash(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	var split []string
	split = strings.Split(information.Map["-129"].(string), ",")
	if len(split) > 2 {
		buf.WriteString(split[2])
		return
	}
	buf.WriteString("<unknown>")
	return
}
