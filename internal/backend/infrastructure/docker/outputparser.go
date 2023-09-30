package docker

import (
	"bufio"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	doneSymbol = "DONE"
)

var (
	outputHeader = regexp.MustCompile(`#\d+ \[.*(?P<progress>\d/\d)] (?P<command>[A-Z]+) .*`)
	outputLine   = regexp.MustCompile(`#\d+ (?P<mark>\S+) (?P<output>.*)$`)
)

type outputParser struct{}

func (p outputParser) parseBuildOutputForRunTarget(output io.Reader) ([]byte, error) {
	scanner := bufio.NewScanner(output)
	for scanner.Scan() {
		submatch := outputHeader.FindStringSubmatch(scanner.Text())
		if len(submatch) != len(outputHeader.SubexpNames()) {
			// It is not a header since fully do not match regexp
			continue
		}

		instruction := submatch[2]
		if instruction != "RUN" { // Check only RUN instructions for output
			continue
		}

		progress := submatch[1]
		completed, err := completed(progress)
		if err != nil {
			return nil, err
		}

		if !completed { // Skip uncompleted stages
			continue
		}

		ok := scanner.Scan() // Skip line with intermediate container hash
		if !ok {
			return nil, errors.New("failed to skip line with intermediate container hash")
		}

		commandOutput, err := scanCommandOutput(scanner)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan command output")
		}

		return []byte(commandOutput), nil
	}

	return nil, errors.New("command output is not fully parsed")
}

func scanCommandOutput(scanner *bufio.Scanner) (string, error) {
	var res []string
	for scanner.Scan() {
		line := scanner.Text()

		submatch := outputLine.FindStringSubmatch(line)
		if len(submatch) != len(outputLine.SubexpNames()) {
			return "", errors.Errorf("invalid output line format: %s", line)
		}

		mark := submatch[1]
		output := submatch[2]

		if mark == doneSymbol {
			return strings.Join(res, "\n"), nil
		}

		res = append(res, output)
	}

	output := strings.Join(res, "\n")
	return "", errors.Errorf("output line is not terminated by %s: current line: %s", doneSymbol, output)
}

func completed(progress string) (bool, error) {
	const progressSeparator = "/"
	parts := strings.Split(progress, progressSeparator)
	if len(parts) != 2 {
		return false, errors.Errorf("incorrect progress format %s", progress)
	}

	readyPart, err := strconv.Atoi(parts[0])
	if err != nil {
		return false, errors.Errorf("incorrect progress format %s", progress)
	}

	allPart, err := strconv.Atoi(parts[1])
	if err != nil {
		return false, errors.Errorf("incorrect progress format %s", progress)
	}

	return readyPart == allPart, nil
}
