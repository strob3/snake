package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

type pt struct{ x, y int }

func main() {
	s, err := tcell.NewScreen()
	if err != nil || s.Init() != nil {
		return
	}
	defer s.Fini()

	const w, h = 40, 20
	snake, dir, score, dead, food := []pt{{20, 10}, {19, 10}, {18, 10}}, pt{1, 0}, 0, false, pt{10, 5}

	spawnFood := func() {
		for hit := true; hit; {
			food, hit = pt{rand.Intn(w-2) + 1, rand.Intn(h-2) + 1}, false
			for _, b := range snake {
				if b == food {
					hit = true
				}
			}
		}
	}
	spawnFood()

	for range time.Tick(100 * time.Millisecond) {
		for s.HasPendingEvent() {
			if k, ok := s.PollEvent().(*tcell.EventKey); ok {
				switch {
				case k.Key() == tcell.KeyEscape || k.Rune() == 'q':
					return
				case k.Rune() == 'r':
					snake, dir, score, dead = []pt{{20, 10}, {19, 10}, {18, 10}}, pt{1, 0}, 0, false
					spawnFood()
				case (k.Key() == tcell.KeyUp || k.Rune() == 'w') && dir.y != 1:
					dir = pt{0, -1}
				case (k.Key() == tcell.KeyDown || k.Rune() == 's') && dir.y != -1:
					dir = pt{0, 1}
				case (k.Key() == tcell.KeyLeft || k.Rune() == 'a') && dir.x != 1:
					dir = pt{-1, 0}
				case (k.Key() == tcell.KeyRight || k.Rune() == 'd') && dir.x != -1:
					dir = pt{1, 0}
				}
			}
		}

		if !dead {
			head := pt{snake[0].x + dir.x, snake[0].y + dir.y}
			dead = head.x <= 0 || head.x >= w || head.y <= 0 || head.y >= h
			for _, b := range snake {
				dead = dead || head == b
			}
			if !dead {
				snake = append([]pt{head}, snake...)
				if head == food {
					score++
					spawnFood()
				} else {
					snake = snake[:len(snake)-1]
				}
			}
		}

		s.Clear()
		def := tcell.StyleDefault
		for x := 0; x <= w; x++ {
			s.SetContent(x, 0, '#', nil, def)
			s.SetContent(x, h, '#', nil, def)
		}
		for y := 0; y <= h; y++ {
			s.SetContent(0, y, '#', nil, def)
			s.SetContent(w, y, '#', nil, def)
		}
		s.SetContent(food.x, food.y, 'O', nil, def.Foreground(tcell.ColorRed))
		for _, b := range snake {
			s.SetContent(b.x, b.y, '█', nil, def.Foreground(tcell.ColorGreen))
		}
		msg := fmt.Sprintf(" score: %d | r: restart | q: quit", score)
		if dead {
			msg = fmt.Sprintf(" game over | score: %d | r: restart | q: quit", score)
		}
		for i, c := range msg {
			s.SetContent(i, h+1, c, nil, def)
		}
		s.Show()
	}
}
