package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type sensorReading struct {
	RawValue int
	Value    int
}

func main() {
	const numReadings = 5

	inChannel := make(chan sensorReading, numReadings)
	readings := getRainfallSensorReadings(numReadings)
	sendReadingsToInputChannel(inChannel, readings)

	// limit the number of go routines to prevent excess context switching
	maxChannels := min(numReadings, runtime.NumCPU())

	// fan out the raw sensor readings to process them in parallel
	multipleChannels := fanOut(inChannel, maxChannels)

	//merge the processed data back into a single channel
	outChannel := fanIn(multipleChannels)

	// Print out the results
	for result := range outChannel {
		outputReading(result)
	}
}

// getRainfallSensorReadings mimics a stream of sensor readings being sent to the application
func getRainfallSensorReadings(numReadings int) []sensorReading {
	readings := make([]sensorReading, numReadings)
	for i := 0; i < numReadings; i++ {
		readings[i].RawValue = i
	}
	fmt.Printf("readings from sensors: %#v", readings)
	return readings
}

// sendReadingsToInputChannel sends a slice of sensor readings to an input channel
func sendReadingsToInputChannel(inputCh chan sensorReading, readings []sensorReading) {
	defer close(inputCh)
	for i, reading := range readings {
		fmt.Printf("\nPushing reading %d value: %#v to in channel", i, reading)
		inputCh <- reading
	}
}

// fanOut receives items on a single input channel and outputs them on to multiple output channels
// this allows part of the processing pipeline to be done in parallel
func fanOut(inChannel <-chan sensorReading, maxChannels int) []<-chan sensorReading {
	outChannels := make([]<-chan sensorReading, 0, maxChannels)

	for i := 0; i < maxChannels; i++ {
		ch := make(chan sensorReading)
		outChannels = append(outChannels, ch)

		go func() {
			defer close(ch)
			for val := range inChannel {
				fmt.Printf("\nSending %#v from in channel to out channel %d", val, i)
				ch <- val
			}
		}()
	}

	return outChannels
}

// fanIn receives values from  multiple input channels, 
// processes each value in its own go routine and then 
// sends the result to single output channel
func fanIn(inChannels []<-chan sensorReading) <-chan sensorReading {
	outChannel := make(chan sensorReading)

	wg := &sync.WaitGroup{} 
	for _, ch := range inChannels {
		wg.Add(1)
		go func(c <-chan sensorReading) {
			defer wg.Done()
			for val := range c {
				processReading(&val)
				outChannel <- val
			}
		}(ch)
	}
	go func() {
		defer close(outChannel)
		wg.Wait()
	}()
	return outChannel
}

// processReading simulates some slower processing
// usings a time delay and a multiplication
func processReading(reading *sensorReading) {
	time.Sleep(200 * time.Millisecond)
	reading.Value = reading.RawValue * 2 //so we can see the processed result
}

// outputReading prints out a processed reading
// In the real world it could be doing som processing that  
// needs to be done sequentially
func outputReading(reading sensorReading) {
	fmt.Printf("\nOutputting reading: %#v", reading)
}
