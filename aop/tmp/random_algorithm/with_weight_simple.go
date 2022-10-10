package random_algorithm

import "math/rand"

func getWeightChoice(weights [][]uint32) uint32 {
	total := uint32(0)
	winner := 0
	for i, v := range weights {
		total += v[1]
		if uint32(rand.Float32()*float32(total)) < v[1] {
			winner = i
		}
	}
	return weights[winner][0]
}
