package main

import (
	"bytes"
	"fmt"
	"os"
	"time"
)

const (
	WALL        = 1
	SPACE       = 0
	MAX_SAMPLES = 100
	PLAYER      = 69
)

type Stats struct {
	start  time.Time
	frames int
	fps    float64
}

func newStats() *Stats {
	return &Stats{
		start: time.Now(),
	}
}

func (s *Stats) update() {
	s.frames++
	if s.frames == MAX_SAMPLES {
		s.fps = float64(s.frames) / time.Since(s.start).Seconds()
		s.frames = 0
		s.start = time.Now()
	}
}

type Level struct {
	width, height int
	data          [][]byte
}

func newLevel(width, height int) *Level {

	data := make([][]byte, width)
	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			data[h] = make([]byte, width)
		}
	}

	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			data[h][w] = byte(SPACE)
		}
	}

	for h := 0; h < height; h++ {
		data[h][0] = WALL
		data[h][width-1] = WALL
	}
	for w := 0; w < width; w++ {
		data[0][w] = WALL
		data[height-1][w] = WALL
	}

	return &Level{
		width:  width,
		height: height,
		data:   data,
	}
}

type Game struct {
	isRunning  bool
	level      *Level
	stats      *Stats
	drawBuffer *bytes.Buffer
}

func newGame(width, height int) *Game {
	level := newLevel(width, height)
	return &Game{isRunning: false, level: level, stats: newStats(), drawBuffer: new(bytes.Buffer)}
}

func (g *Game) start() {
	g.isRunning = true
	g.loop()
}

func (g *Game) update() {

}

func (g *Game) renderLevel() {
	lvl := g.level
	for h := 0; h < lvl.height; h++ {
		for w := 0; w < lvl.width; w++ {
			if lvl.data[h][w] == WALL {
				g.drawBuffer.WriteString(fmt.Sprintf("%c\t", 0x25A0))
			} else {
				g.drawBuffer.WriteString("\t")
			}
		}
		g.drawBuffer.WriteString("\n")
	}
}

func (g *Game) render() {
	g.drawBuffer.Reset()
	_, _ = fmt.Fprint(os.Stdout, "\033c")
	g.renderLevel()
	g.renderStats()
	g.stats.update()
	_, _ = fmt.Fprint(os.Stdout, g.drawBuffer.String())
}

func (g *Game) renderStats() {
	g.drawBuffer.WriteString("-- STATS\n")
	g.drawBuffer.WriteString(fmt.Sprintf("FPS: %.2f\n", g.stats.fps))
}

func (g *Game) loop() {
	for g.isRunning {
		g.update()
		g.render()
		time.Sleep(16 * time.Millisecond)
	}
}

func main() {
	newGame(10, 10).start()
}
