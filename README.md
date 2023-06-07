# Build a cover letter with ChatGPT

This is a super simple CLI for writing cover letters for job applications

## USAGE

* Create a .env file that looks like this

```
OPENAI_KEY=<YOUR_KEY>
RESUME="
A TEXT VERSION OF YOUR RESUME
Just paste it into text and see how it looks. Then correct any formatting necessary.
This will obviously be a multi-line field in this file, so be sure to have that final quote.
"
```

## BUILD / RUN

```
go mod tidy # get all of those tasty deps
go build main.go
./main
```

* Press return a bunch of times at the end of pasting the job description
* Wait patiently
* Paste the output text into a new doc
* WARNING: Carefully read the whole thing (chatgpt will make mistakes and sometimes add completely fictional details)

## TODO

* Build a web interface
    * Probably will want auth with it too so that if I make it public, I can manage costs
