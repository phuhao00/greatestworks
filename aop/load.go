package aop

import "time"

// Interval at which weavelets report their load.
//
// TODO(mwhittaker): Have this be an option that is sent to the weavelet over
// the pipe? Or maybe switch the assigner to poll for load?
const LoadReportInterval = 5 * time.Minute
