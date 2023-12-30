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

type funcRename string

var AllInformation = []any{
	scriptId,
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
		"browser",
		[]any{
			userAgent,
			doNotTrack,
			"webgl",
			[]any{
				webglVendor,
				webglRenderer,
			},
			"canvas",
			[]any{
				rcfp,
				rValue,
				fpValStr,
			},
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
		"-133",
		[]any{
			funcRename("window.Object.getOwnPropertyDescriptors"),
			windowObjectProp,
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
				funcRename("dt_information"),
				dtMact,
				funcRename("dt_list"),
				dtMactList,
				funcRename("acceleration"),
				mouseMovementAcceleration,
				funcRename("velocity"),
				mouseMovementVelocity,
			},
			"kact",
			[]any{
				funcRename("dt_information"),
				dtKact,
				funcRename("dt_list"),
				dtKactList,
			},
			"tact",
			[]any{
				funcRename("dt_information"),
				dtTact,
				funcRename("dt_list"),
				dtTactList,
				funcRename("d_pos_list"),
				dPosTactList,
				funcRename("dt_delta_position"),
				dtdPosTactList,
				funcRename("velocity"),
				tactVelocity,
				funcRename("ratio counter movement"),
				dtdPosCntTactList,
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

func displaySpecificInformation(information utils.OrderedMap, last bool, tab string, fn func(om utils.OrderedMap) (buf *bytes.Buffer), fnName string) (buf *bytes.Buffer) {
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
	if fnName == "" {
		fnName = GetFunctionName(fn)
	}

	addLen := len(fnName) + len(": ")
	avail := terminalWidth - len(tab) - addLen

	information.Map["x-add-indent"] = addLen
	information.Map["x-available-width"] = avail

	if !last {
		information.Map["x-tab"] = tab + "│  "
	} else {
		information.Map["x-tab"] = tab
	}

	b := fn(information)

	buf.WriteString(fnName + ": \u001B[36m")
	buf.WriteString(b.String())
	buf.WriteString("\n")

	// Reset color and return buffer
	buf.WriteString("\u001B[0m")

	if b.Len() > avail {
		buf.WriteString(information.Map["x-tab"].(string))
		buf.WriteString("\n")
	}

	return
}

func display(om utils.OrderedMap, list []any, tab string) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	listLength := len(list)
	for i := 0; i < listLength; i++ {
		k := list[i]

		fnName := ""

		if i > 0 {
			switch list[i-1].(type) {
			case funcRename:
				fnName = string(list[i-1].(funcRename))
			}
		}

		switch k.(type) {
		case func(utils.OrderedMap) (buf *bytes.Buffer):
			buf.WriteString(displaySpecificInformation(om, i == listLength-1, tab, k.(func(utils.OrderedMap) (buf *bytes.Buffer)), fnName).String())

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

func scriptId(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	sid, _ := information.Map["sensor_data"]
	v := sid.([]string)[0]
	buf.WriteString(v[len(v)-24:])
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

	w := information.Map["x-available-width"].(int) - 3
	indent := information.Map["x-add-indent"].(int)
	tab := information.Map["x-tab"].(string)

	//inform about the color code
	buf.WriteString(color.WhiteString(""))
	buf.WriteString("mouse event colors -> move, \u001B[32mclick\u001B[0m, \033[36mmousedown\u001B[0m, \033[33mmouseup\033[0m, \033[90munknown\033[0m")
	buf.WriteString("\n")
	buf.WriteString(tab)
	buf.WriteString(strings.Repeat(" ", indent))

	currentWidthLeft := w
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
			inf = fmt.Sprintf("\033[90m%d\u001B[0m, ", dt)
		}

		lengthWithoutColors := len(fmt.Sprintf("%d, ", dt))
		if lengthWithoutColors > currentWidthLeft {
			buf.WriteString("\n")
			buf.WriteString(tab)
			buf.WriteString(strings.Repeat(" ", indent))
			currentWidthLeft = w
		}

		currentWidthLeft -= lengthWithoutColors
		buf.WriteString(inf)
	}

	buf.Truncate(buf.Len() - 2)

	return
}

func mouseMovementAcceleration(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	_, _, _, split := splitMouseData(information) //ts is in millisecond

	if len(split) == 0 {
		buf.WriteString("no mouse movement")
		return
	}

	w := information.Map["x-available-width"].(int) - 3
	indent := information.Map["x-add-indent"].(int)
	tab := information.Map["x-tab"].(string)

	currentWidthLeft := w

	for i := 0; i < len(split)-2; i++ {
		p1, p2, p3 := strings.Split(split[i], ","), strings.Split(split[i+1], ","), strings.Split(split[i+2], ",")

		var (
			x1, _ = strconv.Atoi(p1[3])
			y1, _ = strconv.Atoi(p1[4])
			t1, _ = strconv.Atoi(p1[2])
			x2, _ = strconv.Atoi(p2[3])
			y2, _ = strconv.Atoi(p2[4])
			t2, _ = strconv.Atoi(p2[2])
			x3, _ = strconv.Atoi(p3[3])
			y3, _ = strconv.Atoi(p3[4])
			t3, _ = strconv.Atoi(p3[2])
		)

		v1 := math.Sqrt(math.Pow(float64(x2-x1), 2)+math.Pow(float64(y2-y1), 2)) / float64(t2-t1)
		v2 := math.Sqrt(math.Pow(float64(x3-x2), 2)+math.Pow(float64(y3-y2), 2)) / float64(t3-t2)

		acceleration := (v2 - v1) / float64(t3-t1)
		inf := strconv.FormatFloat(acceleration, 'f', 2, 64)
		if i < len(split)-3 {
			inf += ", "
		}
		if len(inf) > currentWidthLeft {
			buf.WriteString("\n")
			buf.WriteString(color.HiWhiteString(tab))
			buf.WriteString(strings.Repeat(" ", indent))
			currentWidthLeft = w
		}

		currentWidthLeft -= len(inf)
		buf.WriteString(fmt.Sprintf("\033[36m%s\u001B[0m", inf))
	}

	if buf.Len() > 2 {
		buf.Truncate(buf.Len() - 2)
	} else {
		buf.WriteString("no acceleration")
	}
	return
}

func mouseMovementVelocity(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	_, _, _, split := splitMouseData(information) //ts is in millisecond

	if len(split) == 0 {
		buf.WriteString("no mouse movement")
		return
	}

	w := information.Map["x-available-width"].(int) - 3
	indent := information.Map["x-add-indent"].(int)
	tab := information.Map["x-tab"].(string)

	currentWidthLeft := w

	for i := 0; i < len(split)-1; i++ {
		p1, p2 := strings.Split(split[i], ","), strings.Split(split[i+1], ",")

		var (
			x1, _ = strconv.Atoi(p1[3])
			y1, _ = strconv.Atoi(p1[4])
			t1, _ = strconv.Atoi(p1[2])
			x2, _ = strconv.Atoi(p2[3])
			y2, _ = strconv.Atoi(p2[4])
			t2, _ = strconv.Atoi(p2[2])
		)

		// Calculate the distance between point 1 and point 2
		distance := math.Sqrt(math.Pow(float64(x2-x1), 2) + math.Pow(float64(y2-y1), 2))

		// Calculate the time difference
		timeDiff := float64(t2 - t1)

		// Calculate velocity
		// Make sure to handle the case where timeDiff is 0 to avoid division by zero
		var vel float64
		if timeDiff != 0 {
			vel = distance / timeDiff
		} else {
			vel = 0
		}

		// Format and write the velocity value
		inf := strconv.FormatFloat(vel, 'f', 2, 64)
		if i < len(split)-2 {
			inf += ", "
		}
		if len(inf) > currentWidthLeft {
			buf.WriteString("\n")
			buf.WriteString(color.HiWhiteString(tab))
			buf.WriteString(strings.Repeat(" ", indent))
			currentWidthLeft = w
		}

		currentWidthLeft -= len(inf)
		buf.WriteString(fmt.Sprintf("\033[36m%s\u001B[0m", inf))
	}

	if buf.Len() > 2 {
		buf.Truncate(buf.Len() - 2)
	} else {
		buf.WriteString("no acceleration")
	}
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

func dtKact(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	_, _, _, split := splitKeyboardData(information) //ts is in millisecond

	if len(split) == 0 {
		buf.WriteString("no keyboard activity")
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

		if dt < 1000 && dt > 0 {
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

func dtKactList(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	_, _, _, split := splitKeyboardData(information) //ts is in millisecond

	if len(split) == 0 {
		buf.WriteString("no keyboard activity")
		return
	}

	lastT := 0

	w := information.Map["x-available-width"].(int)
	indent := information.Map["x-add-indent"].(int)
	tab := information.Map["x-tab"].(string)

	currentWidthLeft := w - indent - len(tab) - 8

	//inform about the color code
	buf.WriteString(color.WhiteString(""))
	buf.WriteString("what's type : \033[36mkeydown\u001B[0m, keypress, \033[33mkeyup\033[0m | ")
	for i := 0; i < len(split); i++ {
		s := strings.Split(split[i], ",")
		t, _ := strconv.Atoi(s[2])

		dt := t - lastT
		lastT = t

		var inf string
		switch s[1][0] {
		case '1':
			inf = fmt.Sprintf("\u001B[36m%d\u001B[0m, ", dt)
		case '2':
			inf = fmt.Sprintf("\033[33m%d\033[0m, ", dt)
		case '3':
			inf = fmt.Sprintf("%d, ", dt)
		default:
			continue
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

func windowObjectProp(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	field := information.Map["-133"].(string)
	if field == "" {
		buf.WriteString("\033[90m<nothing>\033[0m")
	} else {
		buf.WriteString(field)
	}
	return buf
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

func webglVendor(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	first := strings.Split(information.Map["-129"].(string), ";wl3")[0]
	second := strings.Split(first, ",")

	if len(second) > 1 {
		buf.WriteString(second[len(second)-1])
		return
	}

	buf.WriteString("<unknown>")
	return
}

func webglRenderer(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	first := strings.Split(information.Map["-129"].(string), ";wl4")[0]
	second := strings.Split(first, ",ANGLE (")

	if len(second) > 1 {
		buf.WriteString("ANGLE (" + second[1])
		return
	}

	buf.WriteString("<unknown>")
	return
}

func doNotTrack(information utils.OrderedMap) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	split := strings.Split(information.Map["-70"].(string), ";")
	buf.WriteString(split[len(split)-1])
	return
}
