package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) Out

func runStage(done In, input In, stage Stage) Out {
	output := make(Bi)
	stageOutput := stage(input)

	go func() {
		defer close(output)
		for {
			select {
			case <-done:
				return
			case val, ok := <-stageOutput:
				if !ok {
					return
				}
				output <- val
			}
		}
	}()
	return output
}

func ExecutePipeline(input In, done In, stages ...Stage) Out {
	current := input
	for _, stage := range stages {
		current = runStage(done, current, stage)
	}
	return current
}
