package exec

import (
	"reflect"
	"testing"
)

type testStruct struct {
	FirstParam   string
	SecondParam  string
	ThirdParam   bool
	FourthPatam  bool
	FifthParam   int
	SixthParam   string
	ZeventhParam int
}

func TestBuildArgs(t *testing.T) {

	tstrct := testStruct{
		FirstParam:  "val1",
		SecondParam: "val2",
		ThirdParam:  true,
		FourthPatam: false,
		FifthParam:  8,
	}

	got := BuildArgs(
		"command1",
		"--opt1",
		567,
		tstrct,
		[]string{"param1", "param2"},
		[]int{7465, 23},
	)

	want := []string{
		"command1",
		"--opt1",
		"567",
		"--first-param", "val1",
		"--second-param", "val2",
		"--third-param",
		"--fifth-param", "8",
		"param1",
		"param2",
		"7465",
		"23",
	}

	for i := 0; i < len(got); i++ {
		if got[i] != want[i] {
			t.Errorf("\ngot:  %v\nwant: %v\n", got[i], want[i])
		}
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot:  %v\nwant: %v\n", got, want)
	}
}
