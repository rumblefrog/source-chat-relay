package filter

import (
	"bufio"
	"os"
	"regexp"

	log "github.com/sirupsen/logrus"
)

func init() {
	file, err := os.Open("filter.txt")

	if err != nil {
		log.Warn("Unable to open filter file. Skipping.")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		regex, err := regexp.Compile(scanner.Text())

		if err != nil {
			continue
		}

		Filter = append(Filter, regex)
	}

	if err := scanner.Err(); err != nil {
		log.WithField("error", err).Warn("Unable to scan filter file")
	}
}
