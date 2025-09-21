package main

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/term"
)

const targetW, targetH = 1600, 900

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "--version" || os.Args[1] == "-v" {
			fmt.Println("tweet-me", version)
			return
		}
	}
	if len(os.Args) > 1 && os.Args[1] == "config" {
		if err := HandleConfig(os.Args[2:]); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		return
	}

	var text, user, handle, retweetsStr, quotesStr, likesStr string

	reader := bufio.NewReader(os.Stdin)

	// Load defaults from config if present
	cfg, _ := loadConfig()

	// For tweet text, prefer an interactive TTY prompt that shows live word count.
	text = readLineInteractive(reader, "Enter tweet text: ")
	if cfg.User != "" {
		user = cfg.User
	} else {
		user = readLine(reader, "Enter display name: ")
	}
	if cfg.Handle != "" {
		handle = cfg.Handle
	} else {
		handle = readLine(reader, "Enter Twitter handle: ")
	}
	retweetsStr = readLine(reader, "Enter retweet count (leave empty for random): ")
	quotesStr = readLine(reader, "Enter quote tweet count (leave empty for random): ")
	likesStr = readLine(reader, "Enter like count (leave empty for random): ")

	// Randomize counts if empty
	retweets := randomOrParse(retweetsStr, 100, 10000)
	quotes := randomOrParse(quotesStr, 10, 1000)
	likes := randomOrParse(likesStr, 200, 20000)

	if text == "" {
		text = "Hello Twitter!"
	}
	if user == "" {
		user = "Hein"
	}
	if handle == "" {
		handle = "@hein_dev"
	}

	// Resolve avatar: prefer ~/.config/tweet-me/profile.png if it exists
	var avatar template.URL
	if dir, err := configDir(); err == nil {
		candidate := filepath.Join(dir, "profile.png")
		if st, err := os.Stat(candidate); err == nil && !st.IsDir() {
			// Use file:// scheme for absolute paths
			abs := candidate
			if !filepath.IsAbs(abs) {
				if a, err := filepath.Abs(abs); err == nil {
					abs = a
				}
			}
			// Use forward slashes for URL
			avatar = template.URL("file://" + strings.ReplaceAll(abs, "\\", "/"))
		}
	}
	if avatar == "" {
		avatar = template.URL("https://randomuser.me/api/portraits/men/81.jpg")
	}

	data := TweetData{
		User:      user,
		Handle:    handle,
		Text:      text,
		Timestamp: time.Now().Format("Jan 2, 2006 at 3:04 PM"),
		Retweets:  socialCounter(retweets),
		Quotes:    socialCounter(quotes),
		Likes:     socialCounter(likes),
		Avatar:    avatar,
		BodyClass: func() string {
			if cfg.Dark {
				return "dark"
			}
			return ""
		}(),
	}

	outputDir := cfg.OutDir
	if outputDir == "" {
		if cwd, err := os.Getwd(); err == nil {
			outputDir = cwd
		} else {
			outputDir = "."
		}
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		panic(err)
	}
	ts := time.Now().Format("20060102-150405")
	htmlFile := "tweet.html"
	output := filepath.Join(outputDir, fmt.Sprintf("tweet-%s.png", ts))

	f, err := os.Create(htmlFile)
	if err != nil {
		panic(err)
	}
	if err := template.Must(template.New("tweet").Parse(htmlTemplate)).Execute(f, data); err != nil {
		panic(err)
	}
	f.Close()

	// Run wkhtmltoimage if available, else just leave the HTML output
	if _, lookErr := exec.LookPath("wkhtmltoimage"); lookErr == nil {
		// Target dimensions must match the CSS canvas in the template
		args := []string{
			"--enable-local-file-access",
			"--width", fmt.Sprint(targetW),
			"--height", fmt.Sprint(targetH),
			"--crop-w", fmt.Sprint(targetW),
			"--crop-h", fmt.Sprint(targetH),
			htmlFile,
			output,
		}
		cmd := exec.Command("wkhtmltoimage", args...)
		if out, err := cmd.CombinedOutput(); err != nil {
			fmt.Printf("⚠️ Failed to run wkhtmltoimage: %v\nOutput: %s\n", err, string(out))
			fmt.Printf("💾 Saved HTML to %s\n", htmlFile)
			return
		}

		fmt.Printf("✅ Saved tweet screenshot to %s\n", output)
	} else {
		fmt.Printf("ℹ️ wkhtmltoimage not found in PATH. Skipping image generation.\n")
		fmt.Printf("💾 Saved HTML to %s\n", htmlFile)
	}
}

