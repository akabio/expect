package expect

// Test implements gos testing.T methods used by expect.
// Used to also allow testing.B instances, as well as testing of expect.
type Test interface {
	Fatalf(f string, i ...interface{})
	Errorf(f string, i ...interface{})
	Error(p ...interface{})
}
