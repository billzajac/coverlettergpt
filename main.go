package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"

	"github.com/ayush6624/go-chatgpt"
)

// use viper package to read .env file
// return the value of the key
func viperEnvVariable(key string) string {

	// SetConfigFile explicitly defines the path, name and extension of the config file.
	// Viper will use this and not check any of the config paths.
	// .env - It will search for the .env file in the current directory
	viper.SetConfigFile(".env")

	// Find and read the config file
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	// viper.Get() returns an empty interface{}
	// to get the underlying type of the key,
	// we have to do the type assertion, we know the underlying value is string
	// if we type assert to other type it will throw an error
	value, ok := viper.Get(key).(string)

	// If the type is a string then ok will be true
	// ok will make sure the program not break
	if !ok {
		log.Fatalf("Invalid type assertion")
	}

	return value
}

func main() {
	key := viperEnvVariable("OPENAI_KEY")
	resume := viperEnvVariable("RESUME")

	c, err := chatgpt.NewClient(key)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	// // Use default gpt3-turbo
	// res, err := c.SimpleSend(ctx, "Hey, Explain GoLang to me in 2 sentences.")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// a, _ := json.MarshalIndent(res, "", "  ")
	// log.Println(string(a))

	fmt.Printf("Job Title: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err = scanner.Err()
	if err != nil {
		log.Fatal(err)
	}
	job_title := scanner.Text()

	fmt.Printf("Company Name: ")
	scanner = bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err = scanner.Err()
	if err != nil {
		log.Fatal(err)
	}
	company_name := scanner.Text()

	fmt.Printf("Job Description: ")
	scanner = bufio.NewScanner(os.Stdin)
	last_line_was_blank := false
	blank_line_count := 0
	var job_description []string
	for {
		scanner.Scan()
		line := scanner.Text()
		if len(line) == 0 {
			if last_line_was_blank == true {
				blank_line_count = blank_line_count + 1
			}
			if blank_line_count > 3 {
				// Hopefully there aren't three consecutive line breaks in the JD - otherwise increase this
				break
			}
			last_line_was_blank = true
		} else {
			last_line_was_blank = false
			blank_line_count = 0
		}
		job_description = append(job_description, line)
	}

	err = scanner.Err()
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println("output:")
	//for _, l := range job_description {
	//	fmt.Println(l)
	//}

	content := fmt.Sprintf("Write me a personalized cover letter explaining why I'm a great candidate for this job. My resume is here: %s\n\nThe job title is %s, the company is %s, and here is the job description: %s", resume, job_title, company_name, strings.Join(job_description, " "))

	fmt.Println(content)
	fmt.Println("------------------------------------------------\n\n")
	// os.Exit(0)

	res, err := c.Send(ctx, &chatgpt.ChatCompletionRequest{
		Model: chatgpt.GPT4,
		Messages: []chatgpt.ChatMessage{
			{
				Role:    chatgpt.ChatGPTModelRoleSystem,
				Content: content,
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	a, _ := json.MarshalIndent(res, "", "  ")
	log.Println(string(a))

	//var response chatgpt.ChatResponse

	//err = json.Unmarshal(res, &response)
	//if err != nil {
	//		fmt.Println("error:", err)
	//}

	log.Printf("\n\n%s", res.Choices[0].Message.Content)
}
