package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ch3mz-za/znox/znox"
)

func main() {
	fmt.Println("Znox V2 - Enjoy!")

	if len(os.Args) == 1 {
		showHelp()
		os.Exit(0)
	}

	enc := flag.NewFlagSet("enc", flag.ExitOnError)
	enci := enc.String("i", "", "Provide an input file to encrypt.")
	enco := enc.String("o", "", "Provide an output filename.")

	dec := flag.NewFlagSet("dec", flag.ExitOnError)
	deci := dec.String("i", "", "Provide an input file to decrypt.")
	deco := dec.String("o", "", "Provide an output filename.")

	pw := flag.NewFlagSet("pw", flag.ExitOnError)
	pwsize := pw.Int("s", 15, "Generate password of given length.")

	switch os.Args[1] {
	case "enc":
		if err := enc.Parse(os.Args[2:]); err != nil {
			log.Println("Error when parsing arguments to enc")
			panic(err)
		}
		if *enci == "" {
			fmt.Println("Provide an input file to encrypt.")
			os.Exit(1)
		}
		if *enco != "" {
			znox.Encryption(*enci, *enco)
		} else {
			znox.Encryption(*enci, *enci+".enc")
		}

	case "dec":
		if err := dec.Parse(os.Args[2:]); err != nil {
			log.Println("Error when parsing arguments to dec")
			panic(err)
		}
		if *deci == "" {
			fmt.Println("Provide an input file to decrypt.")
			os.Exit(1)
		}
		if *deco != "" {
			znox.Decryption(*deci, *deco)
		} else {
			dd := *deci
			o := "decrypted-" + *deci
			if dd[len(dd)-4:] == ".enc" {
				o = "decrypted-" + dd[:len(dd)-4]
			}
			znox.Decryption(*deci, o)
		}

	case "pw":
		if err := pw.Parse(os.Args[2:]); err != nil {
			log.Println("Error when parsing arguments to pw")
			panic(err)
		}
		fmt.Println("Password :", znox.GeneratePassword(*pwsize))

	default:
		showHelp()
	}

	log.Println("Done!")
}

func showHelp() {
	fmt.Println("Example commands:")
	fmt.Println("Encrypt a file : crypto_demo enc -i plaintext.txt -o ciphertext.enc")
	fmt.Println("Decrypt a file : crypto_demo dec -i ciphertext.enc -o decrypted-plaintext.txt")
	fmt.Println("Generate a password : crypto_demo pw -s 15")
}
