package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"strconv"
	"github.com/ajstarks/svgo"
)

// SVG struct for decoding
type SVG struct {
	XMLName xml.Name `xml:"svg"`
	Line    []Line   `xml:"line"`
}

// Line struct for decoding
type Line struct {
	X1 string `xml:"x1,attr"`
	Y1 string `xml:"y1,attr"`
	X2 string `xml:"x2,attr"`
	Y2 string `xml:"y2,attr"`
}

// decode function to convert SVG to hex string
func decode() {
	// Read the SVG string from stdin
	bytes, _ := ioutil.ReadAll(os.Stdin)

	// Decode the SVG string into a struct
	var svg SVG
	xml.Unmarshal(bytes, &svg)

	// Convert the coordinates back into a hex string
	var hexString string

	// Add the first coordinates to hexString
	x1, _ := strconv.Atoi(svg.Line[0].X1)
	y1, _ := strconv.Atoi(svg.Line[0].Y1)
	hexString += fmt.Sprintf("%02x%02x", x1, y1)

	// Iterate over the remaining coordinates and add them to hexString
	for _, line := range svg.Line {
		x2, _ := strconv.Atoi(line.X2)
		y2, _ := strconv.Atoi(line.Y2)
		hexString += fmt.Sprintf("%02x%02x", x2, y2)
	}

	// Replace all whitespaces, newlines, and carriage returns in hexString
	hexString = strings.ReplaceAll(hexString, " ", "")
	hexString = strings.ReplaceAll(hexString, "\n", "")
	hexString = strings.ReplaceAll(hexString, "\r", "")

	// Output the hex string only if there are coordinates
	if hexString != "" {
		// Print the hex string without whitespaces
		fmt.Printf("%s\n", hexString)
	} else {
		fmt.Println("Error: No coordinates found.")
	}
}

// drawSVG function to draw SVG from hex string
func drawSVG(hexString string) {
	// Check if the hex string is even
	if len(hexString)%2 != 0 {
		fmt.Println("Error: Hex string must be even.")
		return
	}

	// Check if the number of hex bytes is even
	if len(hexString)/2%2 != 0 {
		fmt.Println("Error: The number of hex bytes must be even.")
		return
	}

	// Initialize the SVG canvas
	canvas := svg.New(os.Stdout)
	canvas.Start(255, 255) // Adjust the canvas size to cover the range from 00 to ff

	// Parse the hex string into coordinates and draw the lines
	var x1, y1 int64
	for i := 0; i+3 < len(hexString); i += 4 {
		x2, _ := strconv.ParseInt(hexString[i:i+2], 16, 64)
		y2, _ := strconv.ParseInt(hexString[i+2:i+4], 16, 64)
		if i != 0 {
			canvas.Line(int(x1), int(y1), int(x2), int(y2), "stroke:blue;stroke-width:2")
		}
		x1, y1 = x2, y2
	}

	canvas.End()
}

func main() {
	// Check if the correct number of arguments are provided
	if len(os.Args) > 1 {
		if os.Args[1] == "-d" {
			decode()
		} else {
			fmt.Println("Usage: $ echo <hex string (min. 2 hex bytes)> | polyenc [-d]")
			return
		}
	} else {
		// If encoding, read hex string from stdin and draw the SVG
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		hexString := scanner.Text()
		drawSVG(hexString)
	}
}
