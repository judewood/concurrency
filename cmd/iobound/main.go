package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
)

func main() {
	const numberOfTasks = 100
	fileNames := createFiles(numberOfTasks)
	defer removeFiles(fileNames)

	maxGoroutines := min(numberOfTasks, runtime.NumCPU()) 
	// Empty struct because this is a signalling channel and empty struct uses 0 bytes of memory
	limiterChan := make(chan struct{}, maxGoroutines)

	var wg sync.WaitGroup
	for _, fileName := range fileNames {
		wg.Add(1)
		limiterChan <- struct{}{} //block when limit maxGoroutines is reached
		go func() {
			readFile(fileName)
			<-limiterChan //receive from channel allow another goroutine to be run
			wg.Done()
		}()
	}
	wg.Wait() //prevent removal of files until they have been read
}

func createFiles(numTasks int) []string {
	fileNames := make([]string, 0, numTasks)
	for i := 0; i < numTasks; i++ {
		i := i
		filename := fmt.Sprintf("temp  %d", i)
		f, err := os.Create(filename)
		checkErr(err)
		defer f.Close()
		bytes := []byte(fmt.Sprintf("file number %d", i))
		_, err = f.Write(bytes)
		checkErr(err)
		fileNames = append(fileNames, filename)
	}
	return fileNames
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func readFile(fileName string) {
	dat, err := os.ReadFile(fileName)
	checkErr(err)
	fmt.Println(string(dat))
}

func removeFiles(fileNames []string) {
	for _, filename := range fileNames {
		err := os.Remove(filename)
		checkErr(err)
	}
}
