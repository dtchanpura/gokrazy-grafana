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
	lastModifiedTimeStatic := getLastModifiedTime(filepath.Join("tmp", ".timestamp.static"))
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
		lastModifiedTimeStatic = time.Time{}
		err = addTimestampFile(filepath.Join("tmp", ".timestamp.upstream"), newTime)
		if err != nil {
			t.Logf("error adding timestamp file %v", err) // non fatal error
		}
	}
	// extract static binaries
	err = extractFile(
		"tmp",
		fmt.Sprintf(staticBinaryURL, *flagDownloadVersion),
		fmt.Sprintf("grafana-%s/", *flagDownloadVersion),
		lastModifiedTimeStatic,
	)

	if err != nil && err.Error() != "no change" {
		t.Fatalf("error extracting file: %v", err)
	}
	if err == nil {
		err = addTimestampFile(filepath.Join("tmp", ".timestamp.static"), newTime)
		if err != nil && err.Error() != "no change" {
			t.Errorf("error extracting file: %v", err)
		}
	}
}
