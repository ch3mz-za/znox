package main

import (
	"flag"
	"log"
	"time"

	"github.com/ch3mz-za/znox/znox"
)

func main() {
	encryptPtr := flag.Bool("enc", false, "Encryption")
	decryptPtr := flag.Bool("dec", false, "Decryption")
	compressPtr := flag.Bool("c", false, "Compression")

	flag.Parse()

	if len(flag.Args()) < 2 {
		log.Fatal("Please provide an input and target path")
	}

	sourcePath := flag.Arg(0)
	targetPath := flag.Arg(1)

	zn := znox.NewZnox(sourcePath, targetPath, *compressPtr)

	// Give source path
	// Check if path or file

	var startTime time.Time

	if *encryptPtr {
		log.Println("Starting encryption")
		startTime = time.Now()
		zn.MakeEncryption()
	}

	if *decryptPtr {
		log.Println("Starting decryption")
		startTime = time.Now()
		zn.MakeDecryption()
	}

	log.Println("Done - Time elapsed:", time.Since(startTime).Seconds())
}
