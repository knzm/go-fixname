package rename

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

var DiffCmd = "diff"

func Diff(filename string, content []byte) error {
	renamed := fmt.Sprintf("%s.%d.renamed", filename, os.Getpid())
	if err := ioutil.WriteFile(renamed, content, 0644); err != nil {
		return err
	}
	defer os.Remove(renamed)

	diff, err := exec.Command(DiffCmd, "-u", filename, renamed).CombinedOutput()
	if len(diff) > 0 {
		// diff exits with a non-zero status when the files don't match.
		// Ignore that failure as long as we get output.
		os.Stdout.Write(diff)
		return nil
	}
	if err != nil {
		return fmt.Errorf("computing diff: %v", err)
	}
	return nil
}
