package connection

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/Fabucik/wordlist-distribution/checkerr"
)

func InitializeConnection(serverIP string, port string) net.Conn {
	server, err := net.Listen("tcp", serverIP+":"+port)
	checkerr.CheckError(err)

	con, err := server.Accept()
	checkerr.CheckError(err)

	return con
}

func RecieveWordlistLines(connection net.Conn) int {
	buffer := make([]byte, 1024)
	_, err := connection.Read(buffer)
	checkerr.CheckError(err)

	strippedBuffer := strings.TrimSuffix(string(buffer), "\n")

	fmt.Print(strippedBuffer)

	numOfLines, err := strconv.Atoi(strippedBuffer)
	checkerr.CheckError(err)

	return numOfLines
}

func GetWordlist(connection net.Conn, wordlistLines int) string {
	buffer := make([]byte, wordlistLines)

	_, err := connection.Read(buffer)
	checkerr.CheckError(err)

	fmt.Println("received")

	return string(buffer)
}

func KillCrackingProcess(cmd exec.Cmd, connection net.Conn, didOthersCrack *chan bool) {
	*didOthersCrack <- false

	buffer := make([]byte, 256)

	_, err := connection.Read(buffer)
	checkerr.CheckError(err)

	err = cmd.Process.Kill()
	checkerr.CheckError(err)

	*didOthersCrack <- true
}

func StartCrack(hashMode string, hashFile string, outputFile string, wordlist string, connection net.Conn) (result string, success bool) {
	err := os.WriteFile(outputFile, []byte(""), 0700)
	checkerr.CheckError(err)

	args := []string{"--quiet", "--outfile", outputFile, "-m", hashMode, "-a", "0", hashFile, wordlist}

	fmt.Println("Starting hashcat...")

	crackCheckChannel := make(chan bool, 1)

	command := exec.Command("hashcat", args...)
	go KillCrackingProcess(*command, connection, &crackCheckChannel)

	err = command.Run()
	checkerr.CheckError(err)

	fmt.Println("Finished cracking")

	didOthersCrack := <-crackCheckChannel
	if didOthersCrack {
		fmt.Println("Someone else cracked")
		return "", false
	}

	data, err := os.ReadFile(outputFile)
	checkerr.CheckError(err)

	contents := string(data)
	if contents == "" {
		fmt.Println("Cracking unsuccessful")
		return "", false
	} else {
		contents = strings.SplitN(contents, ",", 2)[1]
		fmt.Println("Cracking successful - " + contents)
		return contents, true
	}
}

func ReturnResult(connection net.Conn, result string, success bool) {
	defer connection.Close()

	formattedSuccess := strconv.FormatBool(success)

	separator := "$@#%"

	stringToSend := result + separator + formattedSuccess

	reader := strings.NewReader(stringToSend)

	io.Copy(connection, reader)
}

func ConfirmReceivedWordlist(connection net.Conn) {
	_, err := connection.Write([]byte("RECEIVED OK"))
	checkerr.CheckError(err)
}
