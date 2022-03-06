package main

import "aws_translator/internal/aws_translator/aws"

func main() {
	aws.InitAWS()
	aws.FindFrequency("data/it_in.txt", "data/it_out_clean.txt", 9, false)
	aws.TranslateFile("data/it_out_clean.txt", "data/it_out.txt", "ru")
}