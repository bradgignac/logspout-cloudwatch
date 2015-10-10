package cloudwatch

func filter(in <-chan Log) <-chan Log {
	out := make(chan Log)

	go func() {
		defer close(out)

		for log := range in {
			if log.Body() != "" {
				out <- log
			}
		}
	}()

	return out
}
