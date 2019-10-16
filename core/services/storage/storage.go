package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/leapfrogtechnology/shift/core/structs"
)

func failOnError(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

// Read parses data from shift.json
func Read() structs.Project {
	currentDir, _ := os.Getwd()
	fileName := currentDir + "/shift.json"
	data, err := ioutil.ReadFile(fileName)

	failOnError(err, "Error reading file.")

	jsonData := structs.Project{}

	json.Unmarshal(data, &jsonData)

	failOnError(err, "Error parsing json.")

	return jsonData
}

// Save persists project data in shift.json.
func Save(project structs.Project) {
	jsonData, _ := json.MarshalIndent(project, " ", " ")

	currentDir, _ := os.Getwd()
	fileName := currentDir + "/shift.json"

	error := ioutil.WriteFile(fileName, jsonData, 0644)

	failOnError(error, "Could not save to file.")
}
