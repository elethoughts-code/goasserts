package assertion

type FsExpectation interface {
	FileExists()
}

func (exp *expectation) FileExists() {
	exp.t.Helper()
	exp.Matches(FileExists())
}
