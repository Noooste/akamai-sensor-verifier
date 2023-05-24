package tools

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net/url"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
)

func CheckChallenge(startTime int64, abck string, challenge1, challenge2, challenge3 *string) {
	if strings.HasSuffix(abck, "~-1~-1~-1") || strings.HasSuffix(abck, "~-1~||-1||~-1") {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			log.Print(r)
			debug.PrintStack()
		}
	}()

	bmak := Bmak{}
	bmak.StartTs = startTime
	bmak.CookieValue = abck
	bmak.MnAbck = abck
	c1, c2, c3 := bmak.SolveChallenge()

	*challenge1 = c1
	*challenge2 = c2
	*challenge3 = c3
}

func (bm *Bmak) SolveChallenge() (string, string, string) {
	bm.MnTs = strconv.Itoa(int(bm.StartTs))
	bm.MnMcLmt = 10
	bm.MnTout = 100
	bm.MnStout = 1000
	bm.MnCd = 10000
	bm.StartTs += int64(randIntRange(100, 200))
	bm.MnLc = make([]string, 1)
	bm.MnLd = make([]int, 1)
	bm.MnAl = make([]string, 1)
	bm.MnIl = make([]int, 1)
	bm.MnR = make(map[string]string)
	bm.MnTcl = make([]int, 1)

	params := bm.getMnParamsFromAbck()

	if len(params) == 0 {
		return "", "", ""
	}

	newParams := bm.mnGetNewChallengeParams(params)

	if newParams != nil {
		bm.mnUpdateChallengeDetails(newParams)
		if bm.MnSen != 0 {
			bm.MnState = 1
			bm.MnAl = make([]string, 10)
			bm.MnIl = make([]int, 10)
			bm.MnTcl = make([]int, 10)
			bm.Mnrts = int(getCfDate())
			bm.MnRt = bm.Mnrts - int(bm.StartTs)
		}
	}

	bm.mnW()

	C := bm.mnGetCurrentChallenges(params)
	B := ""
	x := ""
	M := ""

	if len(C) >= 2 && len(C[1]) != 0 {
		j := C[1]
		if len(bm.MnR[j]) != 0 {
			B = bm.MnR[j]
		}
	}
	if len(C) >= 3 && len(C[2]) != 0 {
		A := C[2]
		if len(bm.MnR[A]) != 0 {
			x = bm.MnR[A]
		}
	}
	if len(C) >= 4 && len(C[3]) != 0 {
		L := C[3]
		if len(bm.MnR[L]) != 0 {
			M = bm.MnR[L]
		}
	}

	return B, x, M
}

func (bm *Bmak) getMnParamsFromAbck() [][]any {
	var t [][]any

	var a = bm.CookieValue

	if true {
		aUnescaped, _ := url.QueryUnescape(a)
		var e = strings.Split(aUnescaped, "~")

		if len(e) >= 5 {
			n := e[0]
			o := e[4]
			m := strings.Split(o, "||")

			if len(m) > 0 {
				for r := 0; r < len(m); r++ {
					i := m[r]
					c := strings.Split(i, "-")

					if len(c) > 5 {
						b, _ := strconv.Atoi(c[0])
						d := c[1]
						s, _ := strconv.Atoi(c[2])
						k, _ := strconv.Atoi(c[3])
						l, _ := strconv.Atoi(c[4])

						u := 1
						if len(c) >= 6 {
							u, _ = strconv.Atoi(c[5])
						}

						var v_ = []any{b, n, d, s, k, l, u}

						if u == 2 {
							t = append(t, []any{})
							t = append([][]any{v_}, t...)
						} else {
							t = append(t, v_)
						}
					}
				}
			}
		}
	}
	return t
}

