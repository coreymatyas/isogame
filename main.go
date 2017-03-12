package main

import "github.com/gen2brain/raylib-go/raylib"

func main() {
	screenWidth := int32(1300)
	screenHeight := int32(700)

	tileWidth := int(34) // effective size after scaling by 0.5x, orthogonal
	tileHeight := int(34)
	tileDepth := int(82) - tileHeight // 82 (full height) - tileHeight

	worldWidth := int(64)
	worldHeight := int(64)
	worldDepth := int(3)

	drawOffsetX := screenWidth/2 - int32(tileWidth)
	drawOffsetY := screenHeight/2 - int32(tileHeight)

	raylib.InitWindow(screenWidth, screenHeight, "ISOGAME")
	defer raylib.CloseWindow()

	raylib.SetTargetFPS(60)

	camera := raylib.Camera2D{
		Target:   raylib.NewVector2(float32(screenWidth/2), float32(screenHeight/2)),
		Offset:   raylib.NewVector2(0, 0),
		Rotation: 0.0,
		Zoom:     0.5,
	}

	img := raylib.LoadImage("grass.png")
	defer raylib.UnloadImage(img)

	raylib.ImageResize(img, img.Width/2, img.Height/2)
	texture := raylib.LoadTextureFromImage(img)
	defer raylib.UnloadTexture(texture)

	world := Map{
		Name:    "World",
		Objects: make([]Object, 0),
		Tiles:   makeTiles(worldWidth, worldHeight, worldDepth),
	}

	for !raylib.WindowShouldClose() {
		if raylib.IsKeyDown(raylib.KeyRight) {
			camera.Offset.X -= 5 // Camera displacement with player movement
		} else if raylib.IsKeyDown(raylib.KeyLeft) {
			camera.Offset.X += 5 // Camera displacement with player movement
		}

		if raylib.IsKeyDown(raylib.KeyDown) {
			camera.Offset.Y -= 5 // Camera displacement with player movement
		} else if raylib.IsKeyDown(raylib.KeyUp) {
			camera.Offset.Y += 5 // Camera displacement with player movement
		}

		// Camera rotation controls
		if raylib.IsKeyDown(raylib.KeyA) {
			camera.Rotation--
		} else if raylib.IsKeyDown(raylib.KeyS) {
			camera.Rotation++
		}

		// Limit camera rotation to 80 degrees (-40 to 40)
		if camera.Rotation > 40 {
			camera.Rotation = 40
		} else if camera.Rotation < -40 {
			camera.Rotation = -40
		}

		// Camera zoom controls
		camera.Zoom += float32(raylib.GetMouseWheelMove()) * 0.05

		if camera.Zoom > 3.0 {
			camera.Zoom = 3.0
		} else if camera.Zoom < 0.1 {
			camera.Zoom = 0.1
		}

		// Camera reset (zoom and rotation)
		if raylib.IsKeyPressed(raylib.KeySpace) {
			camera.Zoom = 1.0
			camera.Rotation = 0.0
		}

		raylib.BeginDrawing()
		raylib.ClearBackground(raylib.RayWhite)
		raylib.Begin2dMode(camera)

		for ix := worldWidth - 1; ix >= 0; ix-- {
			for iy := 0; iy < worldHeight; iy++ {
				for iz := 0; iz < worldDepth; iz++ {
					itile := world.Tiles[ix][iy][iz]

					if itile.Enabled {
						v := raylib.Vector3{
							float32(ix * tileWidth),
							float32(iy * tileHeight),
							float32(iz * -tileDepth)}
						VectorISO(&v)
						raylib.DrawTexture(texture,
							int32(v.X)+drawOffsetX,
							int32(v.Y)+drawOffsetY,
							raylib.White)
					}
				}
			}
		}

		raylib.End2dMode()
		raylib.EndDrawing()
	}

}

type Map struct {
	Name    string
	Objects []Object
	Tiles   [][][]Tile
}

type Tile struct {
	Class   uint16
	Enabled bool
}

type Object struct {
	Position raylib.Vector3
}

func makeTiles(x, y, z int) [][][]Tile {
	tiles := make([][][]Tile, x)
	for i := 0; i < x; i++ {
		tiles[i] = make([][]Tile, y)
		for j := 0; j < y; j++ {
			tiles[i][j] = make([]Tile, z)

			// TODO: Procedurally generate some real way
			tiles[i][j][0] = Tile{
				Class:   0,
				Enabled: true,
			}

			if (i+j)%2 == 0 {
				tiles[i][j][1] = Tile{
					Class:   0,
					Enabled: true,
				}
			}
			if (i+j)%2 == 1 {
				tiles[i][j][2] = Tile{
					Class:   0,
					Enabled: true,
				}
			}
		}
	}
	return tiles
}
