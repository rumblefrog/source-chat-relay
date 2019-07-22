package filter

import (
	"bufio"
	"os"
	"regexp"

	"github.com/rumblefrog/source-chat-relay/server/config"
	"github.com/sirupsen/logrus"
)

func init() {
	if !config.Conf.General.Filter {
		return
	}

	file, err := os.Open("filter.txt")

	if err != nil {
		logrus.Warn("Unable to open filter file. Skipping.")
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
		logrus.WithField("error", err).Warn("Unable to scan filter file")
	}
}
