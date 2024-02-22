# Space Invaders Game in Golang

## Overview

This is a simple Space Invaders game implemented in Golang using the Pixel library. It features a player-controlled spaceship, alien invaders, and bullets for shooting down the aliens. The game window is created using PixelGL.

## Features

- Player can move the spaceship left and right using keyboard arrow keys.
- Player can shoot bullets to destroy incoming alien invaders using Spacebar.
- Alien invaders shoot back at the player.
- Score is tracked based on the number of aliens destroyed.

## Prerequisites

- Go 1.22 installed on your machine

## How to Run

1. Clone the repository:

    ```bash
    https://github.com/svidlak/go-space-invaders.git
    ```

2. Navigate to the project directory:

    ```bash
    cd go-space-invaders
    ```
    
3. Install dependencies

    ```bash
    go install
    ```
    
4. Run the game:

    ```bash
    go run main.go
    ```

## Controls

- Move left: Left arrow key
- Move right: Right arrow key
- Shoot: Spacebar

## Gameplay

- Avoid getting hit by alien bullets.
- Destroy as many aliens as possible to increase your score.

## Acknowledgments

This game is inspired by the classic Space Invaders arcade game. Special thanks to the Pixel library and the Golang community.

## License

This project is licensed under the [MIT License](LICENSE).

Feel free to contribute, report issues, or suggest improvements!