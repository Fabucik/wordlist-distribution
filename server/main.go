package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/Fabucik/wordlist-distribution/checkerr"
	"github.com/Fabucik/wordlist-distribution/server/connection"
	"github.com/Fabucik/wordlist-distribution/server/open"
)

func main() {
	var howMuchDevices int
	var ipAddresses string
	var ports string
	var wordlistPath string

	flag.IntVar(&howMuchDevices, "devices", 0, "Specify participants in wordlist distribution")
	flag.StringVar(&ipAddresses, "addresses", "", "Specify IP addresses of participants separated by comma without spaces")
	flag.StringVar(&ports, "ports", "8888", "Specify listening ports of participants separated by comma without spaces")
	flag.StringVar(&wordlistPath, "wordlist", "", "Specify wordlist path")

	flag.Parse()

	ipAddressesSplit := strings.Split(ipAddresses, ",")
	portsSplit := strings.Split(ports, ",")

	fullWordlist := open.SplitNTimes(wordlistPath, howMuchDevices)

	fmt.Println(open.CountWordlistLines(fullWordlist[1]))
	fmt.Println(string([]byte(strconv.Itoa(open.CountWordlistLines(fullWordlist[0])))))

	connections := []net.Conn{}

	for i := 0; i <= howMuchDevices-1; i++ {
		con := connection.InitConection(ipAddressesSplit[i], portsSplit[i])
		connections = append(connections, con)
		go connection.InitSpecificParticipant(con, fullWordlist[i])
	}

	fmt.Println("Connected to all participants")

	resultChannel := make(chan string, 256)
	for i := 0; i <= howMuchDevices-1; i++ {
		con := connections[i]

		go connection.WriteResultToChannel(con, &resultChannel)
	}

	var result string

	for {
		resultAndStatus := <-resultChannel

		separator := "$@#%"

		result = strings.Split(resultAndStatus, separator)[0]
		success := strings.Split(resultAndStatus, separator)[1]

		convertedSuccess, err := strconv.ParseBool(success)
		checkerr.CheckError(err)

		if convertedSuccess {
			for i := 0; i <= howMuchDevices-1; i++ {
				r := strings.NewReader("stop")

				io.Copy(connections[i], r)
			}
			break
		}
	}

	fmt.Println("Cracked password")
	fmt.Println("----------------")
	fmt.Println(result)
}
