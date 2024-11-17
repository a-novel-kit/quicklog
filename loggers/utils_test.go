package loggers_test

type fakeAnimated struct {
	outTerm chan string
	outJSON chan map[string]interface{}
}

func (fake *fakeAnimated) RunTerminal(_ bool) <-chan string {
	return fake.outTerm
}

func (fake *fakeAnimated) RunJSON() <-chan map[string]interface{} {
	return fake.outJSON
}

func (fake *fakeAnimated) Close() {
	if fake.outTerm != nil {
		close(fake.outTerm)
	}
	if fake.outJSON != nil {
		close(fake.outJSON)
	}
}
