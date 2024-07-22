# A Simple Go Promise Library
![Coverage](https://raw.githubusercontent.com/wilfredtan/promise/actions-badge/badge.svg)
This library is a wrapper for writing Go routine functions. It is inspired by Javascript's "Promise" feature.

Example code using the library:
```go
func MyFunction() (string, error) {
	// Create a new promise with resolver callback containing logic
	// Panic recovery is automatically handled :D
	promiseMyValue := promise.New(func(
		resolve func(string),
		reject func(error)
	) {
		// ... some code logic that defines a string `someValue` and an error `someErr`
		// Error handling in callback
		if someErr != nil {
			reject(someErr)
			return
		}
		// Resolve value
		resolve(someValue)
	})
	// Await the data
	myValue, err := promiseMyValue.Await();
	// Error handling
	if err != nil {
		return nil, err
	}
	return myValue, nil
}
```

Before this, whenever I needed to run multiple functions concurrently, I find myself repeating a lot of code with lots of variables to name and keep track of. The same code before I wrote this library:

```go
func MyFunction() (string, error) {
	// Data to be captured from Go routine function
	type myData struct {
		value string
		err   error
	}
	// A channel for receiving the data
	chMyData := make(chan *myData)
	// Go routine function with logic
	go func(ch chan *myData) {
		data := myData{}
		// Panic recovery
		defer func() {
			if err := recover(); err != nil {
				data.err, _ = err.(error)
				ch <- &data
			}
		}()
		// ... some code logic that defines `someValue` and `someErr`
		// Error handling in go routine
		if someErr != nil {
			data.err = someErr
			ch <- &data
			return
		}
		// Send the data through the channel
		data.value = someValue
		ch <- &data
	}(chMyData)
	// Receive the data from the channel
	data := <-chMyData
	// Error handling
	if data.err != nil {
		return "", data.err
	}
	return data.value, nil
}
```