func (bm *Bmak) mnGetNewChallengeParams(t [][]any) []any {
	var a = -99999
	var e = -99999
	var n = -99999

	//checks if array is nil
	if t != nil {
		//makes iterator that mathces the length of the t array
		for o := 0; o < len(t); o++ {
			//assigns m to the first o-th of the t array, which is another array
			m := t[o]
			//checks if m is bigger than 0, so it always works only for the 1st element of the array as its always [[..], []]
			if len(m) > 0 {
				//we do type assertion as the 1st element of the inner array is always an int
				r := m[0].(int)
				//not sure what i is but its the same as the akamai script
				i := bm.MnAbck + fmt.Sprint(bm.StartTs) + fmt.Sprint(m[2])
				//we do type assertion as the 6th element of the inner array is always an int
				c := m[6].(int)

				//we do type assertion as we already know the value type
				var b = 0
				for b = 0; b < 0 && (r == 1 && bm.MnLc[b] != i); b++ {
				}
				if 0 == b {
					a = o
					if 2 == c {
						e = o
					}
					if 3 == c {
						n = o
					}
				}
			}
		}

		if n != -99999 && bm.Pstate {
			return t[n]
		} else {
			if e == -99999 || bm.Pstate {
				if a == -99999 || bm.Pstate {
					return nil
				} else {
					return t[a]
				}
			} else {
				return t[e]
			}
		}
	}
	return nil
}

func (bm *Bmak) mnUpdateChallengeDetails(t []any) {
	bm.MnSen = t[0].(int)
	bm.MnAbck = t[1].(string)
	bm.MnPsn = t[2].(string)
	bm.MnCd = t[3].(int)
	bm.MnTout = t[4].(int)
	bm.MnStout = t[5].(int)
	bm.MnTs = strconv.Itoa(int(bm.StartTs))
	bm.MnCc = bm.MnAbck + strconv.Itoa(int(bm.StartTs)) + bm.MnPsn
}

func Randint(min, max int) int {
	return min + rand.Intn(max-min)
}

const (
	dot0 = "0."
)

func randomHex(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (bm *Bmak) mnW() {
	var r bytes.Buffer
	r.WriteString(bm.MnCc)
	truncateLength := r.Len()

	for index := 0; index < bm.MnMcLmt; index++ {
		if index >= 0 {
			r.Truncate(truncateLength)
		}

		m := bm.MnCd + index
		r.WriteString(fmt.Sprint(m))
		r.WriteString(dot0)

		g := randomHex(11)
		r.WriteString(g)
		tl := r.Len()

		var sum [32]byte
		for {
			rd := randomHex(3)
			r.WriteString(rd)
			sum = sha256.Sum256(r.Bytes())
			if bdm(sum, m) == 0 {
				var iter = Randint(100, 1500)
				bm.MnAl[index] = dot0 + g + rd
				bm.MnTcl[index] = iter/10 - Randint(1, 3)
				bm.MnIl[index] = iter
				if index == 0 {
					bm.MnLg = []string{
						bm.MnAbck,
						bm.MnTs,
						bm.MnPsn,
						bm.MnCc,
						fmt.Sprint(bm.MnCd),
						fmt.Sprint(m),
						dot0 + g + rd,
						r.String(),
						strings.Join(ChangeListToString(sum), ","),
						strconv.Itoa(int(bm.StartTs - getCfDate())),
						"0",
						strconv.Itoa(int(getCfDate())),
					}
				}
				break
			} else {
				r.Truncate(tl)
			}
		}
	}
	bm.MnLc[0] = bm.MnCc
	bm.MnLd[0] = bm.MnCd
	bm.MnState = 0
	bm.MnR[bm.MnAbck+bm.MnPsn] = bm.mnPr()
	return
}

func (bm *Bmak) mnGetCurrentChallenges(t [][]any) []string {
	a := make([]string, 3)

	if t != nil {
		for e := 0; e < len(t); e++ {
			n := t[e]
			if len(n) > 0 {
				o := n[1].(string) + n[2].(string)
				m := n[6].(int)
				if m >= len(a) {
					a = append(a, o)
				}
				a[m] = o
			}
		}
	}
	return a
}

func randIntRange(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func ChangeListToString(list any) (finalList []string) {
	finalList = make([]string, reflect.ValueOf(list).Len())
	for i := 0; i < reflect.ValueOf(list).Len(); i++ {
		finalList[i] = fmt.Sprint(reflect.ValueOf(list).Index(i))
	}
	return
}
