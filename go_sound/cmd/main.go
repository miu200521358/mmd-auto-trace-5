package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path"
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
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-list] <sound-name>\n", path.Base(os.Args[0]))
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

	if err := playMP3(data); err != nil {
		exitf("failed to play %s: %v", name, err)
	}
}

func buildSoundIndex() (soundIndex, error) {
	entries, err := fs.ReadDir(soundFS, "sound")
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
		index.byName[base] = path.Join("sound", filename)
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

func exitf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
