package template

import (
	"fmt"
	"github.com/mattkimber/temporarily-late/internal/manifest"
)

const yjoggle = 4

var joggles = map[int][]int {
	0: {1,0},
	45: {0,0},
	90: {6,0},
	135: {6,3},
	180: {1,6},
	225: {-6,3},
	270: {-6,0},
	315: {0,0},
}

var autoSizeJoggles = map[int][]int {
	0: {0,0},
	45: {0,0},
	90: {2,0},
	135: {2,1},
	180: {0,2},
	225: {-2,1},
	270: {-2,0},
	315: {0,0},
}

func WriteTemplates(m manifest.Manifest) {
	//scales := []int{1,2,4}
	lengths := []int{4,5,6,7,8,9,10,11,12,13,14,15,16}

	scales := []int{1,2}

	for _, scale := range scales {
		for _, length := range lengths {
			WriteTemplate(scale, length, m)
		}
	}
}

func WriteTemplate(scale int, length int, m manifest.Manifest) {
	spritemap := make(map[int]manifest.Sprite)
	tx := 0
	for _, spr := range m.Sprites {
		spr.X = tx
		spritemap[int(spr.Angle)] = spr
		tx += spr.Width + 8 // GoRender sprite spacing
	}

	produceTemplate("", 0, scale, length, 0, spritemap)

	produceTemplate("_turn_1", 15, scale, length, 0, spritemap)
	produceTemplate("_turn_2", 30, scale, length, -45, spritemap)

	//produceTemplate("_turn_r1", 15, scale, length, 0, spritemap)
	//produceTemplate("_turn_r2", 30, scale, length, -45, spritemap)
}

func produceTemplate(name string, angleOffset int, scale int, length int, shiftAngle int, spritemap map[int]manifest.Sprite) {
	fmt.Printf("template template_auto%s_%d_%dx() {\n", name, length, scale)

	// Basic template
	for i := angleOffset; i < 360; i += 45 {
		spr := spritemap[(i+360+shiftAngle) % 360]

		x := spr.X * scale
		w := spr.Width * scale
		h := spr.Height * scale

		xrel := -(w / 2)
		yrel := -(h / 2)

		yrel -= yjoggle * scale

		joggle := joggles[i - angleOffset]
		xrel += joggle[0] * scale
		yrel += joggle[1] * scale

		jogglingSize := (12 - length) / 2

		// Special cases
		if length == 5 {
			jogglingSize = 2
		}

		if length == 4 {
			jogglingSize = 3
		}

		if length == 12 {
			jogglingSize = -1
		}

		if length == 13 {
			jogglingSize = 0
		}


		autoSizeJoggle := autoSizeJoggles[i - angleOffset]
		xrel += autoSizeJoggle[0] * jogglingSize * scale
		yrel += autoSizeJoggle[1] * jogglingSize * scale


		// Depot hack
		if i - angleOffset == 90 || i - angleOffset == 270 {
			xrel += 2
		}

		// Tile adjustment as per https://newgrf-specs.tt-wiki.net/wiki/RealSprites
		if scale == 2 {
			yrel -= 1
		}

		if scale == 4 {
			yrel -= 2
		}

		fmt.Printf("  [ %d, 0, %d, %d, %d, %d ]\n", x, w, h, xrel, yrel)
	}

	fmt.Printf("}\n\n")
}
