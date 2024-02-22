package constants

const (
	PlayerSpeed          = 300
	AlienMovementTimeout = 0.5
	AlienMovementOffset  = 5
)

type Direction string

type Entity string

const (
	PlayerEntity = "player"
	AlienEntity  = "alien"
)
const (
	Right Direction = "right"
	Left  Direction = "left"
	Down  Direction = "down"
)

const (
	WindowWidth  = 800
	WindowHeigth = 600
)

const (
	PlayerAsset = "/assets/player.png"
	AlienAsset  = "/assets/alien.png"
	BulletAsset = "/assets/bullet.png"
)
