package utils

import (
    "archive/zip"
    "fmt"
    "io"
    "os"
    "path/filepath"
	"strings"
)

func Unzip(src, dest string) error {
    r, err := zip.OpenReader(src)
    if err != nil {
        return fmt.Errorf("failed to open zip file: %w", err)
    }
    defer r.Close()

    for _, f := range r.File {
        fpath := filepath.Join(dest, f.Name)
        if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
            return fmt.Errorf("illegal file path: %s", fpath)
        }

        if f.FileInfo().IsDir() {
            os.MkdirAll(fpath, os.ModePerm)
            continue
        }

        if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
            return fmt.Errorf("failed to create directory: %w", err)
        }

        outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
        if err != nil {
            return fmt.Errorf("failed to open file: %w", err)
        }

        rc, err := f.Open()
        if err != nil {
            return fmt.Errorf("failed to open file in zip: %w", err)
        }

        _, err = io.Copy(outFile, rc)
        if err != nil {
            return fmt.Errorf("failed to copy file content: %w", err)
        }

        outFile.Close()
        rc.Close()
    }

    return nil
}
