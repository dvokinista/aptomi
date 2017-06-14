package slinga

import (
	"io"
	"os"
	"io/ioutil"
	log "github.com/Sirupsen/logrus"
)

// copyFileContents copies the contents of the file named src to the file named
// by dst. Dst will be overwritten if already exists. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFile(src, dst string) (err error) {
	/*
	if stat, err := os.Stat(dst); err == nil && !stat.IsDir() {
		return fmt.Errorf("File %s already exists", dst)
	}
	*/

	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

// deleteFile deletes a file
func deleteFile(src string) (err error) {
	return os.Remove(src)
}

func writeTempFile(prefix string, content string) *os.File {
	tmpFile, err := ioutil.TempFile("", "aptomi-" + prefix)
	if err != nil {
		debug.WithFields(log.Fields{
			"prefix":  prefix,
			"error": err,
		}).Fatal("Failed to create temp file")
	}

	_, err = tmpFile.Write([]byte(content))
	if err != nil {
		debug.WithFields(log.Fields{
			"file":  tmpFile.Name(),
			"content": content,
			"error": err,
		}).Fatal("Failed to write to temp file")
	}

	return tmpFile
}
