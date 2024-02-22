package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math/rand"
	"os"
	"time"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"github.com/svidlak/go-space-invaders/constants"
	"github.com/svidlak/go-space-invaders/models"
)

var gameState = models.GameState{}

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
		Bounds: pixel.R(0, 0, constants.WindowWidth, constants.WindowHeigth),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	return win
}

func detectPlayerCollision(bullet models.Bullet) bool {
	return gameState.Player.DetectCollision(bullet)
}

func detectAlienCollision(bullet models.Bullet) bool {
	collision := false

	tmpAliens := []models.Alien{}

	for _, alien := range gameState.Aliens {
		if alien.DetectCollision(bullet) {
			collision = true
		} else {
			tmpAliens = append(tmpAliens, alien)
		}
	}

	gameState.Aliens = tmpAliens
	return collision
}

func initBullet(win *pixelgl.Window) (func(pixel.Vec, constants.Entity), func()) {
	bulletImage := loadPicture(constants.BulletAsset)
	bulletSprite := pixel.NewSprite(bulletImage, bulletImage.Bounds())

	bullets := []models.Bullet{}

	return func(position pixel.Vec, entity constants.Entity) {
			bullets = append(bullets, models.Bullet{Position: position, Entity: entity})
		}, func() {
			tmpBullets := []models.Bullet{}

			for idx, bullet := range bullets {
				bulletCollision := false

				if bullet.Entity == constants.AlienEntity {
					bullet.Position.Y--
					bulletCollision = detectPlayerCollision(bullet)
					if bulletCollision {
						fmt.Println("player fallen")
					}
				} else {
					bullet.Position.Y++
					bulletCollision = detectAlienCollision(bullet)

					if bulletCollision {
						fmt.Println("alien fallen")
					}
				}

				bullets[idx] = bullet

				if bullet.Position.Y > 0 && bullet.Position.Y < constants.WindowHeigth && !bulletCollision {
					bulletSprite.Draw(win, pixel.IM.Moved(bullet.Position))
					tmpBullets = append(tmpBullets, bullet)
				}
			}

			bullets = tmpBullets
		}
}

func initPlayer(win *pixelgl.Window, shootBullet func(pixel.Vec, constants.Entity)) func(float64) {
	playerImage := loadPicture(constants.PlayerAsset)
	playerSprite := pixel.NewSprite(playerImage, playerImage.Bounds())
	playerSpriteWidthBound := (playerImage.Bounds().Max).X / 2

	gameState.Player = models.Player{Bounds: playerImage.Bounds(), Position: pixel.V(constants.WindowWidth/2, 50)}

	return func(dt float64) {
		if win.Pressed(pixelgl.KeyLeft) {
			if gameState.Player.Position.X > playerSpriteWidthBound {
				gameState.Player.Position.X -= constants.PlayerSpeed * dt
			}
		}
		if win.Pressed(pixelgl.KeyRight) {
			if gameState.Player.Position.X < constants.WindowWidth-playerSpriteWidthBound {
				gameState.Player.Position.X += constants.PlayerSpeed * dt
			}
		}

		if win.JustPressed(pixelgl.KeySpace) {
			shootBullet(gameState.Player.Position, constants.PlayerEntity)
		}

		playerSprite.Draw(win, pixel.IM.Moved(gameState.Player.Position))
	}
}

func initAliens(win *pixelgl.Window, shootBullet func(pixel.Vec, constants.Entity)) func(float64) {
	alienImage := loadPicture(constants.AlienAsset)
	alienSprite := pixel.NewSprite(alienImage, alienImage.Bounds())

	alienSpriteWidthBound := (alienImage.Bounds().Max).X / 2

	last := time.Now()

	alien := models.Alien{Bounds: alienImage.Bounds(), Direction: constants.Down, Position: pixel.V(alienSpriteWidthBound, constants.WindowHeigth-50)}
	gameState.Aliens = []models.Alien{alien}

	return func(dt float64) {
		passed := time.Since(last).Seconds()

		if passed > constants.AlienMovementTimeout {
			for idx, alienVal := range gameState.Aliens {
				alien := &alienVal

				alien.UpdateDirection(alienSpriteWidthBound, constants.WindowWidth)
				alien.Move((alienImage.Bounds().Max).Y + constants.AlienMovementTimeout)

				gameState.Aliens[idx] = *alien
			}

			gameState.Aliens = append(gameState.Aliens, alien)
			randomAlienNumber := rand.Intn(len(gameState.Aliens) - 0)
			shootBullet(gameState.Aliens[randomAlienNumber].Position, constants.AlienEntity)
			last = time.Now()
		}

		for _, alien := range gameState.Aliens {
			alienSprite.Draw(win, pixel.IM.Moved(alien.Position))
		}
	}
}

func run() {
	win := initGameWindow()
	shootBullet, updatePlayerBullets := initBullet(win)
	updatePlayerMovement := initPlayer(win, shootBullet)
	updateAlienMovement := initAliens(win, shootBullet)

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		win.Clear(color.Black)

		updatePlayerMovement(dt)
		updateAlienMovement(dt)
		updatePlayerBullets()

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
