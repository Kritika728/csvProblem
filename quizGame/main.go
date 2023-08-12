package main

import (
	"bufio"
	"encoding/csv"
	"flag"
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

var filename = flag.String("quizGame.csv", "problem.csv", "csv file that contains quiz questions")
var timeLimit = flag.Int("15", 30, "Time limit for the quiz")

func main() {
	var result response

	//file := readFile()
	//result.total = getTotalQuestion(file)

	file := readFile()
	result = fetchFileContent(file)

	//output
	fmt.Printf("%v/%v ", result.correct, result.total)

}

//readFile : read and parse csv file and return a pointer to a csv.Reader struct
func readFile() *csv.Reader {
	folder, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	file := csv.NewReader(folder)

	return file
}

//fetchFileContent : read content of file
func fetchFileContent(file *csv.Reader) response {
	var result response
	var userAns int
	scanner := bufio.NewScanner(os.Stdin)
	done := make(chan bool)

	//to get the total Number of question which is mentioned in the first row of the file
	//return err if file is empty
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

			for {
				userAns, err = readAns(scanner)
				if err == nil {
					break
				}
				fmt.Println("Please renter ")
			}

			ans, _ := strconv.Atoi(row.answer)

			if userAns == ans {
				result.correct++
			}

		}

		//when all the questions are answered by user, send data to the channel
		done <- true
	}()

	//channel receive the data
	select {
	case <-done:
		fmt.Println("You did it")
	case <-time.After(time.Duration(*timeLimit) * time.Second):
		fmt.Println("TimeOut")
	}
	return result
}

//getTotalQuestion : get the total number of question
func getTotalQuestion(file *csv.Reader) int {
	var count response
	records, err := file.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	count.total = len(records[1:])

	return count.total
}

func readAns(scanner *bufio.Scanner) (int, error) {
	var (
		err error
		num int
	)
	scanner.Scan()
	input := scanner.Text()
	if num, err = strconv.Atoi(input); err != nil {
		fmt.Println("Invalid input. Please enter a valid integer: ", err)
	}
	return num, err
}
