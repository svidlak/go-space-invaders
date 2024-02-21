package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"os"
	"time"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

var windowProperties = map[string]float64{
	"width":  800,
	"height": 600,
}

var assetsPath = map[string]string{
	"player": "/assets/player.png",
	"alien":  "/assets/alien.png",
}

type Direction string

const (
	Right Direction = "right"
	Left  Direction = "left"
	Down  Direction = "down"
)

const playerSpeed = 300

type Alien struct {
	direction Direction
	position  pixel.Vec
}

func (a *Alien) setPosition(heigth float64) {
	if a.direction == Right {
		a.position.X += playerSpeed / 5
	} else if a.direction == Down {
		a.position.Y -= heigth
	} else {
		a.position.X -= playerSpeed / 5
	}
}

func loadPicture(path string) pixel.Picture {
	pwd, _ := os.Getwd()
	file, err := os.Open(pwd + path)

	if err != nil {
		fmt.Printf("failed to open picture file: %v", err)
		panic(err)
	}

	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Printf("failed to decode picture: %v", err)
		panic(err)
	}

	return pixel.PictureDataFromImage(img)
}

func initGameWindow() *pixelgl.Window {
	cfg := pixelgl.WindowConfig{
		Title:  "Svidlak Space Invaders",
		Bounds: pixel.R(0, 0, windowProperties["width"], windowProperties["height"]),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	return win
}

func initPlayer(win *pixelgl.Window) func(float64) {
	playerImage := loadPicture(assetsPath["player"])
	playerSprite := pixel.NewSprite(playerImage, playerImage.Bounds())
	playerSpriteWidthBound := (playerImage.Bounds().Max).X / 2
	playerPos := pixel.V(windowProperties["width"]/2, 50)

	return func(dt float64) {
		if win.Pressed(pixelgl.KeyLeft) {
			if playerPos.X > playerSpriteWidthBound {
				playerPos.X -= playerSpeed * dt
			}
		}
		if win.Pressed(pixelgl.KeyRight) {
			if playerPos.X < windowProperties["width"]-playerSpriteWidthBound {
				playerPos.X += playerSpeed * dt
			}
		}

		playerSprite.Draw(win, pixel.IM.Moved(playerPos))
	}
}

func initAliens(win *pixelgl.Window) func(float64) {
	alienImage := loadPicture(assetsPath["alien"])
	alienSprite := pixel.NewSprite(alienImage, alienImage.Bounds())

	alienSpriteWidthBound := (alienImage.Bounds().Max).X / 2
	alienSpriteHeigth := (alienImage.Bounds().Max).Y

	alienPos := pixel.V(alienSpriteWidthBound, windowProperties["height"]-50)

	last := time.Now()

	alien := Alien{direction: Down, position: alienPos}
	aliens := []Alien{alien}

	return func(dt float64) {
		passed := time.Since(last).Seconds()

		if passed > 0.5 {
			for idx, alienVal := range aliens {
				alien := &alienVal

				if alien.position.X <= alienSpriteWidthBound {
					if alien.direction == Down {
						alien.direction = Right
					} else {
						alien.direction = Down
					}
				}

				if alien.position.X >= windowProperties["width"]-((alienSpriteWidthBound*2)*2) {
					if alien.direction == Down {
						alien.direction = Left
					} else {
						alien.direction = Down
					}
				}

				alien.setPosition(alienSpriteHeigth + 5)
				aliens[idx] = *alien
			}

			aliens = append(aliens, alien)
			last = time.Now()
		}

		for _, alien := range aliens {
			alienSprite.Draw(win, pixel.IM.Moved(alien.position))
		}
	}
}

func run() {
	win := initGameWindow()
	updatePlayerMovement := initPlayer(win)
	updateAlienMovement := initAliens(win)

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		win.Clear(color.Black)

		updatePlayerMovement(dt)
		updateAlienMovement(dt)

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
