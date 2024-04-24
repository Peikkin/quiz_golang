package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Examples struct {
	question string
	answer   string
}

func RunningExamples(fileName string) ([]Examples, error) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Error().Err(err).Msg("Ошибка открытия файла")
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvLines, err := csvReader.ReadAll()
	if err != nil {
		log.Error().Err(err).Msg("Ошибка чтения файла")
		return nil, err
	}
	return ParseExample(csvLines), nil
}

func ParseExample(lines [][]string) []Examples {
	exam := make([]Examples, len(lines))
	for i := 0; i < len(lines); i++ {
		exam[i] = Examples{question: lines[i][0], answer: lines[i][1]}
	}
	return exam
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	fileName := flag.String("file", "quiz.csv", "Путь до csv файла")

	timer := flag.Int("t", 30, "Таймер викторины")

	flag.Parse()

	exam, err := RunningExamples(*fileName)
	if err != nil {
		log.Panic().Err(err).Msg("Ошибка чтения файла")
	}
	correctAnswers := 0

	tQuiz := time.NewTimer(time.Duration(*timer) * time.Second)

	answerChanel := make(chan string)

exampleLoop:
	for i, p := range exam {

		var answer string
		fmt.Printf("Вопрос %v: %v=", i+1, p.question)
		go func() {
			fmt.Scanf("%v", &answer)
			answerChanel <- answer
		}()
		select {
		case <-tQuiz.C:
			log.Info().Msg("Время истекло")
			break exampleLoop
		case iAns := <-answerChanel:
			if iAns == p.answer {
				correctAnswers += 1
			}
			if i == len(exam)-1 {
				close(answerChanel)
			}
		}
	}
	log.Info().Msgf("Ваш результат %v из %v", correctAnswers, len(exam))
	fmt.Print("Нажмите Enter для выхода")
	<-answerChanel
}
