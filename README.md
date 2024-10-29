# concurrency
Examples of golang concurrency use cases and sample implementations

## iobound 
Using concurrency for io bound operation (file access)
The goroutines are blocked by the operating system file operations
While a goroutine is blocked then a goroutine that is in the runnable state and be run on the same CPU thread

## Data race
Data race occurs when multiple goroutines access the same variable and at least one goroutine is writing to the variable. 
Increment  is three machine ops read the variable, increment it, write it back to memory. If the operation is being performed on a simple type eg int the we can change type to a specific type eg int64 and use atomic
 
