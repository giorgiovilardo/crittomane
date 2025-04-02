package internal

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func TarPrivate() (*bytes.Buffer, error) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return nil, err
	}

	privateDir := filepath.Join(cwd, "private")

	tarBuffer, err := tarDirectory(privateDir)
	if err != nil {
		fmt.Println("Error creating tar:", err)
		return nil, err
	}

	return tarBuffer, nil
}

func tarDirectory(sourceDir string) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)

	gw := gzip.NewWriter(buf)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	baseDir := filepath.Dir(sourceDir)
	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(baseDir, path)
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		header.Name = relPath

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(tw, file)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if err := tw.Close(); err != nil {
		return nil, err
	}

	if err := gw.Close(); err != nil {
		return nil, err
	}

	return buf, nil
}

func UntarBytes(tarData []byte) error {
	destDir := "."

	buf := bytes.NewReader(tarData)

	gr, err := gzip.NewReader(buf)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		filePath := filepath.Join(destDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(filePath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", filePath, err)
			}

			if err := os.Chmod(filePath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to set directory permissions for %s: %w", filePath, err)
			}

			if err := os.Chtimes(filePath, header.AccessTime, header.ModTime); err != nil {
				return fmt.Errorf("failed to set directory timestamps for %s: %w", filePath, err)
			}

		case tar.TypeReg, tar.TypeRegA:
			if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory for %s: %w", filePath, err)
			}

			outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create output file %s: %w", filePath, err)
			}

			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return fmt.Errorf("failed to write output file %s: %w", filePath, err)
			}
			outFile.Close()

			if err := os.Chmod(filePath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to set file permissions for %s: %w", filePath, err)
			}

			if err := os.Chtimes(filePath, header.AccessTime, header.ModTime); err != nil {
				return fmt.Errorf("failed to set file timestamps for %s: %w", filePath, err)
			}

		case tar.TypeSymlink:
			if err := os.Symlink(header.Linkname, filePath); err != nil {
				return fmt.Errorf("failed to create symlink %s -> %s: %w", filePath, header.Linkname, err)
			}

		default:
			return fmt.Errorf("unsupported file type for %s: %c", header.Name, header.Typeflag)
		}
	}

	return nil
}
