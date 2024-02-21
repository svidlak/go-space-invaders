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

const (
	windowWidth  = 800
	windowHeigth = 600
)

const (
	playerSpeed          = 300
	alienMovementTimeout = 0.5
	alienMovementOffset  = 5
)

const (
	playerAsset = "/assets/player.png"
	alienAsset  = "/assets/alien.png"
)

type Direction string

const (
	Right Direction = "right"
	Left  Direction = "left"
	Down  Direction = "down"
)

type Alien struct {
	direction Direction
	position  pixel.Vec
}

func (a *Alien) move(height float64) {
	if a.direction == Right {
		a.position.X += playerSpeed / alienMovementOffset
	} else if a.direction == Down {
		a.position.Y -= height
	} else {
		a.position.X -= playerSpeed / alienMovementOffset
	}
}

func (a *Alien) updateDirection(alienSpriteWidthBound, windowWidth float64) {
	if a.position.X <= alienSpriteWidthBound {
		if a.direction == Down {
			a.direction = Right
		} else {
			a.direction = Down
		}
	}

	if a.position.X >= windowWidth-((alienSpriteWidthBound*2)*2) {
		if a.direction == Down {
			a.direction = Left
		} else {
			a.direction = Down
		}
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
		Bounds: pixel.R(0, 0, windowWidth, windowHeigth),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	return win
}

func initPlayer(win *pixelgl.Window) func(float64) {
	playerImage := loadPicture(playerAsset)
	playerSprite := pixel.NewSprite(playerImage, playerImage.Bounds())
	playerSpriteWidthBound := (playerImage.Bounds().Max).X / 2
	playerPos := pixel.V(windowWidth/2, 50)

	return func(dt float64) {
		if win.Pressed(pixelgl.KeyLeft) {
			if playerPos.X > playerSpriteWidthBound {
				playerPos.X -= playerSpeed * dt
			}
		}
		if win.Pressed(pixelgl.KeyRight) {
			if playerPos.X < windowWidth-playerSpriteWidthBound {
				playerPos.X += playerSpeed * dt
			}
		}

		playerSprite.Draw(win, pixel.IM.Moved(playerPos))
	}
}

func initAliens(win *pixelgl.Window) func(float64) {
	alienImage := loadPicture(alienAsset)
	alienSprite := pixel.NewSprite(alienImage, alienImage.Bounds())

	alienSpriteWidthBound := (alienImage.Bounds().Max).X / 2

	alienPos := pixel.V(alienSpriteWidthBound, windowHeigth-50)

	last := time.Now()

	alien := Alien{direction: Down, position: alienPos}
	aliens := []Alien{alien}

	return func(dt float64) {
		passed := time.Since(last).Seconds()

		if passed > alienMovementTimeout {
			for idx, alienVal := range aliens {
				alien := &alienVal

				alien.updateDirection(alienSpriteWidthBound, windowWidth)
				alien.move((alienImage.Bounds().Max).Y + alienMovementOffset)

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
