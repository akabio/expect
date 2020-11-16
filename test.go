package expect

// Test implements testing.T methods used by expect.
// Necessary to:
// - allow usage of testing.T and testing.B instances
// - for running tests
type Test interface {
	Fatalf(f string, i ...interface{})
	Errorf(f string, i ...interface{})
	Error(p ...interface{})
}
