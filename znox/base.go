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

	baseName := returnBaseName(sourcePath)
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

	Encrypt(sourceToEncrypt, targetEncryption, getAESkey(false))

	removeProcessFiles(zn)
}

func (zn *znox) MakeDecryption() {

	targetDecryption := zn.files.targetTar
	compression := strings.HasSuffix(zn.source, fileTypeCompressedEncryption)
	if compression {
		targetDecryption = zn.files.targetTarGzip
	}

	Decrypt(zn.source, targetDecryption, getAESkey(true))

	if compression {
		UnGzip(targetDecryption, zn.target)
	}

	Untar(zn.files.targetTar, zn.target)

	removeProcessFiles(zn)
}

func getAESkey(justRead bool) []byte {
	aesKey, err := ReadAESkey()
	if err != nil {
		if !justRead {
			err = GenerateAESkey()
			if err != nil {
				log.Fatal("Unable to generate key")
			}
			aesKey, err = ReadAESkey()
		}
		if err != nil {
			log.Fatal("Unable to read key")
		}
	}
	return aesKey
}

func removeProcessFiles(zn *znox) {
	os.Remove(zn.files.targetTar)
	if zn.compression {
		os.Remove(zn.files.targetTarGzip)
	}
}

func returnBaseName(path string) string {
	baseName := filepath.Base(path)
	if strings.Contains(baseName, ".") {
		idx := strings.Index(baseName, ".")
		baseName = baseName[:idx]
	}
	return baseName
}
