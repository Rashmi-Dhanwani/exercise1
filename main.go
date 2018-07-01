package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	csvFileName := flag.String("csv", "problem.csv", "enter file name")
	shuffle := flag.Bool("shuffle", true, "want to shuffle the question")
	timeD := flag.Int("timer", 30, "enter quiz duration in secs")
	flag.Parse()
	file, err := os.Open(*csvFileName)
	if err != nil {
		fmt.Println("unable to open")
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		fmt.Println("unable to read all")
	}
	if *shuffle {
		shuffleQuiz(lines)
	}
	playQuiz(lines, *timeD)

}

func shuffleQuiz(lines [][]string) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for n := range lines {
		r := r.Intn(n + 1)
		lines[n], lines[r] = lines[r], lines[n]
	}
}

func playQuiz(lines [][]string, timeD int) {
	var rightCount int
	var ans string
	ansChannel := make(chan string)
	timer := time.NewTimer(time.Duration(timeD) * time.Second)
	for k, v := range lines {
		fmt.Printf("Question %d: %s \n", k+1, v[0])
		go func() {
			fmt.Scanf("%s \n", &ans)
			ansChannel <- ans
		}()

		select {
		case ans := <-ansChannel:
			if strings.EqualFold(ans, strings.TrimSpace(v[1])) {
				rightCount++
			}
			if k+1 == len(lines) {
				showResult(lines, rightCount)
			}
			break
		case <-timer.C:
			showResult(lines, rightCount)
			return
		}
	}
}

func showResult(lines [][]string, rightCount int) {
	fmt.Println("Time Over")
	fmt.Printf("Total Ques= %d  Right Answers= %d Wrong Answers= %d \n", len(lines), rightCount, len(lines)-rightCount)
}
