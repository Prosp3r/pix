package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

//multiple paths to result case
var INPUT = `{"forward": "tiger", "left": {"forward": {"upstairs": "exit"}, "left": "dragon"}, "right": {"forward": "dead end"}, "upstairs": {"forward": {"left": "dead end", "right": {"forward": {"left": "exit"}} }}}`

//single path to result case
// var INPUT = `{"forward": "tiger", "left": {"forward": {"upstairs": "exit"}, "left": "dragon"}, "right": {"forward": "dead end"}}`

//case fail
// var INPUT = `{"forward": "tiger", "left": "ogre", "right":"demon"}`

//case single string
// var INPUT = `{"forward": "exit"}`

type Path struct {
	Direction []string
	Success   bool
}

func main() {
	wg := sync.WaitGroup{}
	Ch := make(chan Path)
	// Dn := make(chan string)
	var Res []Path
	I := make(map[string]interface{})
	err := json.Unmarshal([]byte(INPUT), &I)
	if err != nil {
		log.Printf("%v", err)
	}

	var Tr []string
	wg.Add(len(I))
	for i, v := range I {
		go processPaths(Tr, i, v, &wg, Ch)
	}

	SuccessPath := receivePaths(&Res, Ch)
	time.Sleep(time.Millisecond * 10)

	fmt.Println(SuccessPath)
	wg.Wait()
}

func receivePaths(Res *[]Path, Ch chan Path) []string {
	var SuccessPath [][]string
	for {
		select {
		case Dval := <-Ch:
			if Dval.Success {
				*Res = append(*Res, Dval)
				SuccessPath = append(SuccessPath, Dval.Direction)
			}
		case <-time.After(time.Millisecond * 10):
			if len(*Res) < 1 {
				// fmt.Println("[Sorry]")
				Rez := []string{"Sorry"}
				return Rez
			} else {
				//find shortest
				var ShortestIndex int
				for i, v := range SuccessPath {
					if i == 0 {
						ShortestIndex = 0
					}
					if len(v) < len(SuccessPath[ShortestIndex]) {
						ShortestIndex = i
					}
				}
				return SuccessPath[ShortestIndex]
			}
		}
	}
}

func processPaths(Tr []string, I string, V interface{}, wg *sync.WaitGroup, Ch chan Path) {
	Tr = append(Tr, I)

	switch V.(type) {
	case string:
		if V.(string) == "exit" {
			Ch <- Path{Direction: Tr, Success: true}
			wg.Done()
			return
		} else {
			Ch <- Path{Direction: Tr, Success: false}
			wg.Done()
			return
		}
	case map[string]interface{}:
		// Tr = append(Tr, I) //store the node name
		Nmap := V.(map[string]interface{})
		wg.Add(len(Nmap))
		for i, v := range Nmap {
			go processPaths(Tr, i, v, wg, Ch)
		}
	}
	wg.Done()
}
