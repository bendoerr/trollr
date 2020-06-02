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

type RollsResult struct {
	Definition string `json:",omitempty"`
	NumTimes   int    `json:",omitempty"`
	Runtime    int64  `json:",omitempty"`
	Rolls      Rolls  `json:",omitempty"`
	Err        error  `json:"-"`
	Error      string `json:",omitempty"`
}

type CalcResult struct {
	Cumulative       string        `json:",omitempty"`
	Average          Probability   `json:",omitempty"`
	Definition       string        `json:",omitempty"`
	Err              error         `json:"-"`
	Error            string        `json:",omitempty"`
	MeanDeviation    Probability   `json:",omitempty"`
	ProbabilitiesCum Probabilities `json:",omitempty"`
	ProbabilitiesEq  Probabilities `json:",omitempty"`
	Runtime          int64         `json:",omitempty"`
	Spread           Probability   `json:",omitempty"`
}

var ProbabilitiesMatcher = regexp.MustCompile(`^([^:]*): *([0-9.]*) *([0-9.]*)`)

var StatsMatcher = regexp.MustCompile(`Average = ([0-9.]*) *Spread = ([0-9.]*) *Mean deviation = ([0-9.]*)`)

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

	timing := servertiming.FromContext(ctx)
	metric := timing.NewMetric("submit").Start()
	xr := t.executor(ctx, t.path, &definition, strconv.Itoa(num))
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

		if err != nil {
			break
		}

		r.Rolls = append(r.Rolls, strings.Split(strings.TrimSpace(line), " "))
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

	timing := servertiming.FromContext(ctx)
	metric := timing.NewMetric("submit").Start()
	xr := t.executor(ctx, t.path, &definition, "0", acc.String())
	metric.Stop()
	r.Runtime = metric.Duration.Milliseconds()

	if xr.Err() != nil {
		r.Err = xr.Err()
		r.Error = xr.Err().Error()
		return r
	}

	for {
		line, err := xr.Stdout().ReadString('\n')

		line = strings.TrimSpace(line)
		if len(line) < 0 {
			continue
		}

		if err != nil {
			break
		}

		if m := ProbabilitiesMatcher.FindAllStringSubmatch(line, -1); len(m) > 0 {
			ms := m[0]
			f, err := strconv.ParseFloat(ms[2], 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			r.ProbabilitiesEq[ms[1]] = Probability(f)

			f, err = strconv.ParseFloat(ms[3], 64)
			if err != nil {
				fmt.Println(err)
				continue
			}
			r.ProbabilitiesCum[ms[1]] = Probability(f)
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
		}
	}

	return r
}
