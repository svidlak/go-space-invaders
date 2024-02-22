package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"os"
	"strconv"
	"time"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"github.com/gopxl/pixel/text"
	"github.com/svidlak/go-space-invaders/constants"
	"github.com/svidlak/go-space-invaders/models"
	"golang.org/x/image/font/basicfont"
)

var gameStateController = models.GameStateController{}

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

func detectPlayerCollision(bullet *models.Bullet) bool {
	return gameStateController.Player.DetectCollision(bullet)
}

func detectAlienCollision(bullet *models.Bullet) bool {
	collision := false

	tmpAliens := []models.Alien{}

	for _, alien := range gameStateController.Aliens {
		if alien.DetectCollision(bullet) {
			collision = true
			gameStateController.Score += 10
		} else {
			tmpAliens = append(tmpAliens, alien)
		}
	}

	gameStateController.Aliens = tmpAliens
	return collision
}

func detectCollision(bullet *models.Bullet) bool {
	if bullet.Entity == constants.AlienEntity {
		return detectPlayerCollision(bullet)
	} else {
		return detectAlienCollision(bullet)
	}
}

func initBullet(win *pixelgl.Window) func() {
	bulletImage := loadPicture(constants.BulletAsset)
	bulletSprite := pixel.NewSprite(bulletImage, bulletImage.Bounds())

	gameStateController.Bullets = []models.Bullet{}

	return func() {
		tmpBullets := []models.Bullet{}

		for _, bullet := range gameStateController.Bullets {
			bullet.MoveBullet(&bullet)
			bulletCollision := detectCollision(&bullet)

			if bullet.Position.Y > 0 && bullet.Position.Y < constants.WindowHeigth && !bulletCollision {
				bulletSprite.Draw(win, pixel.IM.Moved(bullet.Position))
				tmpBullets = append(tmpBullets, bullet)
			}
		}

		gameStateController.Bullets = tmpBullets
	}
}

func initPlayer(win *pixelgl.Window) func(float64) {
	playerImage := loadPicture(constants.PlayerAsset)
	playerSprite := pixel.NewSprite(playerImage, playerImage.Bounds())
	playerSpriteWidthBound := (playerImage.Bounds().Max).X / 2

	gameStateController.Player = models.Player{Bounds: playerImage.Bounds(), Position: pixel.V(constants.WindowWidth/2, 50), Health: 3}

	return func(dt float64) {
		if win.Pressed(pixelgl.KeyLeft) {
			if gameStateController.Player.Position.X > playerSpriteWidthBound {
				gameStateController.Player.Position.X -= constants.PlayerSpeed * dt
			}
		}
		if win.Pressed(pixelgl.KeyRight) {
			if gameStateController.Player.Position.X < constants.WindowWidth-playerSpriteWidthBound {
				gameStateController.Player.Position.X += constants.PlayerSpeed * dt
			}
		}

		if win.JustPressed(pixelgl.KeySpace) {
			gameStateController.ShootBullet(constants.PlayerEntity)
		}

		playerSprite.Draw(win, pixel.IM.Moved(gameStateController.Player.Position))
	}
}

func initAliens(win *pixelgl.Window) func(float64) {
	alienImage := loadPicture(constants.AlienAsset)
	alienSprite := pixel.NewSprite(alienImage, alienImage.Bounds())

	alienSpriteWidthBound := (alienImage.Bounds().Max).X / 2
	alienMovement := (alienImage.Bounds().Max).Y + constants.AlienMovementTimeout
	alienPosition := pixel.V(alienSpriteWidthBound, constants.WindowHeigth-50)

	last := time.Now()

	gameStateController.AddAlien(alienImage.Bounds(), constants.Down, alienPosition)

	return func(dt float64) {
		passed := time.Since(last).Seconds()

		if passed > constants.AlienMovementTimeout {
			gameStateController.UpdateAlienMovement(dt, alienSpriteWidthBound, alienMovement)
			gameStateController.AddAlien(alienImage.Bounds(), constants.Down, alienPosition)
			gameStateController.ShootBullet(constants.AlienEntity)

			last = time.Now()
		}

		for _, alien := range gameStateController.Aliens {
			alienSprite.Draw(win, pixel.IM.Moved(alien.Position))
		}
	}
}

func initGameText(win *pixelgl.Window) func() {
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	return func() {
		scoreText := text.New(pixel.V(0, constants.WindowHeigth-25), basicAtlas)
		healthText := text.New(pixel.V(constants.WindowWidth-100, constants.WindowHeigth-25), basicAtlas)

		fmt.Fprintln(scoreText, "Score: "+strconv.Itoa(gameStateController.Score))
		fmt.Fprintln(healthText, "Health: "+strconv.Itoa(gameStateController.Player.Health))

		scoreText.Draw(win, pixel.IM.Scaled(scoreText.Orig, 1.5))
		healthText.Draw(win, pixel.IM.Scaled(healthText.Orig, 1.5))
	}
}

func run() {
	win := initGameWindow()
	updatePlayerMovement := initPlayer(win)
	updateAlienMovement := initAliens(win)
	updateBulletsMovement := initBullet(win)
	updateScoreText := initGameText(win)

	win.SetSmooth(true)

	last := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		win.Clear(color.Black)

		updatePlayerMovement(dt)
		updateAlienMovement(dt)
		updateBulletsMovement()
		updateScoreText()

		if gameStateController.Player.Health < 1 {
			win.Destroy()
		}

		win.Update()

	}
}

func main() {
	pixelgl.Run(run)
}
