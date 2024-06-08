package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"uw/ulog"
)

type Entrie struct {
	Version  int           `json:"version"`
	Resource string        `json:"resource"`
	Entries  EntrieObjects `json:"entries"`
}

type EntrieObject struct {
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Source    string `json:"source,omitempty"`
}

type EntrieObjects []EntrieObject

func (es EntrieObjects) Len() int {
	return len(es)
}

func (es EntrieObjects) Less(i, j int) bool {
	return es[i].Timestamp > es[j].Timestamp
}

func (es EntrieObjects) Swap(i, j int) {
	es[i], es[j] = es[j], es[i]
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Welcome to Restory.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "example: %s <option> <history-dir-path> <output-dir-path>\n", os.Args[0])
		os.Exit(1)
	}

	outputHistoryFilesPtr := flag.Bool("h", false, "output history files")
	debugLoggerPtr := flag.Bool("d", false, "debug logger")
	flag.Parse()

	if !*debugLoggerPtr {
		ulog.GlobalFormat().SetLevel(ulog.GlobalFormat().GetLevel() ^ ulog.LevelDebug)
	}

	historyDirPath, outputDirPath := strings.TrimSpace(flag.Arg(0)), strings.TrimSpace(flag.Arg(1))
	if historyDirPath == "" || outputDirPath == "" {
		flag.Usage()
	}

	outputHistoryFiles := *outputHistoryFilesPtr

	ulog.Debug("Restory:\nHistoryDir: %s\nOutputDir: %s\nOutputFiles: %t",
		historyDirPath, outputDirPath, outputHistoryFiles)

	if _, e := os.Stat(outputDirPath); e != nil && !errors.Is(e, os.ErrNotExist) {
		ulog.Fatal("stat %s: %s", outputDirPath, e)
	} else if errors.Is(e, os.ErrNotExist) {
		if e := os.MkdirAll(outputDirPath, os.ModePerm); e != nil {
			ulog.Fatal("create %s: %s", outputDirPath, e)
		}
	}

	historyChunkPaths := make([]string, 0, 512)
	if e := filepath.Walk(historyDirPath, func(path string, info fs.FileInfo, e error) error {
		if e != nil {
			return fmt.Errorf("walk at %s: %w", path, e)
		}

		if !info.IsDir() || path == historyDirPath {
			return nil
		}

		historyChunkPaths = append(historyChunkPaths, path)

		return nil
	}); e != nil {
		ulog.Fatal("walking %s: %s", historyDirPath, e)
	}

	ulog.Debug("chunk: %d", len(historyChunkPaths))
	ps, files := ulog.Progress(10, float64(len(historyChunkPaths)), "chunk"), 0

	for _, path := range historyChunkPaths {
		ulog.Debug("checking: %s", path)

		if e := func(historyChunkPath string) error {
			entriesJsonPath := filepath.Join(historyChunkPath, "entries.json")

			ulog.Debug("opening: %s", entriesJsonPath)
			f, e := os.OpenFile(entriesJsonPath, os.O_RDONLY, 0)
			if e != nil {
				ulog.Warn("open %s: %s, skipping", entriesJsonPath, e)
				return nil
			}

			defer f.Close()

			entrie := &Entrie{}
			if e := json.NewDecoder(f).Decode(&entrie); e != nil {
				return fmt.Errorf("unmarshal %s: %w", entriesJsonPath, e)
			}

			sort.Sort(entrie.Entries)
			if len(entrie.Entries) < 1 {
				return nil
			}
			if !outputHistoryFiles {
				entrie.Entries = entrie.Entries[:1]
			}

			targetOutputPath := filepath.Join(outputDirPath, strings.TrimPrefix(entrie.Resource, "file://"))
			ulog.Debug("restoring: %s", targetOutputPath)

			for i, entry := range entrie.Entries {
				if e := func(i int, entry EntrieObject, targetPath string) error {
					if i > 0 {
						targetPath = fmt.Sprintf("%s.%d", targetPath, i)
					}

					if e := os.MkdirAll(filepath.Dir(targetPath), os.ModePerm); e != nil {
						return fmt.Errorf("create dir %s: %w", filepath.Dir(targetPath), e)
					}

					tf, e := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
					if e != nil {
						return fmt.Errorf("create %s: %w", targetPath, e)
					}

					defer tf.Close()

					historyPath := filepath.Join(historyChunkPath, entry.ID)
					hf, e := os.OpenFile(historyPath, os.O_RDONLY, 0)
					if e != nil {
						ulog.Warn("open %s: %s, skipping", historyPath, e)
						return nil
					}

					defer hf.Close()

					if _, e := io.Copy(tf, hf); e != nil {
						return fmt.Errorf("copy %s to %s: %w", historyPath, targetPath, e)
					}

					ulog.Debug("copied %s to %s", historyPath, targetPath)
					files++
					return nil
				}(i, entry, targetOutputPath); e != nil {
					return fmt.Errorf("copy %s history (%d): %w", entriesJsonPath, i, e)
				}
			}

			return nil
		}(path); e != nil {
			ulog.Fatal("processing %s: %s", path, e)
		}

		ps.Append(1, path)
	}

	ulog.Info("full done, saved: %d chunk (%d files) to %s, bye!",
		len(historyChunkPaths), files, outputDirPath)
}
