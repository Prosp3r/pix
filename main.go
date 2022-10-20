package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
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

	for x := 0; x < len(I); x++ {
		Dval := <-Ch
		Res = append(Res, Dval)
	}

	if len(Res) > 0 {
		Countr := 0
		for _, v := range Res {
			if v.Success {
				fmt.Println(v.Direction)
				Countr++
			}
		}
		if Countr < 1 {
			fmt.Println("[Sorry]")
		}
	}

	wg.Wait()
}

func processPaths(Tr []string, I string, V interface{}, wg *sync.WaitGroup, Ch chan Path) {

	switch V.(type) {
	case string:
		if V == "exit" {
			Tr = append(Tr, I)
			Ch <- Path{Direction: Tr, Success: true}
			wg.Done()
			return
		} else {
			//return sorry wrong path/dead end
			Tr = append(Tr, "Sorry")
			Ch <- Path{Direction: Tr, Success: false}
			wg.Done()
			return
		}
	case map[string]interface{}:
		Tr = append(Tr, I) //store the node name

		Nmap := V.(map[string]interface{})
		for i, v := range Nmap {
			processPaths(Tr, i, v, wg, Ch)
		}
	}
}

//map[forward:tiger
//left:map[forward:map[upstairs:exit] left:dragon]
//right:map[forward:dead end]]
