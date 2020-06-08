package main_test

import (
	"context"
	"github.com/bendoerr/trollr/exec"
	jsoniter "github.com/json-iterator/go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"io"
	"io/ioutil"
	"path"
	"runtime"

	. "github.com/bendoerr/trollr/app"
)

var _ = Describe("Troll", func() {

	DescribeTable("parsing rolls",
		func(rollName string) {
			rollDef := testContent(rollName, "def")
			stdout := testContent(rollName, "result_roll")
			expectedResult := testContent(rollName, "json_roll")
			executor := exec.NewTestExecutor(stdout, "", 0, nil)
			troll := NewTroll("", executor)

			result := troll.MakeRolls(context.Background(), 1, rollDef)
			Expect(result.Err).To(BeNil())

			resultJson, _ := jsoniter.MarshalToString(result.Rolls)
			Expect(resultJson).To(MatchJSON(expectedResult))
		},
		Entry("4d4 n=1", "4d4 n=1"),
		Entry("4d4 n=3", "4d4 n=3"),
		Entry("blades in the dark n=1", "blades in the dark n=1"),
		Entry("blades in the dark n=4", "blades in the dark n=4"),
		Entry("savage worlds d8", "savage worlds d8"),
		Entry("sum 3d6 n=1", "sum 3d6 n=1"),
		Entry("sum 3d6 n=3", "sum 3d6 n=3"),
	)

	DescribeTable("parsing calcs",
		func(rollName string) {
			rollDef := testContent(rollName, "def")
			stdout := testContent(rollName, "result_calc")
			expectedResult1 := testContent(rollName, "json_calc_eq")
			expectedResult2 := testContent(rollName, "json_calc_ge")
			executor := exec.NewTestExecutor(stdout, "", 0, nil)
			troll := NewTroll("", executor)

			result := troll.CalcRoll(context.Background(), rollDef, "ge")
			Expect(result.Err).To(BeNil())

			resultJson, _ := jsoniter.MarshalToString(result.ProbabilitiesEq)
			Expect(resultJson).To(MatchJSON(expectedResult1))

			resultJson, _ = jsoniter.MarshalToString(result.ProbabilitiesCum)
			Expect(resultJson).To(MatchJSON(expectedResult2))
		},
		Entry("4d4 n=1", "4d4 n=1"),
		Entry("blades in the dark n=1", "blades in the dark n=1"),
		Entry("savage worlds d8", "savage worlds d8"),
		Entry("sum 3d6 n=1", "sum 3d6 n=1"),
	)
})

var testRolls = func() map[string]string {
	_, thisFilename, _, _ := runtime.Caller(1)
	testRollsDir := path.Join(path.Dir(thisFilename), "test_rolls")
	testRollDirs, err := ioutil.ReadDir(testRollsDir)
	if err != nil {
		panic(err)
	}

	results := make(map[string]string)
	for i := range testRollDirs {
		if testRollDirs[i].IsDir() {
			results[testRollDirs[i].Name()] = path.Join(testRollsDir, testRollDirs[i].Name())
		}
	}

	return results
}()

func testContent(rollName string, file string) string {
	bs, err := ioutil.ReadFile(path.Join(testRolls[rollName], file))
	if err != nil && err != io.EOF {
		panic(err)
	}
	return string(bs)
}
