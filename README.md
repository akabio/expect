# expect

Go test expectations:

    func TestExample(t *testing.T) {
        expect.NamedValue(t, "the guy", "Peter").ToBe("Steven")
        // will fail with error:
        //   expected the guy to be 'Steven' but it is 'Peter'
    }
