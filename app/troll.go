package main

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	servertiming "github.com/mitchellh/go-server-timing"

	"github.com/bendoerr/trollr/exec"
)

type Troll struct {
	path     string
	executor exec.Executor
	timeout  time.Duration
	max      int
}

func NewTroll(path string, executor exec.Executor) *Troll {
	return &Troll{
		path:     path,
		executor: executor,
		timeout:  30 * time.Second,
		max:      127,
	}
}

type Roll []string

type Rolls []Roll

type Probability float64

type Probabilities map[string]Probability

type Cumulative string

var (
	CumulativeGreaterThan        Cumulative = "gt"
	CumulativeGreaterThanOrEqual Cumulative = "ge"
	CumulativeLesserThan         Cumulative = "lt"
	CumulativeLesserThanOrEqual  Cumulative = "le"
)

func CumulativeFromString(s string) (Cumulative, error) {
	switch s {
	case string(CumulativeGreaterThan):
		return CumulativeGreaterThan, nil
	case string(CumulativeGreaterThanOrEqual):
		return CumulativeGreaterThanOrEqual, nil
	case string(CumulativeLesserThan):
		return CumulativeLesserThan, nil
	case string(CumulativeLesserThanOrEqual):
		return CumulativeLesserThanOrEqual, nil
	default:
		return Cumulative(""), fmt.Errorf("value '%s' is not a valid accumulate setting", s)
	}
}

func (a Cumulative) String() string {
	return string(a)
}

// swagger:model RollsResult
type RollsResult struct {
	Definition string `json:",omitempty"`
	NumTimes   int    `json:",omitempty"`
	Runtime    int64  `json:",omitempty"`
	Rolls      Rolls  `json:",omitempty"`
	RollsRaw   string `json:"-"`
	Err        error  `json:"-"`
	Error      string `json:",omitempty"`
}

// swagger:model CalcResult
type CalcResult struct {
	Cumulative       string        `json:",omitempty"`
	Average          Probability   `json:",omitempty"`
	Definition       string        `json:",omitempty"`
	Err              error         `json:"-"`
	Error            string        `json:",omitempty"`
	MeanDeviation    Probability   `json:",omitempty"`
	ProbabilitiesCum Probabilities `json:",omitempty"`
	ProbabilitiesEq  Probabilities `json:",omitempty"`
	ProbabilitiesRaw string        `json:"-"`
	Runtime          int64         `json:",omitempty"`
	Spread           Probability   `json:",omitempty"`
}

var ProbabilitiesMatcher = regexp.MustCompile(`^([^:]*): *([0-9.]*) *([0-9.]*)`)

var StatsMatcher = regexp.MustCompile(`Average = ([0-9.]*) *Spread = ([0-9.]*) *Mean deviation = ([0-9.]*)`)

var NumbersMatcher = regexp.MustCompile(`^[0-9\s]*$`)

func (t *Troll) MakeRolls(ctx context.Context, num int, definition string) RollsResult {
	r := RollsResult{
		Definition: definition,
		NumTimes:   num,
		Rolls:      make([]Roll, 0),
	}

	if num < 1 {
		num = 1
		r.NumTimes = 1
	}

	if num > 127 {
		r.Err = fmt.Errorf("exceeded max rolls %v > 127", num)
		r.Error = r.Err.Error()
		return r
	}

	timeout, timeoutCancel := context.WithTimeout(ctx, 15*time.Second)
	defer timeoutCancel()

	timing := servertiming.FromContext(ctx)
	metric := timing.NewMetric("submit").Start()
	xr := t.executor(timeout, t.path, &definition, strconv.Itoa(num))
	metric.Stop()
	r.Runtime = metric.Duration.Milliseconds()

	if xr.Err() != nil {
		r.Err = xr.Err()
		r.Error = xr.Err().Error()
		return r
	}

	var line string
	var err error
	for {
		line, err = xr.Stdout().ReadString('\n')
		r.RollsRaw = r.RollsRaw + line
		if err != nil {
			break
		}
		line := strings.TrimSpace(line)

		if len(line) < 1 {
			continue
		}

		if NumbersMatcher.MatchString(line) {
			r.Rolls = append(r.Rolls, strings.Split(line, " "))
		} else {
			r.Rolls = append(r.Rolls, []string{line})
		}
	}

	return r
}

func (t *Troll) CalcRoll(ctx context.Context, definition string, cumulative string) CalcResult {
	r := CalcResult{
		Definition:       definition,
		ProbabilitiesCum: make(Probabilities),
		ProbabilitiesEq:  make(Probabilities),
	}

	var err error

	acc := CumulativeGreaterThanOrEqual
	if len(cumulative) > 0 {
		acc, err = CumulativeFromString(cumulative)
		if err != nil {
			r.Err = err
			r.Error = err.Error()
			return r
		}
	}
	r.Cumulative = acc.String()

	timeout, timeoutCancel := context.WithTimeout(ctx, 15*time.Second)
	defer timeoutCancel()

	timing := servertiming.FromContext(ctx)
	metric := timing.NewMetric("submit").Start()
	xr := t.executor(timeout, t.path, &definition, "0", acc.String())
	metric.Stop()
	r.Runtime = metric.Duration.Milliseconds()

	if xr.Err() != nil {
		r.Err = xr.Err()
		r.Error = xr.Err().Error()
		return r
	}

	var last_line string
	var header_done bool

	for {
		line, err := xr.Stdout().ReadString('\n')
		if err != nil {
			break
		}
		r.ProbabilitiesRaw = r.ProbabilitiesRaw + line
		line = strings.TrimSpace(line)

		if len(line) < 0 {
			continue
		}

		if !header_done {
			header_done = true
			continue
		}

		if m := ProbabilitiesMatcher.FindAllStringSubmatch(line, -1); len(m) > 0 {
			ms := m[0]
			key := ms[1]
			if len(last_line) > 0 {
				key = strings.TrimSpace(last_line + "\n" + key)
				last_line = ""
			}
			f, err := strconv.ParseFloat(ms[2], 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			r.ProbabilitiesEq[key] = Probability(f)

			f, err = strconv.ParseFloat(ms[3], 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			r.ProbabilitiesCum[key] = Probability(f)
		} else if m := StatsMatcher.FindAllStringSubmatch(line, -1); len(m) > 0 {
			ms := m[0]
			f, err := strconv.ParseFloat(ms[1], 64)
			if err != nil {
				fmt.Println(err)
				continue
			}

			r.Average = Probability(f)

			f, err = strconv.ParseFloat(ms[2], 64)
			if err != nil {
				fmt.Println(err)
				continue
			}

			r.Spread = Probability(f)

			f, err = strconv.ParseFloat(ms[3], 64)
			if err != nil {
				fmt.Println(err)
				continue
			}

			r.MeanDeviation = Probability(f)
		} else {
			if len(last_line) > 1 {
				last_line = last_line + "\n" + line
			} else {
				last_line = line
			}
		}
	}

	return r
}
