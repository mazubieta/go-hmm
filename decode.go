package main

import (
	/*"bytes"*/
	"encoding/csv"
	"flag"
	"fmt"
	"git.bigodev.com/mazubieta/go-hmm/markov"
	"os"
	"strconv"
)

func main() {
	var inputFile, outputFile string

	// set flags for simulation
	flag.StringVar(&inputFile, "i", "output_hmm", "Output file name(Default: \"output_hmm\")")
	flag.StringVar(&outputFile, "o", "output_hmm", "Output file name(Default: \"output_hmm\")")
	flag.Parse()

	// create fair die
	fair := markov.CreateState("f")
	fair.AddEmition("1", 1.0)
	fair.AddEmition("2", 1.0)
	fair.AddEmition("3", 1.0)
	fair.AddEmition("4", 1.0)
	fair.AddEmition("5", 1.0)
	fair.AddEmition("6", 1.0)
	fair.AddNeighbor("f", 0.5)
	fair.AddNeighbor("u", 0.5)

	// craete unfair die
	unfair := markov.CreateState("u")
	unfair.AddEmition("1", 1.0)
	unfair.AddEmition("2", 1.0)
	unfair.AddEmition("3", 1.0)
	unfair.AddEmition("4", 1.0)
	unfair.AddEmition("5", 1.0)
	unfair.AddEmition("6", 5.0)
	unfair.AddNeighbor("f", 0.5)
	unfair.AddNeighbor("u", 0.5)

	// lookup map of states
	viterbi := markov.Viterbi{
		markov.Model{States: []*markov.State{&fair, &unfair}},
	}
	posterior := markov.Posterior{
		markov.Model{States: []*markov.State{&fair, &unfair}},
	}

	viterbi.Initialize()
	posterior.Initialize()

	// read csv file of all tests
	file, err := os.Open(inputFile)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	nRecords := len(lines)
	seeds := make([]int64, nRecords)
	tprobs := make([]float64, nRecords)
	observations := make([][]string, nRecords)
	hiddenStates := make([][]string, nRecords)

	var numObservations int64
	for i, line := range lines {
		seeds[i], _ = strconv.ParseInt(line[0], 10, 64)
		tprobs[i], _ = strconv.ParseFloat(line[1], 64)
		numObservations, _ = strconv.ParseInt(line[2], 10, 0)
		observations[i] = line[3 : 3+numObservations]
		hiddenStates[i] = line[3+numObservations : 3+2*numObservations]
	}

	n := int(numObservations)

	for i := 0; i < nRecords; i++ {
		observation := observations[i]
		states := hiddenStates[i]
		fair.EditNeighbor("f", 1.0-tprobs[i])
		fair.EditNeighbor("u", tprobs[i])
		unfair.EditNeighbor("f", tprobs[i])
		unfair.EditNeighbor("u", 1.0-tprobs[i])

		if len(observation) != len(states) {
			fmt.Println("Mismatch record [", i, "] ...skipping...")
			continue
		}

		// decode via Viterbi
		pi_v := viterbi.Decode(observation)
		pi_p := posterior.Decode(observation)

		/*var stateBuffer, piBuffer bytes.Buffer*/
		correct_v := 0.0
		correct_p := 0.0
		for j := 0; j < n; j++ {
			/*piBuffer.WriteString(pi[j])*/
			/*stateBuffer.WriteString(states[j])*/
			if pi_v[j] == states[j] {
				correct_v++
			}
			if pi_p[j] == states[j] {
				correct_p++
			}
		}
		/*fmt.Println(stateBuffer.String())*/
		/*fmt.Println(piBuffer.String())*/
		fmt.Println(tprobs[i], n, correct_v, correct_v/float64(n), "||", correct_p, correct_p/float64(n))

	}
}
