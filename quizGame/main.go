package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
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
	file := readFile()
	result = fetchFileContent(file)
	fmt.Println("Total: ", result)

}

func readFile() *csv.Reader {
	folder, err := os.Open("quizGame.csv")
	if err != nil {
		log.Fatal(err)
	}
	file := csv.NewReader(folder)

	return file
}

func fetchFileContent(file *csv.Reader) response {
	var result response
	records, err := file.ReadAll()
	fmt.Println("record", records)
	if err != nil {
		log.Fatal(err)
	}
	result.total = len(records[1:])

	for key, value := range records {
		if key == 0 {
			continue
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		row := unit{question: value[0], answer: value[1]}

		//split the question by "+" sign
		nums := strings.Split(row.question, "+")

		ans, _ := strconv.Atoi(row.answer)
		num1, _ := strconv.Atoi(nums[0])
		num2, _ := strconv.Atoi(nums[1])
		if num1+num2 == ans {
			result.correct++
		} else {
			result.incorrect++
		}

	}
	return response{correct: result.correct, incorrect: result.incorrect, total: result.total}
}
