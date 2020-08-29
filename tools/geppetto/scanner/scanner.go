package scanner

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/geppetto/types"
)

type Scanner struct{}

func New() *Scanner {
	return &Scanner{}
}

func (s *Scanner) ScanDir(targetDir string) ([]types.Repo, error) {
	filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		ignore := []string{"node_modules", ".git"}
		isIgnore := func(name string) bool {
			for _, i := range ignore {
				if name == i {
					return true
				}
			}
			return false
		}

		name := info.Name()

		if !info.IsDir() {
			return nil
		}

		if isIgnore(name) {
			return filepath.SkipDir
		}

		if hasFile(path, "chart.yaml") {
			logrus.Infof("HELM: %s", path)
			return filepath.SkipDir
		}

		if hasFile(path, "package.json") {
			logrus.Infof("NPM: %s", path)
			return filepath.SkipDir
		}

		if hasFile(path, "go.mod") {
			logrus.Infof("GOLANG: %s", path)
			return filepath.SkipDir
		}

		logrus.Debugf("- Unrecognised: %s", path)
		return nil
	})

	return nil, nil
}

func hasFile(path string, filename string) bool {
	file := path + "/" + filename

	info, err := os.Stat(file)
	if err != nil {
		return false
	}

	if info.IsDir() {
		return false
	}

	return true
}
