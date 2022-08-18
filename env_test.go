package env_test

import (
	"os"
	"testing"
	"time"

	"github.com/hugjobk/go-env"
)

func TestParseEnv(t *testing.T) {
	var Env = struct {
		Env1  int    `env:"ENV1"`
		Env2  string `env:"ENV2"`
		Env34 struct {
			Env3 bool          `env:"ENV3,true"`
			Env4 time.Duration `env:"ENV4,1h30m"`
		}
	}{}
	os.Setenv("ENV1", "123")
	os.Setenv("ENV2", "abc")
	os.Setenv("ENV4", "2h22m22s")
	if err := env.ParseEnv(&Env); err != nil {
		t.Error(err)
	} else {
		t.Logf("%+v", Env)
	}
}
