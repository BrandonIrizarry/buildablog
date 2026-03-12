package environment

import (
	"fmt"

	"github.com/joho/godotenv"
)

var expectedEnv = []string{
	"BLOGDIR",
	"PORT",
	"SITEURL",
	"TIMEZONE",
	"IS_REPO",
}

func New() (Env, error) {
	// Load environment data into the program, and make sure they
	// all exist and are nonempty.
	var err error

	envMap, err := godotenv.Read()
	if err != nil {
		return Env{}, fmt.Errorf("can't read from environment: %w", err)
	}

	for _, v := range expectedEnv {
		_, ok := envMap[v]
		if !ok {
			return Env{}, fmt.Errorf("Missing environment variable %s", v)
		}
	}

	var isRepo bool

	switch envMap["IS_REPO"] {
	case "true":
		isRepo = true
	case "false":
		isRepo = false
	default:
		return Env{}, fmt.Errorf("Invalid IS_REPO configuration (should be 'true' or 'false')")
	}

	env := Env{
		SiteURL:  envMap["SITEURL"],
		Port:     envMap["PORT"],
		BlogDir:  envMap["BLOGDIR"],
		Timezone: envMap["TIMEZONE"],
		IsRepo:   isRepo,
	}

	return env, nil
}
