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

### ToHavePrefix

Asserts that the string begins with the provided string.

### NotToBe

Asserts that the value is not deeply equal to the provided value.

### ToBeAbout

Asserts that the number is about expected value with a margib of error of provided delta.
