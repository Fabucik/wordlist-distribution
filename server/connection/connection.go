package connection

import (
	"net"
	"strconv"

	"github.com/Fabucik/wordlist-distribution/checkerr"
	"github.com/Fabucik/wordlist-distribution/server/open"
	"github.com/Fabucik/wordlist-distribution/server/send"
)

func InitConection(ipAddress string, port string) net.Conn {
	con, err := net.Dial("tcp", ipAddress+":"+port)
	checkerr.CheckError(err)

	return con
}

func WriteResultToChannel(connection net.Conn, resultChannel *chan string) {
	buffer := make([]byte, 256)

	_, err := connection.Read(buffer)
	checkerr.CheckError(err)

	*resultChannel <- string(buffer)
}

func GetConfirmMessage(connection net.Conn) bool {
	buffer := make([]byte, 20)

	_, err := connection.Read(buffer)
	checkerr.CheckError(err)

	if string(buffer) == "RECEIVED OK" {
		return true
	} else {
		return false
	}
}

func InitSpecificParticipant(connection net.Conn, wordlistPart string) {
	connection.Write([]byte(strconv.Itoa(open.CountWordlistLines(wordlistPart)) + "\n"))

	success := GetConfirmMessage(connection)
	if success {
		send.SendWordlistPart(connection, wordlistPart)
	}
}
