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

//readFile : read and parse csv file and return a pointer to a csv.Reader struct
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

//fetchFileContent : read content of file
func fetchFileContent(file *csv.Reader) response {
	var result response
	var userAns int
	done := make(chan bool)

	//to get the total Number of question which is mentioned in the first row of the file
	data, err := file.Read()
	if err != nil {
		log.Fatal(err, " File is empty")
	}

	result.total, _ = strconv.Atoi(data[1])

	//to ignore the header
	_, err = file.Read()

	if err != nil {
		log.Fatal(err)
	}

	//to start the quiz, user has to press Enter key
	fmt.Println("Press Enter to start the quiz")
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

			row := unit{question: record[0], answer: record[1]}

			fmt.Println("Question: ", row.question)

			//write the answer of given question
			fmt.Scanln(&userAns)

			if err != nil {
				log.Fatal(err)
			}

			ans, _ := strconv.Atoi(row.answer)

			if userAns == ans {
				result.correct++
			}

		}

		//when all the questions are answered by user, send data to the channel
		done <- true
	}()

	//channel receive the signal
	select {
	case <-done:
		fmt.Println("You did it")
	case <-time.After(10 * time.Second):
		fmt.Println("TimeOut")
	}
	return result
}
