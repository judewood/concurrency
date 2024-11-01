package main

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type cup struct {
	ID    int
	value string
}

// coin represents a coin to be placed under a cup
const coin = "coin"

func main() {
	totalCups := 10
	cups := setupCups(make([]cup, totalCups))
	guess := getGuess(len(cups))
	winningCup := liftCups(cups)
	playerFeedback(guess, winningCup)
}

// setupCups creates a slice of cups and puts a coin under one of them
func setupCups(cups []cup) []cup {
	insertPoint := rand.Intn(len(cups))
	for i := 0; i < len(cups); i++ {
		cups[i] = cup{ID: i + 1}
		if i == insertPoint {
			cups[i].value = coin
		}
	}
	return cups
}

// getGuess is an infinite loop that breaks out once the play has entered
// an integer in range of the number of cups
func getGuess(totalCups int) int {
	for {
		fmt.Printf("\nFind the cup with the coin under it. Guess between 1 and %d\n", totalCups)
		guess, err := getIntFromStdIn()
		if err == nil && guess >= 1 && guess <= totalCups {
			return guess
		}
	}
}

// getIntFromStdIn gets the player input and tries to convert it to an int
func getIntFromStdIn() (int, error) {
	var input string
	_, err := fmt.Scanln(&input)
	checkErr(err)
	fmt.Print(input)
	return strconv.Atoi(input)
}

// liftCups lifts cups until the coin is found
func liftCups(cups []cup) int {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc() //prevent possible context leak
	var wg sync.WaitGroup
	winner := 0
Loop:
	for _, liftedCup := range cups {
		time.Sleep(500 * time.Millisecond) //simulate delay between lifting cups so we can see some go routines not being started
		select {
		case <-ctx.Done():
			break Loop //break from loop not just the select
		default:
			wg.Add(1)
			go func() {
				defer wg.Done()
				winningCup := checkCup(liftedCup, cancelFunc)
				if winningCup {
					winner = liftedCup.ID
				}
			}()
		}
	}
	wg.Wait() //wait until all running goroutines are done
	return winner
}

// checkCup lifts a single cup and returns whether there was a coin under it
func checkCup(liftedCup cup, cancelFunc context.CancelFunc) bool {
	hasCoin := false
	fmt.Printf("\nLifting cup %d...", liftedCup.ID)
	if liftedCup.value == coin {
		hasCoin = true
		cancelFunc() //signal to stop lifting any more cups
	}
	return hasCoin
}

// playerFeedback tell the player what cup the coin was under and whether they won
func playerFeedback(guess int, winningCup int) {
	fmt.Printf("\nThe coin was under cup %d", winningCup)
	if guess == winningCup {
		fmt.Printf("\nYou win !")
		return
	}
	fmt.Println("\nBetter luck next time")
}

// checkErr handles errors that should cause panic
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
