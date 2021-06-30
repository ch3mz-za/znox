package znox

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type znox struct {
	source       string
	target       string
	compression  bool
	baseFilename string
	files        processFiles
}

type processFiles struct {
	targetTar           string
	targetTarGzip       string
	targetEnc           string
	targetCompressedEnc string
}

const (
	fileTypeCompressedEncryption   = ".encc"
	fileTypeUncompressedEncryption = ".enc"
	fileTypeTarGzip                = ".tar.gz"
	fileTypeTar                    = ".tar"
	fileTypeGzip                   = ".gz"
)

func NewZnox(sourcePath, targetPath string, compression bool) *znox {

	baseName := filepath.Base(sourcePath)
	if strings.Contains(baseName, ".") {
		idx := strings.Index(baseName, ".")
		baseName = baseName[:idx]
	}

	baseTargetPath := filepath.Join(targetPath, baseName)
	files := processFiles{
		targetTar:           baseTargetPath + fileTypeTar,
		targetTarGzip:       baseTargetPath + fileTypeTarGzip,
		targetEnc:           baseTargetPath + fileTypeUncompressedEncryption,
		targetCompressedEnc: baseTargetPath + fileTypeCompressedEncryption,
	}

	return &znox{
		source:       sourcePath,
		target:       targetPath,
		compression:  compression,
		baseFilename: baseName,
		files:        files,
	}
}

func (zn *znox) MakeEncryption() {

	Tar(zn.source, zn.target)
	sourceToEncrypt := zn.files.targetTar
	targetEncryption := zn.files.targetEnc
	if zn.compression {
		if err := Gzip(zn.files.targetTar, zn.target); err != nil {
			log.Fatal("Compression failed:", err)
		}
		sourceToEncrypt = zn.files.targetTarGzip
		targetEncryption = zn.files.targetCompressedEnc
	}

	Encrypt(sourceToEncrypt, targetEncryption, ReadAESkey())

	removeProcessFiles(zn)
}

func (zn *znox) MakeDecryption() {

	targetDecryption := zn.files.targetTar
	compression := strings.HasSuffix(zn.source, fileTypeCompressedEncryption)
	if compression {
		targetDecryption = zn.files.targetTarGzip
	}

	Decrypt(zn.source, targetDecryption, ReadAESkey())

	if compression {
		UnGzip(targetDecryption, zn.target)
	}

	Untar(zn.files.targetTar, zn.target)

	removeProcessFiles(zn)
}

func removeProcessFiles(zn *znox) {
	os.Remove(zn.files.targetTar)
	if zn.compression {
		os.Remove(zn.files.targetTarGzip)
	}
}
