package goprettypackets

import "fmt"

const (
	TableWidth = 10
)

func FormatRawPacket(array []byte) string {
	var output = ""

	if len(array) == 0 {
		return output
	}

	chunkedOut := DivideIntoChunks(array)

	for _, chunk := range chunkedOut {

		//display bytes
		output += DisplayBytes(chunk)

		//add separator
		output += "| "
		output += ChunkToString(chunk)

		//end of line
		output += "\n\t\t"
	}

	return output
}

func DisplayBytes(chunk []byte) string {
	var buffer string
	charactersPerRow := len(chunk)

	for _, _byte := range chunk {
		buffer += fmt.Sprintf("%03d ", _byte)
	}

	for ; charactersPerRow < TableWidth; charactersPerRow++ {
		buffer += "... "
	}
	return buffer
}

func ChunkToString(chunk []byte) string {
	var buffer string
	for _, char := range chunk {
		//display text
		if (char >= 0x21 && char <= 0x7e) || char >= 0xA1 {
			buffer += fmt.Sprintf("%c ", char)
		} else {
			buffer += ". "
		}

	}
	return buffer
}

func DivideIntoChunks(array []byte) [][]byte {
	var divided [][]byte

	for i := 0; i < len(array); i += TableWidth {
		end := i + TableWidth

		if end > len(array) {
			end = len(array)
		}

		divided = append(divided, array[i:end])
	}

	return divided
}
