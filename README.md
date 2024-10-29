# concurrency
Examples of golang concurrency use cases and sample implementations

## iobound 
folder `./iobound`
Using concurrency for io bound operation (file access)
The goroutines are blocked by the operating system file operations
While a goroutine is blocked then a goroutine that is in the runnable state and be run on the same CPU thread

## Data race atomic
folder `./datarace/atomic`
Data race occurs when multiple goroutines attempt to access the same variable and at least one goroutine is writing to the variable. 
In this example the code runs without error but the result may be wrong.
This is because increment  is three machine ops read the variable, increment it, write it back to memory. Another goroutine could access the variable while this is happening.  If the operation is being performed on a simple type eg int the we can change type to a specific type eg int64 and use atomic

## Data race mutex
folder `./datarace/mutex`
We can only use atomic for a few data types. 
This example has a map containing a single item that gets incremented
Without the mutex there is a 'fatal error: concurrent map writes' at runtime.
