package main

import (
	"flag"
	"fmt"

	"github.com/Fabucik/wordlist-distribution/client/connection"
)

func main() {
	var serverIP string
	var localPort string
	var hashMode string
	var hashFile string
	var outputFile string

	flag.StringVar(&serverIP, "server", "", "Specify server's IP")
	flag.StringVar(&localPort, "port", "8888", "Specify listening port")
	flag.StringVar(&hashMode, "hash-type", "0", "Specify hash type for hashcat")
	flag.StringVar(&hashFile, "hash-file", "", "Specify path to file containing hashed password")
	flag.StringVar(&outputFile, "outfile", "", "Specify path to output file (cracked hash will be stored here")

	flag.Parse()

	con := connection.InitializeConnection(serverIP, localPort)

	fmt.Println("Receiving wordlist...")

	wordlist := connection.GetWordlist(con, connection.RecieveWordlistLines(con))

	connection.ConfirmReceivedWordlist(con)
	fmt.Println("Wordlist received!")

	result, success := connection.StartCrack(hashMode, hashFile, outputFile, wordlist, con)

	connection.ReturnResult(con, result, success)
}
