package send

import (
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/Fabucik/wordlist-distribution/checkerr"
	"github.com/Fabucik/wordlist-distribution/server/open"
)

func SendWordlistPart(connection net.Conn, wordlistPart string) {
	wordlistReader := strings.NewReader(wordlistPart)

	lines := open.CountWordlistLines(wordlistPart)

	_, err := io.CopyN(connection, wordlistReader, int64(lines))
	fmt.Println("sent")
	checkerr.CheckError(err)
}
