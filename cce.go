package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"strconv"
	"math/rand"
	"time"
	"github.com/ajstarks/svgo"
)

// SVG struct for decoding
type SVG struct {
	XMLName xml.Name `xml:"svg"`
	Circle  []Circle `xml:"circle"`
}

// Circle struct for decoding
type Circle struct {
	Cx string `xml:"cx,attr"`
	Cy string `xml:"cy,attr"`
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

	// Iterate over the coordinates and add them to hexString
	for _, circle := range svg.Circle {
		cx, _ := strconv.Atoi(circle.Cx)
		cy, _ := strconv.Atoi(circle.Cy)
		hexString += fmt.Sprintf("%02x%02x", cx, cy)
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

// randomColor function to generate a random color
func randomColor() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("rgb(%d,%d,%d)", rand.Intn(256), rand.Intn(256), rand.Intn(256))
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

	// Check if the number of coordinate pairs is even
	// if len(hexString)/4%2 != 0 {
	//	fmt.Println("Error: The number of coordinate pairs must be even.")
	//	return
	//}

	// Initialize the SVG canvas
	canvas := svg.New(os.Stdout)
	canvas.Start(255, 255) // Adjust the canvas size to cover the range from 00 to ff

	// Parse the hex string into coordinates and draw the circles
	for i := 0; i+3 < len(hexString); i += 4 {
		cx, _ := strconv.ParseInt(hexString[i:i+2], 16, 64)
		cy, _ := strconv.ParseInt(hexString[i+2:i+4], 16, 64)
		canvas.Circle(int(cx), int(cy), 20, fmt.Sprintf("fill:%s;stroke:black;stroke-width:1", randomColor()))
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

