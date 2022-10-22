package main

import (
	"encoding/json"
	"log"
	"reflect"
	"sync"
	"testing"
	"time"
)

type PathsTest struct {
	Name  string
	Input string
	Want  []string
}

func (pt PathsTest) TestreceivePaths(t *testing.T) {

	var testCases = []PathsTest{
		{
			Name:  "Single Input",
			Input: `{"forward": "exit"}`,
			Want:  []string{"[forward]"},
		},
		{
			Name:  "No exit path",
			Input: `{"forward": "tiger", "left": "ogre", "right":"demon"}`,
			Want:  []string{"[Sorry]"},
		},
		{
			Name:  "Single exit",
			Input: `{"forward": "tiger", "left": {"forward": {"upstairs": "exit"}, "left": "dragon"}, "right": {"forward": "dead end"}}`,
			Want:  []string{"[left forward upstairs]"},
		},
		{
			Name:  "Multiple exit points",
			Input: `{"forward": "tiger", "left": {"forward": {"upstairs": "exit"}, "left": "dragon"}, "right": {"forward": "dead end"}, "upstairs": {"forward": {"left": "dead end", "right": {"forward": {"left": "exit"}} }}}`,
			Want:  []string{"[left forward upstairs]"},
		},
	}

	wg := sync.WaitGroup{}
	Ch := make(chan Path)
	// Dn := make(chan string)
	var Res []Path

	for _, tc := range testCases {
		I := make(map[string]interface{})
		err := json.Unmarshal([]byte(tc.Input), &I)
		if err != nil {
			log.Printf("%v", err)
		}

		var Tr []string
		wg.Add(len(I))
		for i, v := range I {
			go processPaths(Tr, i, v, &wg, Ch)
		}

		got := receivePaths(&Res, Ch)
		time.Sleep(time.Millisecond * 10)
		if !reflect.DeepEqual(tc.Want, got) {
			t.Errorf("Line %q, wanted %v, got %v", tc.Input, tc.Want, got)
		}
	}
}
