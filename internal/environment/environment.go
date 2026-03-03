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

	env := Env{
		SiteURL:  envMap["SITEURL"],
		Port:     envMap["PORT"],
		BlogDir:  envMap["BLOGDIR"],
		Timezone: envMap["TIMEZONE"],
	}

	return env, nil
}

func (env Env) DraftsDir(genre string) string {
	return fmt.Sprintf("%s/drafts/%s", env.BlogDir, genre)
}

func (env Env) PublishedDir(genre string) string {
	return fmt.Sprintf("%s/published/%s", env.BlogDir, genre)
}
