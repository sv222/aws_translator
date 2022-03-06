package aws

import (
	"bufio"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/translate"
	"github.com/joho/godotenv"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

var (
	tr           *translate.Translate
	translations []*Text
)

type Text struct {
	Key string
	En  string
	Ru  string
}

func InitAWS() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading environment variables: %v", err)
	}

	keyAWS := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAWS := os.Getenv("AWS_SECRET_ACCESS_KEY")

	creds := credentials.NewStaticCredentials(keyAWS, secretAWS, "")

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: creds,
	})

	if err != nil {
		log.Fatal(err)
	}

	tr = translate.New(sess)

}

func TranslateFile(in, out, lang string) {
	inFile, err := os.Open(in)
	if err != nil {
		log.Fatal(err)
	}
	defer inFile.Close()

	outFile, err := os.OpenFile(out, os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	reader := bufio.NewReader(inFile)

	for {
		line, _, err := reader.ReadLine()

		if err == io.EOF {
			break
		}

		l := string(line)

		if len(l) > 1 {

			translated := Translate(l, lang)
			res := l + " -> " + translated + "\n"
			outFile.WriteString(res)
			fmt.Println(res)
		}

	}
}

func Translate(key, lang string) string {
	input := &translate.TextInput{
		SourceLanguageCode: aws.String("en"),
		TargetLanguageCode: aws.String(lang),
		Text:               aws.String(key),
	}

	req, out := tr.TextRequest(input)

	if err := req.Send(); err != nil {
		panic(req.Error)
	}

	return *out.TranslatedText
}

func FindFrequency(in, out string, num int, showFreq bool) {
	unique := make(map[string]int)

	inFile, err := os.Open(in)
	if err != nil {
		log.Fatal(err)
	}
	defer inFile.Close()

	outFile, err := os.OpenFile(out, os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	reader := bufio.NewReader(inFile)

	for {
		line, _, err := reader.ReadLine()

		if err == io.EOF {
			break
		}

		l := strings.ToLower(string(line))

		if len(l) > 1 {

			if _, ok := unique[l]; !ok {
				unique[l] = 1
			} else {
				unique[l]++
			}
		}
	}

	type kv struct {
		Key   string
		Value int
	}

	resData := make([]kv, 0)

	for k, v := range unique {
		resData = append(resData, kv{k, v})
	}

	sort.Slice(resData, func(i, j int) bool {
		return resData[i].Value > resData[j].Value
	})

	for _, kv := range resData {
		if kv.Value > num {
			match, _ := regexp.MatchString("^[a-zA]+$", kv.Key)
			if match && showFreq {
				fmt.Fprintf(outFile, "%s -> %d\n", kv.Key, kv.Value)
			}
			if match && showFreq == false {
				fmt.Fprintf(outFile, "%s\n", kv.Key)
			}
		}
	}
}