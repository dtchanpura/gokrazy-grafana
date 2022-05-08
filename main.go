package main

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	downloadURL = "https://dl.grafana.com/oss/release/grafana-%s.linux-arm64.tar.gz"
)

var (
	flagDownloadFiles   = flag.Bool("download-files", false, "Download files before starting")
	flagDownloadVersion = flag.String("download-version", "8.5.2", "Version to download")
	flagBaseDir         = flag.String("base-dir", os.Getenv("HOME"), "Base Directory to use")
	flagTimestampFile   = flag.String("timestamp-file", filepath.Join(os.Getenv("HOME"), ".timestamp"), "Last updated timestamp file")
)

func main() {
	flag.Parse()

	if !*flagDownloadFiles {
		if err := start(); err != nil {
			log.Fatal(err)
		}
		return
	}
	err := os.MkdirAll(*flagBaseDir, 0755)
	if err != nil {
		log.Fatal(err)
	}
	lastModifiedTime := getLastModifiedTime(*flagTimestampFile + ".upstream")
	// lastModifiedTimeStatic := getLastModifiedTime(*flagTimestampFile + ".static")
	newTime := time.Now()
	err = extractFile(
		*flagBaseDir,
		fmt.Sprintf(downloadURL, *flagDownloadVersion),
		fmt.Sprintf("grafana-%s/", *flagDownloadVersion),
		lastModifiedTime,
	)

	if err != nil && err.Error() != "no change" {
		log.Fatal(err)
	}
	if err == nil {
		// forcefully update static file as there is a change in upstream
		// lastModifiedTimeStatic = time.Time{}
		err = addTimestampFile(*flagTimestampFile+".upstream", newTime)
		if err != nil {
			log.Println(err) // non fatal error
		}
	}

	if err := start(); err != nil {
		log.Println(err)
		os.Exit(125)
	}
}

func start() error {
	var bin = filepath.Join("/usr", "local", "bin", "grafana-server")
	if err := syscall.Exec(bin, []string{bin, "-homepath=/perm/grafana"}, nil); err != nil {
		return err
	}
	return nil
}

// reference https://stackoverflow.com/a/31491033/4414721
func extractFile(baseDir string, url string, stripPrefix string, modifiedSince time.Time) error {
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	r.Header.Set("If-Modified-Since", modifiedSince.UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT"))
	d, err := http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	defer d.Body.Close()

	if d.StatusCode == http.StatusNotModified {
		log.Println("contents not modified, skipping download")
		return errors.New("no change")
	}
	// uncompress
	gr, err := gzip.NewReader(d.Body)
	if err != nil {
		return err
	}
	// unarchive
	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			return err
		}
		filename := strings.TrimPrefix(hdr.Name, stripPrefix)
		switch hdr.Typeflag {
		case tar.TypeDir:
			// create a directory
			fmt.Println("creating:   " + filename)
			err = os.MkdirAll(filepath.Join(baseDir, filename), hdr.FileInfo().Mode())
			if err != nil {
				return err
			}
		case tar.TypeReg, tar.TypeRegA:
			// write a file
			fmt.Println("extracting: " + filename)
			w, err := os.OpenFile(filepath.Join(baseDir, filename), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, hdr.FileInfo().Mode())
			if err != nil {
				return err
			}
			_, err = io.Copy(w, tr)
			if err != nil {
				return err
			}
			w.Close()
		}
	}
	return nil
}

func getLastModifiedTime(filename string) time.Time {
	i, err := os.Stat(filename)
	if err != nil {
		return time.Time{}
	}
	return i.ModTime()
}

func addTimestampFile(filename string, t time.Time) error {
	if err := os.WriteFile(filename, []byte(strconv.FormatInt(t.UnixNano(), 10)), 0644); err != nil {
		return err
	}
	if err := os.Chtimes(filename, t, t); err != nil {
		return err
	}
	return nil
}
