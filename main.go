package main

import (
	"flag"
	"log"
	"time"

	"github.com/ch3mz-za/znox/znox"
)

var (
	sourcePath  string
	targetPath  string
	encryptPtr  = flag.Bool("enc", false, "Encryption")
	decryptPtr  = flag.Bool("dec", false, "Decryption")
	compressPtr = flag.Bool("c", false, "compression")
)

func init() {
	flag.Parse()

	if len(flag.Args()) < 2 {
		log.Fatal("Please provide an input and target path")
	}

	sourcePath = flag.Arg(0)
	targetPath = flag.Arg(1)
}

func main() {
	zn := znox.NewZnox(sourcePath, targetPath, *compressPtr)

	var startTime time.Time

	if *encryptPtr {
		startTime = time.Now()
		zn.MakeEncryption()
	}

	if *decryptPtr {
		startTime = time.Now()
		zn.MakeDecryption()
	}

	log.Println("Done - Time elapsed:", time.Since(startTime))
}
