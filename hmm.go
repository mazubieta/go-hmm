package main

import (
	"bytes"
	"flag"
	"git.bigodev.com/mazubieta/go-hmm/markov"
	"io/ioutil"
	"math/Rand"
	"strconv"
	"time"
)

func main() {
	var seed int64
	var numIter int
	var transProb float64
	var outputFile string

	// set flags for simulation
	flag.Int64Var(&seed, "s", time.Now().Unix(), "Seed for PRNG (Default: time.Time())")
	flag.IntVar(&numIter, "n", 1000, "Number of Time Steps (Default: 1000)")
	flag.Float64Var(&transProb, "t", 0.1, "Transition Porbability (Default: 0.1)")
	flag.StringVar(&outputFile, "o", "output_hmm", "Output file name(Default: \"output_hmm\")")
	flag.Parse()

	// seed PRNG
	rand.Seed(seed)

	// create fair die
	fair := markov.CreateState("f")
	fair.AddEmition("1", 1.0)
	fair.AddEmition("2", 1.0)
	fair.AddEmition("3", 1.0)
	fair.AddEmition("4", 1.0)
	fair.AddEmition("5", 1.0)
	fair.AddEmition("6", 1.0)
	fair.AddNeighbor("f", 1.0-transProb)
	fair.AddNeighbor("u", transProb)

	// craete unfair die
	unfair := markov.CreateState("u")
	unfair.AddEmition("1", 1.0)
	unfair.AddEmition("2", 1.0)
	unfair.AddEmition("3", 1.0)
	unfair.AddEmition("4", 1.0)
	unfair.AddEmition("5", 1.0)
	unfair.AddEmition("6", 2.0)
	unfair.AddNeighbor("f", transProb)
	unfair.AddNeighbor("u", 1.0-transProb)

	// lookup map of states
	states := map[string]*markov.State{"f": &fair, "u": &unfair}

	// randomly start at a state (assumed equal probability of start state)
	ridx := rand.Intn(len(states))
	i := 0

	var currentState *markov.State
	for _, state := range states {
		if ridx == i {
			currentState = state
			break
		}
		i++
	}

	var emitionBuffer, stateBuffer, fileBuffer bytes.Buffer

	for i := 0; i < numIter; i++ {
		e, _ := currentState.GetRandomEmition()
		emitionBuffer.WriteString(e)
		stateBuffer.WriteString(currentState.Name)
		nextState, _ := currentState.GetRandomTransition()
		currentState = states[nextState]
	}

	fileBuffer.WriteString(strconv.FormatInt(seed, 10))
	fileBuffer.WriteString(",")
	fileBuffer.WriteString(strconv.FormatFloat(transProb, 'f', 6, 64))
	fileBuffer.WriteString(",")
	fileBuffer.WriteString(emitionBuffer.String())
	fileBuffer.WriteString(",")
	fileBuffer.WriteString(stateBuffer.String())
	fileBuffer.WriteString("\n")

	ioutil.WriteFile(outputFile, fileBuffer.Bytes(), 0644)
}
