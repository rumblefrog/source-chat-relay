package filter

import (
	"bufio"
	"os"
	"regexp"

	"github.com/rumblefrog/source-chat-relay/server/config"
	"github.com/sirupsen/logrus"
)

func ParseFilters() {
	if !config.Config.General.Filter {
		return
	}

	file, err := os.Open("filter.txt")

	defer file.Close()

	if err != nil {
		logrus.Warn("Unable to open filter file. Skipping.")

		return
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		regex, err := regexp.Compile(scanner.Text())

		if err != nil {
			continue
		}

		Filter = append(Filter, regex)
	}

	logrus.WithField("count", len(Filter)).Info("Compiled regular expressions")

	if err := scanner.Err(); err != nil {
		logrus.WithField("error", err).Warn("Unable to scan filter file")
	}
}
