package open

import (
	"bufio"
	"os"
	"strings"

	"github.com/Fabucik/wordlist-distribution/checkerr"
)

func GetWordlistContents(wordlistPath string) string {
	data, err := os.ReadFile(wordlistPath)
	checkerr.CheckError(err)

	return string(data)
}

func CountWordlistLines(wordlistContents string) int {
	return strings.Count(wordlistContents, "\n") + 1
}

func GetLineN(wordlist string, n int) string {
	reader := strings.NewReader(wordlist)
	scanner := bufio.NewScanner(reader)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		if lineNum == n {
			return "\n" + scanner.Text()
		}
	}

	return ""
}

func SplitNTimes(wordlistPath string, n int) []string {
	wordlistContents := GetWordlistContents(wordlistPath)
	numOfLines := CountWordlistLines(wordlistContents)
	whenToSplit := int(numOfLines / n)

	splittedWordlist := strings.SplitAfter(wordlistContents, GetLineN(wordlistContents, whenToSplit))

	n--

	for i := 1; i < n; {
		numOfLines = CountWordlistLines(splittedWordlist[i])
		whenToSplit = int(numOfLines / n)

		temporarySlice := strings.SplitAfter(splittedWordlist[i], GetLineN(splittedWordlist[i], whenToSplit))
		splittedWordlist = splittedWordlist[:len(splittedWordlist)-1]
		splittedWordlist = append(splittedWordlist, temporarySlice...)

		n--
	}

	return splittedWordlist
}
