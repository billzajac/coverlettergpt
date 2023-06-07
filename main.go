package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
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
	var debug bool
	// Used to indicate the job description has been completely pasted
	// Assumption: there aren't newline_max line breaks in the JD
	var newline_max int
	var help bool

	flag.IntVar(&newline_max, "n", 3, "Newlines to indicate the end of the JD")
	flag.BoolVar(&debug, "d", false, "DEBUG - Only show the prompt, do not submit it")
	flag.BoolVar(&help, "h", false, "Help")
	flag.Parse()

	if help {
		flag.PrintDefaults()
		fmt.Println("\nAlso be sure to create a .env file with\n\nOPENAI_KEY=<YOUR_KEY>\nRESUME=\"<MULTI-LINE TEXT RESUME\" (don't forget the end quote)\n")
		os.Exit(0)
	}

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
	newline_count := 0
	var job_description []string
	for {
		scanner.Scan()
		line := scanner.Text()
		if len(line) == 0 {
			newline_count = newline_count + 1
			if newline_count >= newline_max {
				fmt.Println("\nSUBMITTING THE FOLLOWING PROMPT ==---->\n")
				break
			} else if newline_max-newline_count == 1 {
				fmt.Printf("Press RETURN %d last time to submit the prompt...", newline_max-newline_count)
			} else {
				fmt.Printf("Press RETURN %d more times to submit the prompt...", newline_max-newline_count)
			}
		} else {
			// reset the counter because we found more non-empty lines
			newline_count = 0
		}
		job_description = append(job_description, line)
	}

	err = scanner.Err()
	if err != nil {
		log.Fatal(err)
	}

	content := fmt.Sprintf("Write me a personalized cover letter explaining why I'm a great candidate for this job. My resume is here: %s\n\nThe job title is %s, the company is %s, and here is the job description: %s", resume, job_title, company_name, strings.Join(job_description, " "))

	fmt.Println(content)
	fmt.Println("------------------------------------------------\n\n")
	if debug {
		os.Exit(0)
	}

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
