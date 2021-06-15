package filter

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/cloudskiff/driftctl/pkg/resource"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
	"github.com/sirupsen/logrus"
)

type DriftIgnore struct {
	matcher gitignore.Matcher
}

func NewDriftIgnore() *DriftIgnore {
	d := DriftIgnore{
		gitignore.NewMatcher(nil),
	}
	err := d.readIgnoreFile()
	if err != nil {
		logrus.Debug(err)
	}
	return &d
}

func (r *DriftIgnore) readIgnoreFile() error {
	file, err := os.Open(".driftignore")
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []gitignore.Pattern
	scanner := bufio.NewScanner(file)
	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		line := scanner.Text()
		lines = append(lines, gitignore.ParsePattern(line, nil))
		if !strings.HasSuffix(line, "*") {
			line := fmt.Sprintf("%s.*", line)
			lines = append(lines, gitignore.ParsePattern(line, nil))
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	r.matcher = gitignore.NewMatcher(lines)

	return nil
}

func (r *DriftIgnore) IsResourceIgnored(res resource.Resource) bool {
	strRes := fmt.Sprintf("%s.%s", res.TerraformType(), res.TerraformId())

	return r.matcher.Match([]string{strRes}, false)
}

func (r *DriftIgnore) IsFieldIgnored(res resource.Resource, path []string) bool {
	sprintf := fmt.Sprintf("%s.%s", res.TerraformType(), res.TerraformId())
	p := strings.Join(path, ".")
	full := strings.Join([]string{sprintf, p}, ".")
	return r.matcher.Match([]string{full}, false)
}
