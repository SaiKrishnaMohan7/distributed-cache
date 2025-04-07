package config

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// LoadDotEnv loads environment variables from a .env file at the specified path.
func LoadDotEnv(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	log.Printf("loading .env at path: %s", path)

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		// skip malformed lines
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), `"`) // remove quotes and whitespace

		_ = os.Setenv(key, value)
	}

	log.Print("Done setting env vars")

	return scanner.Err()
}

// InitLogger initializes the logger with specific flags for detailed logging.
func InitLogger() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Logger initialized!")
}
