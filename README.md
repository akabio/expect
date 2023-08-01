# expect

Go test expectations.

Intention of the library is to be able to write simple readable tests which produce easy to understand messages when the expectation fails.

```go
func TestExample(t *testing.T) {
    expect.Value(t, "the guy", "Peter").ToBe("Steven")
    // will fail with error:
    //   expected the guy to be 'Steven' but it is 'Peter'
}
```

## Expectations

### ToBe

Asserts that the value is deeply equal to the provided value.

```go
expect.Value(t, "the house", "big").ToBe("small")
// expected the house to be 'small' but it is 'big'
```

#### checking time.Time
When comparing times it can happen that two time instances with the exact same date/time can be different. This happens when one of them has location set to nil and the other to UTC. Although the documntation states that nil must be used instead of UTC some 3th party libs manage to return such instances.

#### checking nil

A check for nil will always be ok if the value is nil or a interface that points to a nil value.
This allows that following works:

    var a = *AType
    expect.Value(t, "a nil", a).ToBe(nil)

#### checking error
For error comparison the Error strings are returned. This can lead to messages like `expected Error to be 'foo' but it is 'foo'`.

#### checking structs, slices, maps
It will print complex data types formated as yaml:

```go
expect.Value(t, "array", []int{3, 1}).ToBe([]int{1, 3})
// expected array to be:
//   > - 1
//   > - 3
// but it is:
//   > - 3
//   > - 1
```

It will check for exact numbers:

```go
expect.Value(t, "liters", 3.4500000000001).ToBe(3.45)
// expected liters to be 3.45 but it is 3.4500000000001
```

### ToCount

Asserts that the list/map/chan/string has c elements.

```go
expect.Value(t, "token", "F7gTr7y").ToCount(8)
// expected token to have 8 elements but it has 7 elements
```

### ToHavePrefix/Suffix

Asserts that the string begins with the provided string or ends with it.

### NotToBe

Asserts that the value is not deeply equal to the provided value.

### ToBeAbout

Asserts that the number is about expected value with a margin of error of provided delta.

### ToBeType

Asserts that the type of the value is the same of the value given as parameter.

### ToBeSnapshot(filename)

ToBeSnapshot checks if the value is the same as what's in the given file.

- If the file isn't there, it will make a new one. You can look at it
  and change it if you need to.

- If the value doesn't match what's in the file, the test will fail.
  It will also create a new file with the same name but with a ".current"
  extension. This file will contain the failed content.