package trollr

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bendoerr/trollr/util"
)

type Troll struct {
	path     string
	executor util.Executor
	timeout  time.Duration
	max      int
}

func NewTroll(path string, executor util.Executor) *Troll {
	return &Troll{
		path:     path,
		executor: executor,
		timeout:  30 * time.Second,
		max:      127,
	}
}

type Roll []string

type Rolls []Roll

type RollsResult struct {
	Definition string `json:",omitempty"`
	NumTimes   int    `json:",omitempty"`
	Runtime    int64  `json:",omitempty"`
	Rolls      Rolls  `json:",omitempty"`
	Err        error  `json:"-"`
	Error      string `json:",omitempty"`
}

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

	start := time.Now()
	xr := t.executor(ctx, t.path, &definition, strconv.Itoa(num))
	r.Runtime = time.Since(start).Milliseconds()

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
