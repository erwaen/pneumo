package minivalkey

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func SendRespCommand(valkeyConn net.Conn, command string, args ...string) (string, error) {
	// for the first line of RESP format (*{len}\r\n)
	elements := append([]string{command}, args...)
	arrayLen := len(elements)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("*%d\r\n", arrayLen))

	// start writing every command with a bulk string
	//(${lenString}\r\n)
	//(${string}\r\n)
	for _, e := range elements {
		sb.WriteString(fmt.Sprintf("$%d\r\n", len(e)))
		sb.WriteString(fmt.Sprintf("%s\r\n", e))
	}

	// send the command
	_, err := valkeyConn.Write([]byte(sb.String()))
	if err != nil {
		return "", fmt.Errorf("erro sending a RESP command: %w", err)
	}

	// read response
	reader := bufio.NewReader(valkeyConn)
	response, err := readRespReply(reader)
	if err != nil {
		return "", err
	}

	return response, nil
}

func readRespReply(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading response from valkey server: %w", err)
	}

	if len(line) < 3 {
		return "", fmt.Errorf("bad format line response, too short: %s", line)
	}

	switch line[0] {
	case '+': // simple string
		return strings.TrimSpace(line[1:]), nil

	case '-': // error
		return "", fmt.Errorf("error from valkey server: %s", strings.TrimSpace(line[1:]))

	case ':': // integer
		return strings.TrimSpace(line[1:]), nil

	case '$': // bulk string
		// get len
		if line[1] == '-' && line[2] == '1' { // $-1 means nil response
			return "", nil // nil response for bulk string
		}

		var bulkLen int
		n, err := fmt.Sscanf(line[1:], "%d\r\n", &bulkLen)
		if n != 1 || err != nil {
			return "", fmt.Errorf("Invalid bulk string length: %s", line)
		}

		// read the data + 2 (\r\n)
		bulkData := make([]byte, bulkLen+2)
		_, err = reader.Read(bulkData)
		if err != nil {
			return "", fmt.Errorf("error reading bulk string data: %w", err)
		}

		// return without the \r\n
		return string(bulkData[:bulkLen]), nil

	// case '*': // array
	// 	var arrayLen int
	// 	n, err:=fmt.Sscanf(line[1:], "%d\r\n", &arrayLen)
	// 	if n != 1 || err != nil {
	// 		return "", fmt.Errorf("Invalid array length: %s", line)
	// 	}
	//
	// 	if arrayLen == 0 {
	// 		return "", nil
	// 	}
	//
	// 	// not handle everythign yet

	default:
		return "", fmt.Errorf("unknown RESP type: %c", line[0])

	}
}
