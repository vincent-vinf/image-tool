package zip

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/klauspost/pgzip"
)

type Compression int

const (
	Uncompressed Compression = 0 // Uncompressed represents the uncompressed.
	Bzip2        Compression = 1 // Bzip2 is bzip2 compression algorithm.
	Gzip         Compression = 2 // Gzip is gzip compression algorithm.
	Xz           Compression = 3 // Xz is xz compression algorithm.
	Zstd         Compression = 4 // Zstd is zstd compression algorithm.
	Zip          Compression = 5
)

var (
	bzip2Magic = []byte{0x42, 0x5A, 0x68}
	gzipMagic  = []byte{0x1F, 0x8B, 0x08}
	xzMagic    = []byte{0xFD, 0x37, 0x7A, 0x58, 0x5A, 0x00}
	zstdMagic  = []byte{0x28, 0xb5, 0x2f, 0xfd}
	zipMagic   = []byte{0x50, 0x4B, 0x03, 0x04}
)

type matcher = func([]byte) bool

func DetectFileCompression(filePath string) (Compression, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	buffer := make([]byte, 10)
	_, err = io.ReadFull(file, buffer)
	if err != nil {
		return 0, err
	}

	return DetectCompression(buffer), nil
}

func DetectCompression(source []byte) Compression {
	compressionMap := map[Compression]matcher{
		Bzip2: magicNumberMatcher(bzip2Magic),
		Gzip:  magicNumberMatcher(gzipMagic),
		Xz:    magicNumberMatcher(xzMagic),
		Zip:   magicNumberMatcher(zipMagic),
		Zstd:  zstdMatcher(),
	}
	for _, compression := range []Compression{Gzip, Zip} {
		fn := compressionMap[compression]
		if fn(source) {
			return compression
		}
	}

	return Uncompressed
}

func magicNumberMatcher(m []byte) matcher {
	return func(source []byte) bool {
		return bytes.HasPrefix(source, m)
	}
}

const (
	zstdMagicSkippableStart = 0x184D2A50
	zstdMagicSkippableMask  = 0xFFFFFFF0
)

func zstdMatcher() matcher {
	return func(source []byte) bool {
		if bytes.HasPrefix(source, zstdMagic) {
			// Zstandard frame
			return true
		}
		// skippable frame
		if len(source) < 8 {
			return false
		}
		// magic number from 0x184D2A50 to 0x184D2A5F.
		if binary.LittleEndian.Uint32(source[:4])&zstdMagicSkippableMask == zstdMagicSkippableStart {
			return true
		}
		return false
	}
}

func Compress(src, dst string, compression Compression) error {
	file, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer file.Close()

	// 创建目标文件
	outFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer outFile.Close()

	var writer io.WriteCloser
	switch compression {
	case Gzip:
		writer = pgzip.NewWriter(outFile)
	default:
		return fmt.Errorf("unsupported compression: %v", compression)
	}
	defer writer.Close()

	// 将解压缩的数据写入到目标文件
	_, err = io.Copy(writer, file)
	if err != nil {
		return fmt.Errorf("failed to write compressed data to file: %v", err)
	}

	return nil
}

func Decompress(src, dst string, compression Compression) error {
	file, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer file.Close()

	var reader io.ReadCloser
	switch compression {
	case Gzip:
		reader, err = pgzip.NewReader(file)
		if err != nil {
			return fmt.Errorf("failed to create gzip reader: %v", err)
		}
	default:
		return fmt.Errorf("unsupported compression: %v", compression)
	}
	defer reader.Close()

	// 创建目标文件
	outFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer outFile.Close()

	// 将解压缩的数据写入到目标文件
	_, err = io.Copy(outFile, reader)
	if err != nil {
		return fmt.Errorf("failed to write uncompressed data to file: %v", err)
	}

	return nil
}
