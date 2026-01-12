package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"
)

//go:embed data/*.mp3
var soundFS embed.FS

type soundIndex struct {
	names  []string
	byName map[string]string
}

func main() {
	index, err := buildSoundIndex()
	if err != nil {
		exitf("failed to load embedded sounds: %v", err)
	}

	list := flag.Bool("list", false, "list available sound names")
	exportPath := flag.String("export", "", "export selected sound to a file or directory")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-list] [-export <path>] <sound-name>\n", path.Base(os.Args[0]))
		if len(index.names) == 0 {
			fmt.Fprintln(os.Stderr, "No embedded sounds found.")
			return
		}
		fmt.Fprintln(os.Stderr, "Available sounds:")
		for _, name := range index.names {
			fmt.Fprintf(os.Stderr, "  - %s\n", name)
		}
	}
	flag.Parse()

	if *list {
		for _, name := range index.names {
			fmt.Println(name)
		}
		return
	}

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(2)
	}

	name := normalizeSoundName(flag.Arg(0))
	soundPath, ok := index.byName[name]
	if !ok {
		exitf("sound %q not found (use -list to see available names)", name)
	}

	data, err := soundFS.ReadFile(soundPath)
	if err != nil {
		exitf("failed to read %s: %v", soundPath, err)
	}

	if *exportPath != "" {
		outPath, err := resolveExportPath(*exportPath, name)
		if err != nil {
			exitf("failed to resolve export path: %v", err)
		}
		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			exitf("failed to create export directory: %v", err)
		}
		if err := os.WriteFile(outPath, data, 0o644); err != nil {
			exitf("failed to write %s: %v", outPath, err)
		}
		fmt.Printf("exported: %s\n", outPath)
		return
	}

	if err := playMP3(data); err != nil {
		exitf("failed to play %s: %v", name, err)
	}
}

func buildSoundIndex() (soundIndex, error) {
	entries, err := fs.ReadDir(soundFS, "data")
	if err != nil {
		return soundIndex{}, err
	}

	index := soundIndex{
		names:  make([]string, 0, len(entries)),
		byName: make(map[string]string, len(entries)),
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		filename := entry.Name()
		if path.Ext(filename) != ".mp3" {
			continue
		}
		base := strings.TrimSuffix(filename, ".mp3")
		if _, exists := index.byName[base]; exists {
			return soundIndex{}, fmt.Errorf("duplicate sound name: %s", base)
		}
		index.byName[base] = path.Join("data", filename)
		index.names = append(index.names, base)
	}

	sort.Strings(index.names)
	return index, nil
}

func normalizeSoundName(arg string) string {
	if ext := path.Ext(arg); ext != "" {
		return strings.TrimSuffix(arg, ext)
	}
	return arg
}

func playMP3(data []byte) error {
	decoder, err := mp3.NewDecoder(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("decode mp3: %w", err)
	}

	ctx, ready, err := oto.NewContext(decoder.SampleRate(), 2, oto.FormatSignedInt16LE)
	if err != nil {
		return fmt.Errorf("open audio device: %w", err)
	}
	<-ready

	player := ctx.NewPlayer(decoder)
	defer player.Close()
	player.Play()

	for {
		time.Sleep(100 * time.Millisecond)
		if !player.IsPlaying() {
			break
		}
	}
	return nil
}

func resolveExportPath(target string, name string) (string, error) {
	if target == "" {
		return "", fmt.Errorf("export path is empty")
	}
	info, err := os.Stat(target)
	if err == nil && info.IsDir() {
		return filepath.Join(target, name+".mp3"), nil
	}
	if err != nil {
		if os.IsNotExist(err) {
			if strings.HasSuffix(target, string(os.PathSeparator)) || strings.HasSuffix(target, "/") {
				clean := filepath.Clean(target)
				return filepath.Join(clean, name+".mp3"), nil
			}
			return target, nil
		}
		return "", err
	}
	return target, nil
}

func exitf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