func readLine(r *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	line, _ := r.ReadString('\n')
	return strings.TrimSpace(line)
}

func randomOrParse(s string, min, max int) int {
	if s == "" {
		return min + int(time.Now().UnixNano()%int64(max-min+1))
	}
	var v int
	_, err := fmt.Sscanf(s, "%d", &v)
	if err != nil {
		return min + int(time.Now().UnixNano()%int64(max-min+1))
	}
	return v
}

func socialCounter(number int) string {
	switch {
	case number < 1000:
		return fmt.Sprintf("%d", number)
	case number < 1000000:
		return fmt.Sprintf("%.1fK", float64(number)/1000)
	case number < 1000000000:
		return fmt.Sprintf("%.1fM", float64(number)/1000000)
	default:
		return fmt.Sprintf("%.1fB", float64(number)/1000000000)
	}
}

// provides a simple TTY prompt that shows a live word count while typing
func readLineInteractive(r *bufio.Reader, prompt string) string {
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		return readLine(r, prompt)
	}

	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return readLine(r, prompt)
	}
	defer term.Restore(fd, oldState)

	// Terminal width
	cols, _, err := term.GetSize(fd)
	if err != nil || cols <= 0 {
		cols = 80
	}

	var runes []rune

	render := func() (out string, lines int) {
		s := string(runes)
		count := utf8.RuneCountInString(s)
		prefix := fmt.Sprintf("%s(%d) ", prompt, count)
		out = prefix + s
		visible := utf8.RuneCountInString(prefix) + utf8.RuneCountInString(s)
		lines = (visible + cols - 1) / cols
		if lines < 1 {
			lines = 1
		}
		return
	}

	// Print initial state
	out, prevLines := render()
	fmt.Print(out)

	buf := make([]byte, 1)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			break
		}
		b := buf[0]
		if b == '\r' || b == '\n' {
			// Clear previous wrapped lines
			for i := 0; i < prevLines-1; i++ {
				fmt.Print("\x1b[1A")
			}
			for i := 0; i < prevLines; i++ {
				fmt.Print("\x1b[2K\r")
				if i < prevLines-1 {
					fmt.Print("\x1b[1B")
				}
			}
			for i := 0; i < prevLines-1; i++ {
				fmt.Print("\x1b[1A")
			}
			final, _ := render()
			fmt.Print(final)
			fmt.Print("\r\n")
			break
		}
		// backspace (127) or ctrl-h
		if b == 127 || b == 8 {
			if len(runes) > 0 {
				runes = runes[:len(runes)-1]
			}
		} else if b == 3 { // Ctrl-C
			os.Exit(1)
		} else if b >= 32 {
			runes = append(runes, rune(b))
		}

		for i := 0; i < prevLines-1; i++ {
			fmt.Print("\x1b[1A")
		}
		for i := 0; i < prevLines; i++ {
			fmt.Print("\x1b[2K\r")
			if i < prevLines-1 {
				fmt.Print("\x1b[1B")
			}
		}
		for i := 0; i < prevLines-1; i++ {
			fmt.Print("\x1b[1A")
		}
		out, prevLines = render()
		fmt.Print(out)
	}

	return strings.TrimSpace(string(runes))
}
