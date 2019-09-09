package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/leapfrogtechnology/shift/deployment/domain/project"
)

// Data is stored as shift.json.
type Data map[string]project.Response

var saveFilePath = "/var/lib/shift"

func failOnError(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

// Read parses data from shift.json
func read() Data {
	data, err := ioutil.ReadFile(saveFilePath + "/shift.json")
	failOnError(err, "Error reading file.")

	jsonData := Data{}

	json.Unmarshal(data, &jsonData)

	failOnError(err, "Error parsing json.")

	return jsonData
}

// Save persists project data in shift.json.
func Save(project project.Response) {
	jsonData := read()

	jsonData[project.ProjectName] = project

	data, _ := json.MarshalIndent(jsonData, "", " ")

	ioutil.WriteFile(saveFilePath+"/shift.json", data, 0644)
}
