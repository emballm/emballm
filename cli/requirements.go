package cli

import (
	"errors"
	"os/exec"
	"runtime"
)

func CheckRequirements() (err error) {
	switch runtime.GOOS {
	case "darwin", "linux", "windows":
	default:
		err = errors.New("unsupported operating system")
		return
	}

	_, err = exec.LookPath("ollama")
	if err != nil {
		err = errors.New("ollama is required")
		return
	}

	return
}
