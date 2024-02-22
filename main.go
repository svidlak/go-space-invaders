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
	"github.com/svidlak/go-space-invaders/constants"
	"github.com/svidlak/go-space-invaders/models"
)

const (
	windowWidth  = 800
	windowHeigth = 600
)

const (
	playerAsset = "/assets/player.png"
	alienAsset  = "/assets/alien.png"
	bulletAsset = "/assets/bullet.png"
)

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

// func initBullet () {
// 	bulletImage := loadPicture(bulletAsset)
// }

func initPlayer(win *pixelgl.Window) func(float64) {
	playerImage := loadPicture(playerAsset)
	playerSprite := pixel.NewSprite(playerImage, playerImage.Bounds())
	playerSpriteWidthBound := (playerImage.Bounds().Max).X / 2

	player := models.Player{Position: pixel.V(windowWidth/2, 50)}

	return func(dt float64) {
		if win.Pressed(pixelgl.KeyLeft) {
			if player.Position.X > playerSpriteWidthBound {
				player.Position.X -= constants.PlayerSpeed * dt
			}
		}
		if win.Pressed(pixelgl.KeyRight) {
			if player.Position.X < windowWidth-playerSpriteWidthBound {
				player.Position.X += constants.PlayerSpeed * dt
			}
		}

		playerSprite.Draw(win, pixel.IM.Moved(player.Position))
	}
}

func initAliens(win *pixelgl.Window) func(float64) {
	alienImage := loadPicture(alienAsset)
	alienSprite := pixel.NewSprite(alienImage, alienImage.Bounds())

	alienSpriteWidthBound := (alienImage.Bounds().Max).X / 2

	last := time.Now()

	alien := models.Alien{Direction: constants.Down, Position: pixel.V(alienSpriteWidthBound, windowHeigth-50)}
	aliens := []models.Alien{alien}

	return func(dt float64) {
		passed := time.Since(last).Seconds()

		if passed > constants.AlienMovementTimeout {
			for idx, alienVal := range aliens {
				alien := &alienVal

				alien.UpdateDirection(alienSpriteWidthBound, windowWidth)
				alien.Move((alienImage.Bounds().Max).Y + constants.AlienMovementTimeout)

				aliens[idx] = *alien
			}

			aliens = append(aliens, alien)
			last = time.Now()
		}

		for _, alien := range aliens {
			alienSprite.Draw(win, pixel.IM.Moved(alien.Position))
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
