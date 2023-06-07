package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	gosdk "github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/tidwall/gjson"
)

func main() {
	auth := gosdk.Auth{
		UserName: "goadmin",
		Password: "goadmin",
	}
	client := gosdk.NewClient("https://gocd.example.com/go", auth, "info", []byte{})
	//fmt.Println(client.GetServerHealth())

	env, _ := client.GetEnvironment("dev")
	pipelines := env.Pipelines
	file, err := os.Create("pipelineInfo.csv")
	if err != nil {
		fmt.Println(err)
	}
	writer := csv.NewWriter(file)
	for _, pipeline := range pipelines {
		config, err := client.GetPipelineConfig(pipeline.Name)
		if err != nil {
			fmt.Println(err)
		}
		bytes, _ := json.Marshal(config.Config)
		pipelineConfig := string(bytes)
		if gjson.Valid(pipelineConfig) {
			for _, stage := range gjson.Get(pipelineConfig, "stages").Array() {
				stageName := gjson.Get(stage.String(), "name")
				for _, job := range gjson.Get(stage.String(), "jobs").Array() {
					jobName := gjson.Get(job.String(), "name")

					//fmt.Println(pipeline.Name, stageName, jobName)
					resources := gjson.Get(job.String(), "resources.#")

					if resources.Uint() != 0 {
						err := writer.Write([]string{pipeline.Name, stageName.String(), jobName.String()})
						if err != nil {
							fmt.Println(err)
						}
					}

				}

			}
		}
	}

	writer.Flush()
}
