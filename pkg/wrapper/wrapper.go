package wrapper

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

// OctantRunner wraps logic for figuring out where octant is.
type OctantRunner struct {
}

// NewOctantRunner creates an instance of OctantRunner.
func NewOctantRunner() *OctantRunner {
	return &OctantRunner{}
}

func (or *OctantRunner) BinPath() (string, error) {
	ch, err := configHome()
	if err != nil {
		return "", fmt.Errorf("find config home: %w", err)
	}

	configRoot := filepath.Join(ch, "k8s-lab")
	binDir := filepath.Join(configRoot, "bin")

	fi, err := os.Stat(binDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}

		if err := os.MkdirAll(binDir, 0755); err != nil {
			return "", err
		}
	} else if !fi.IsDir() {
		return "", fmt.Errorf("%s should be a directory", fi.Name())
	}

	bin, err := binName()
	if err != nil {
		return "", fmt.Errorf("can't find bin name: %w", err)
	}

	binPath := filepath.Join(binDir, bin)
	fi, err = os.Stat(binPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}

		if err := fetchOctant(binPath); err != nil {
			return "", fmt.Errorf("fetch octant: %w", err)
		}
	} else {
		if fi.IsDir() {
			return "", fmt.Errorf("octant bin is a directory")
		}
	}

	return binPath, nil
}

func fetchOctant(binPath string) error {
	bu, err := binURL()
	if err != nil {
		return fmt.Errorf("build bin url: %w", err)
	}

	logrus.WithField("url", bu).Info("downloading octant")
	resp, err := http.Get(bu)
	if err != nil {
		return fmt.Errorf("fetch %s: %q", bu, err)
	}
	defer func() {
		if cErr := resp.Body.Close(); cErr != nil {
			logrus.WithError(cErr).Error("close response body")
		}
	}()

	f, err := os.Create(binPath)
	if err != nil {
		return err
	}
	defer func() {
		if cErr := f.Close(); cErr != nil {
			logrus.WithError(cErr).Error("close octant bin file")
		}
	}()
	if _, err := io.Copy(f, resp.Body); err != nil {
		return fmt.Errorf("write octant bin")
	}

	if err := os.Chmod(binPath, 0755); err != nil {
		return err
	}

	return nil
}

func configHome() (string, error) {
	switch runtime.GOOS {
	case "windows":
		return os.Getenv("APPDATA"), nil
	case "linux", "darwin":
		return filepath.Join(os.Getenv("HOME"), ".config"), nil
	}

	return "", fmt.Errorf("unsupported os %q", runtime.GOOS)
}

func binName() (string, error) {
	switch runtime.GOOS {
	case "windows":
		return "octant.exe", nil
	case "linux":
		return "octant-linux-amd64", nil
	case "darwin":
		return "octant-darwin-amd64", nil
	}

	return "", fmt.Errorf("unsupported os %q", runtime.GOOS)
}

func binURL() (string, error) {
	u, err := url.Parse("https://storage.googleapis.com/bryanl-k8s-lab/v0.1.0")
	if err != nil {
		return "", err
	}

	filename, err := binName()
	if err != nil {
		return "", err
	}

	u.Path = path.Join(u.Path, filename)

	return u.String(), nil
}
