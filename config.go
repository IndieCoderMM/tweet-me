package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	User   string `json:"user"`
	Handle string `json:"handle"`
	OutDir string `json:"output_directory"`
	Dark   bool   `json:"dark"`
}

func configDir() (string, error) {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "tweet-me"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "tweet-me"), nil
}

func configPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

func loadConfig() (Config, error) {
	var c Config
	path, err := configPath()
	if err != nil {
		return c, err
	}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return c, nil
		}
		return c, err
	}
	defer f.Close()
	dec := json.NewDecoder(f)
	if err := dec.Decode(&c); err != nil {
		return Config{}, err
	}
	return c, nil
}

func saveConfig(c Config) error {
	dir, err := configDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	path := filepath.Join(dir, "config.json")
	tmp := path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(c); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// HandleConfig manages the `config` subcommand
func HandleConfig(args []string) error {
	fs := flag.NewFlagSet("config", flag.ExitOnError)
	userFlag := fs.String("user", "", "Display name (e.g. 'Hein')")
	handleFlag := fs.String("handle", "", "Handle (e.g. '@hein_dev')")
	outDirFlag := fs.String("outdir", "", "Output directory for generated files")
	darkFlag := fs.String("dark", "", "Dark mode: 'true'/'false' or 'toggle' to flip")
	_ = fs.Parse(args)

	existing, _ := loadConfig()
	c := existing

	// If no flags provided, display current config and exit
	if *userFlag == "" && *handleFlag == "" && *outDirFlag == "" && *darkFlag == "" {
		cfgDir, _ := configDir()
		avatar := filepath.Join(cfgDir, "profile.png")
		avatarStatus := "not found"
		if st, err := os.Stat(avatar); err == nil && !st.IsDir() {
			avatarStatus = "ok"
		}
		p, _ := configPath()
		fmt.Printf("Config file: %s\n", p)
		fmt.Printf("User: %s\n", existing.User)
		fmt.Printf("Handle: %s\n", existing.Handle)
		fmt.Printf("Avatar: %s (%s)\n", avatar, avatarStatus)
		if existing.OutDir != "" {
			if st, err := os.Stat(existing.OutDir); err == nil && st.IsDir() {
				fmt.Printf("OutputDir: %s (ok)\n", existing.OutDir)
			} else {
				fmt.Printf("OutputDir: %s (missing)\n", existing.OutDir)
			}
		} else {
			fmt.Printf("OutputDir: (not set)\n")
		}
		fmt.Printf("Dark: %v\n", existing.Dark)
		return nil
	}

	if *userFlag != "" {
		c.User = *userFlag
	}
	if *handleFlag != "" {
		c.Handle = *handleFlag
	}
	if *outDirFlag != "" {
		dir := *outDirFlag
		if strings.HasPrefix(dir, "~") {
			if home, err := os.UserHomeDir(); err == nil {
				dir = filepath.Join(home, strings.TrimPrefix(dir, "~"))
			}
		}
		if !filepath.IsAbs(dir) {
			if abs, err := filepath.Abs(dir); err == nil {
				dir = abs
			}
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create output dir: %w", err)
		}
		c.OutDir = dir
	}

	if *darkFlag != "" {
		v := strings.ToLower(strings.TrimSpace(*darkFlag))
		switch v {
		case "true", "1", "yes", "on":
			c.Dark = true
		case "false", "0", "no", "off":
			c.Dark = false
		case "toggle":
			c.Dark = !c.Dark
		default:
			return fmt.Errorf("invalid -dark value: %q (use true/false/toggle)", *darkFlag)
		}
	}

	if err := saveConfig(c); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	p, _ := configPath()
	fmt.Printf("✅ Saved config to %s\n", p)
	return nil
}
