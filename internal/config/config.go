package config

import (
	"bufio"
	"os"
	"strings"
)

type Config struct {
	Port      string
	Domain    string
	LocalPath string

	S3Bucket string
	S3Region string
	S3URL    string
}

func Load() *Config {
	loadDotEnv(".env")

	return &Config{
		Port:      get("PORT", "8080"),
		Domain:    strings.TrimRight(get("DOMAIN", "http://localhost:8080"), "/"),
		LocalPath: get("LOCAL_PATH", "./img"),

		S3Bucket: get("S3_BUCKET", ""),
		S3Region: get("S3_REGION", ""),
		S3URL:    get("S3_URL", ""),
	}
}

func loadDotEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, "\"'")

		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
}

func get(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
