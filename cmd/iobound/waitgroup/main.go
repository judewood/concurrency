package main

import (
	"fmt"
	"os"
	"sync"
)

func main() {
	const numberOfTasks = 100
	fileNames := createFiles(numberOfTasks)
	defer removeFiles(fileNames)

	var wg sync.WaitGroup
	for _, fileName := range fileNames {
		wg.Add(1)
		go func() {
			defer wg.Done()
			readFile(fileName)
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
