package tools

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Bmak struct {
	StartTs                                                     int64
	CookieValue, MnAbck, MnTs, MnPsn, MnCc                      string
	MnMcLmt, MnTout, MnStout, MnCd, MnSen, MnState, Mnrts, MnRt int
	Pstate                                                      bool

	MnAl, MnLg, MnLc  []string
	MnTcl, MnIl, MnLd []int
	MnR               map[string]string
}

// getCfDate, its basically an almost mirror function of akamai's get_cf_date() func
func getCfDate() int64 {
	return time.Now().UTC().UnixNano() / 1e6
}

func CalDis(list [6]float64) int {
	a := list[0] - list[1]
	e := list[2] - list[3]
	n := list[4] - list[5]
	return int(math.Sqrt(a*a + e*e + n*n))
}

func Jrs(t int64) (int, int) {
	var (
		rd    = math.Floor(1e5*rand.Float64() + 1e4)
		e     = strconv.FormatFloat(float64(t)*rd, 'f', -1, 64)
		n     uint8
		oList [6]float64
		m     = len(e) >= 18
	)

	for i := 0; i < 6; i++ {
		value, _ := strconv.Atoi(e[n : n+2])
		oList[i] = float64(value)
		if m {
			n += 3
		} else {
			n += 2
		}
	}

	return int(rd), CalDis(oList)
}

func bdm(t [32]byte, a int) int {
	e := int(t[0])
	for n := 1; n < 32; n++ {
		e = (e<<8 + int(t[n])) % a
	}
	return e
}

func (bm *Bmak) mnPr() string {
	return fmt.Sprintf("%s;%s;%s;%s;", strings.Join(bm.MnAl, ","), strings.Join(convertIntToString(bm.MnTcl), ","), strings.Join(convertIntToString(bm.MnIl), ","), strings.Join(bm.MnLg, ","))
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

func Od(a, t string) string {
	var e []string
	n := len(t)
	if n > 0 {
		e = make([]string, len(a))
		for o, v := range a {
			r := string(v)
			i := t[o%n]

			vInt := int(v)
			m := rir(vInt, 47, 57, int(i))
			if m != vInt {
				r = string(rune(m))
			}

			e[o] = r
		}

		if len(e) > 0 {
			return strings.Join(e, "")
		}
	}
	return a
}

func rir(a, b, c, d int) int {
	if a > 47 && a <= 57 {
		a += d % (c - b)
		if a > c {
			a = a - c + b
		}
	}
	return a
}

func Cc(t int) func(t, a int) int64 {
	var a = t % 4
	if a == 2 {
		a = 3
	}
	var e = 42 + a
	n := func(t, a int) int64 {
		return 0
	}
	if e == 42 {
		n = func(t, a int) int64 {
			return int64(t * a)
		}
	} else if e == 43 {
		n = func(t, a int) int64 {
			return int64(t + a)
		}
	} else {
		n = func(t, a int) int64 {
			return int64(t - a)
		}
	}

	return n
}
