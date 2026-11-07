package version

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

const (
	AppName = "personal-finance-tracker"
)

type Info struct {
	App      string `json:"app"`
	Version  string `json:"version"`
	Commit   string `json:"commit"`
	BuiltAt  string `json:"built_at"`
}

func Current() Info {
	commit := os.Getenv("COMMIT")
	if commit == "" {
		commit = "dev"
	}
	builtAt := os.Getenv("BUILT_AT")
	if builtAt == "" {
		builtAt = time.Now().UTC().Format(time.RFC3339)
	}
	return Info{
		App:     AppName,
		Version: "1.0.0",
		Commit:  commit,
		BuiltAt: builtAt,
	}
}

func WriteJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(Current())
}
