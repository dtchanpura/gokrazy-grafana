package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestExtract(t *testing.T) {
	os.Mkdir("tmp", 0755)
	lastModifiedTime := getLastModifiedTime(filepath.Join("tmp", ".timestamp.upstream"))
	newTime := time.Now()
	err := extractFile(
		"tmp",
		fmt.Sprintf(downloadURL, *flagDownloadVersion),
		fmt.Sprintf("grafana-%s/", *flagDownloadVersion),
		lastModifiedTime,
	)
	if err != nil && err.Error() != "no change" {
		t.Fatalf("error extracting file: %v", err)
	}
	if err == nil {
		err = addTimestampFile(filepath.Join("tmp", ".timestamp.upstream"), newTime)
		if err != nil {
			t.Logf("error adding timestamp file %v", err) // non fatal error
		}
	}
}
