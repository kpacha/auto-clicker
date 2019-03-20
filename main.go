package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/go-vgo/robotgo"
)

func main() {
	framesSrc := flag.String("f", "frames.json", "path of the file containning the animation")
	sleep := flag.Duration("s", 45*time.Second, "time to sleep between clicks")
	flag.Parse()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	log.Println("select the app to click on")
	time.Sleep(10 * time.Second)

	pointsToClick := [][]int{}

	for {
		log.Println("place the mouse over the pos to click and press k")
		k := robotgo.AddEvent("k")
		if k {
			x, y := robotgo.GetMousePos()
			log.Println("new point:", x, y)
			pointsToClick = append(pointsToClick, []int{x, y})
		}
		last := len(pointsToClick) - 1
		if last > 1 && distance(pointsToClick[last], pointsToClick[last-1]) < 25 {
			pointsToClick = pointsToClick[:last-1]
			break
		}
	}

	log.Printf("starting with %d targets", len(pointsToClick))

	wg := new(sync.WaitGroup)
	wg.Add(2)

	logCh := make(chan string)
	go func() {
		defer wg.Done()

		render(ctx, *framesSrc, logCh)
	}()

	go func() {
		defer wg.Done()

		totalClicks := 0
		for {
			for i, p := range pointsToClick {
				logCh <- fmt.Sprintf("%s moving mouse to point %d [%d,%d]", time.Now().Format("15:04:05"), i+1, p[0], p[1])
				robotgo.MoveMouseSmooth(p[0], p[1], 1.0, 100.0)
				logCh <- fmt.Sprintf("%s clicking", time.Now().Format("15:04:05"))
				robotgo.MouseClick("left", false)
				logCh <- fmt.Sprintf("%s waiting after clicking #%d", time.Now().Format("15:04:05"), i+1)
				totalClicks++
				if totalClicks%10 == 0 {
					logCh <- fmt.Sprintf("%s clicks so far: %d", time.Now().Format("15:04:05"), totalClicks)
				}

				select {
				case <-ctx.Done():
					return
				case <-time.After(*sleep):
				}
			}
		}
	}()

	<-sigs
	cancel()
}

func distance(p0, p1 []int) int {
	dx := p0[0] - p1[0]
	dy := p0[1] - p1[1]
	return dx*dx + dy*dy
}

func render(ctx context.Context, framesSrc string, in chan string) {
	// Set output character
	const outputChar = "  "

	// Set colors
	colors := map[string]string{
		"+": "226",
		"@": "223",
		",": "17",
		"-": "205",
		"#": "82",
		".": "15",
		"$": "219",
		"%": "217",
		";": "99",
		"&": "214",
		"=": "39",
		"'": "0",
		">": "196",
		"*": "245",
	}

	// Import frames from data file
	framesFile, _ := filepath.Abs(framesSrc)
	data, _ := ioutil.ReadFile(framesFile)

	var frames [][]string
	json.Unmarshal(data, &frames)

	// Get TTY size
	termWidth, termHeight := termSize()

	// Calculate the width in terms of the output char
	termWidth = termWidth / len(outputChar)

	minRow := 0
	maxRow := len(frames[0])

	minCol := 0
	maxCol := len(frames[0][0])

	if maxRow > termHeight {
		minRow = (maxRow - termHeight) / 2
		maxRow = minRow + termHeight
	}

	if maxCol > termWidth {
		minCol = (maxCol - termWidth) / 2
		maxCol = minCol + termWidth
	}

	// Initialize term
	fmt.Print("\033[H\033[2J\033[?25l")

	logs := make([]string, 3)
	for {
		for _, frame := range frames {
			// Print the next frame
			for _, line := range frame[minRow:maxRow] {
				for _, char := range line[minCol:maxCol] {
					fmt.Printf("\033[48;5;%sm%s", colors[string(char)], outputChar)
				}
				fmt.Println("\033[m")
			}

			for _, msg := range logs {
				fmt.Printf("\033[1;37;17m%s", msg)
				fmt.Println("\033[m")
			}

			// Reset the frame and sleep
			fmt.Print("\033[H")
			select {
			case msg := <-in:
				logs = append(logs[1:], msg)
			case <-ctx.Done():
				return
			default:
			}
			<-time.After(150 * time.Millisecond)
		}
	}
}

func termSize() (int, int) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	size, _ := cmd.Output()

	termWidth, _ := strconv.Atoi(strings.TrimSpace(strings.Split(string(size), " ")[1]))
	termHeight, _ := strconv.Atoi(strings.Split(string(size), " ")[0])
	return termWidth, termHeight - 3
}
