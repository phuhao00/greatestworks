package metrics_test

import (
	"greatestworks/aop/metrics"
	"greatestworks/aop/protos"
)

var (
	// Unlabeled counter.
	catCounter = metrics.Register(
		protos.MetricType_COUNTER,
		"example_cats",
		"Number of cats.",
		nil,
	)

	// Labeled counters.
	dogCounters = metrics.RegisterMap[dogLabels](
		protos.MetricType_COUNTER,
		"example_dogs",
		"Number of dogs, by breed.",
		nil,
	)
	corgiCounter     = dogCounters.Get(dogLabels{"corgi"})
	poodleCounter    = dogCounters.Get(dogLabels{"poodle"})
	dachshundCounter = dogCounters.Get(dogLabels{"dachshund"})
	dalmatianCounter = dogCounters.Get(dogLabels{"dalmatians"})
)

type dogLabels struct {
	Breed string
}

func Example() {
	catCounter.Add(9.0)
	corgiCounter.Add(2.0)
	poodleCounter.Add(1.0)
	dachshundCounter.Add(10.0)
	dalmatianCounter.Add(101.0)
}
