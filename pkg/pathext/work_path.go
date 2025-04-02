package pathext

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fekoneko/piximan/pkg/work"
)

func FormatWorkPath(path string, work *work.Work) (string, error) {
	replacer := strings.NewReplacer(
		"{title}", work.Title,
		"{id}", strconv.FormatUint(work.Id, 10),
		"{user}", work.UserName,
		"{userid}", strconv.FormatUint(work.UserId, 10),
	)

	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	sections := strings.Split(path, string(filepath.Separator))
	if len(sections) != 0 && len(sections[0]) == 0 {
		sections[0] = string(filepath.Separator)
	}

	for i, section := range sections {
		if i == 0 || section == "." || section == ".." {
			continue
		}
		filename := replacer.Replace(section)
		sections[i] = ToValidFilename(filename)
	}

	return filepath.Join(sections...), nil
}
