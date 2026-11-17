package tests

import (
	"log"

	"api-tests-template/internal/utils"
)

func SetupSuite() {
	log.Println("Init environment variables")
	utils.LoadEnv()
}

func TearDownSuite() {
	log.Println("Tear down suite")
}

func Precondition(text string) {
	utils.LogWithLabelAndTimeStamp("Precondition", text)
}

func Step(text string) {
	utils.LogWithLabelAndTimeStamp("Step", text)
}

func Check(text string) {
	utils.LogWithLabelAndTimeStamp("Check", text)
}
