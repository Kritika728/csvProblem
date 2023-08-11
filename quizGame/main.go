package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

type unit struct {
	question string
	answer   string
}
type response struct {
	total     int
	correct   int
	incorrect int
}

func main() {
	var result response
	//file := readFile()
	//result.total = getTotalQuestion(file)
	samefile := readFile()
	result = fetchFileContent(samefile)

	fmt.Printf("%v/%v ", result.correct, result.total)

}

func readFile() *csv.Reader {
	folder, err := os.Open("quizGame.csv")
	if err != nil {
		log.Fatal(err)
	}
	file := csv.NewReader(folder)

	return file
}

func getTotalQuestion(file *csv.Reader) int {
	var count response
	records, err := file.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	count.total = len(records[1:])

	return count.total
}
func fetchFileContent(file *csv.Reader) response {
	var result response
	var userAns int
	done := make(chan bool)

	data, err := file.Read()
	result.total, _ = strconv.Atoi(data[1])

	//to ignore the header
	_, err = file.Read()

	if err != nil {
		log.Fatal(err)
	}
	//to start the quiz, user has to press any key
	fmt.Println("Enter any key to start the quiz")
	fmt.Scanln()

	go func() {
		for {
			record, err := file.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			//timer
			row := unit{question: record[0], answer: record[1]}

			fmt.Println("Question: ", row.question)

			//give the answer of given question
			fmt.Scanln(&userAns)

			if err != nil {
				log.Fatal(err)
			}

			ans, _ := strconv.Atoi(row.answer)

			if userAns == ans {
				result.correct++
			}

		}
		done <- true
	}()

	//timer will start
	select {
	case <-done:
		fmt.Println("You did it")
	case <-time.After(10 * time.Second):
		fmt.Println("TimeOut")
	}
	return result
}
