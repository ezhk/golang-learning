package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)
type I interface{}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	outCh := make(Bi)

	// border case when in channel is nil
	if in == nil {
		close(outCh)
		return outCh
	}

	queueChannels := make([]Bi, 0)
	for v := range in {
		stagesOut := make(Bi)
		queueChannels = append(queueChannels, stagesOut)

		go processStages(v, done, stagesOut, stages...)
	}

	// process answers by order
	go func(out Bi) {
		defer close(out)

		for _, ch := range queueChannels {
			select {
			case <-done:
				return
			case v, ok := <-ch:
				if ok {
					out <- v
				}
			}
		}
	}(outCh)

	return outCh
}

func processStages(value I, done In, out chan<- interface{}, stages ...Stage) {
	defer close(out)

	for _, stage := range stages {
		bufferedInterCh := make(chan interface{})

		// goroutine put value to channel for stage
		go func() {
			defer close(bufferedInterCh)

			select {
			case <-done:
				return
			case bufferedInterCh <- value:
			}
		}()

		// read stage result
		select {
		case <-done:
			return
		case value = <-stage(bufferedInterCh):
		}
	}

	select {
	case <-done:
		return
	case out <- value:
	}
}
