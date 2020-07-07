package template

import (
	"fmt"
	"github.com/mattkimber/temporarily-late/internal/manifest"
)

const yjoggle = 4

// TTD bounding boxes

// diagonal = 4 pixels per unit (horizontal)
//          = 2 pixels per unit (vertical):
// 2x:
// 3 = 24 pixels
// 4 = 32 pixels
// 6 = 72 pixels

// horizontal = 8 pixels per unit (horizontal):
// 3 = 48 pixels

// vertical = 4 pixels per unit (vertical):
// 3 = 24 pixels


var unitShifts = map[int][]float64 {
	0: {0,-2},
	45: {2,-1},
	90: {4,0},
	135: {2,1},
	180: {0,2},
	225: {-2,1},
	270: {-4,0},
	315: {-2,-1},
}

var boundingBoxJoggles = map[int][]float64 {
	0: {1,0},
	45: {-3,1},
	90: {0,-1},
	135: {3,0.5},
	180: {1,0},
	225: {-3,1},
	270: {0,-1},
	315: {3,0.5},
}


// sprite top left is always 4 from front of relevant unit

// 8,8,8 doesn't shift on reversing
// 2,4,2 doesn't shift on reversing
// 3,4,3 shifts on reversing
// 2,6,1 shifts on reversing
// 2,5,1 shifts on reversing
// 3,6,1 shifts on reversing (2 units)
// 3,6,2 shifts on reversing
// 4,6,2 shifts on reversing
// 3,6,3 shifts on reversing
// 4,6,4 doesn't shift on reversing
// 4,4,4 doesn't shift on reversing
// 2,6,2 doesn't shift on reversing
// 2,5,2 doesn't shift on reversing
// 2,3,2 doesn't shift on reversing
// 4,2,4 doesn't shift on reversing

var configs = map[int][]int {
	3: {1,1,1},
	4: {1,1,2},
	5: {2,1,2},
	6: {2,2,2},
	7: {2,3,2},
	8: {2,4,2},
	9: {2,5,2},
	10: {4,2,4},
	11: {4,3,4},
	12: {4,4,4},
	13: {4,5,4},
	14: {4,6,4},
	15: {4,7,4},
	16: {4,8,4},
}

func WriteTemplates(m manifest.Manifest) {
	//scales := []int{1,2,4}
	lengths := []int{4,5,6,7,8,9,10,11,12,13,14,15,16}

	scales := []float64{1,2}

	for _, scale := range scales {
		for _, length := range lengths {
			WriteTemplate(scale, length, m)
		}
	}
}

func WriteTemplate(scale float64, length int, m manifest.Manifest) {
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

func produceTemplate(name string, angleOffset int, scale float64, length int, shiftAngle int, spritemap map[int]manifest.Sprite) {
	fmt.Printf("template template_auto%s_%d_%dx() {\n", name, length, int(scale))


	// Basic template
	for i := angleOffset; i < 360; i += 45 {
		direction := i - angleOffset
		spr := spritemap[(i+360+shiftAngle) % 360]

		x := float64(spr.X) * scale
		w := float64(spr.Width) * scale
		h := float64(spr.Height) * scale

		fscale := float64(scale)

		// Set xrel and yrel to the middle of the object
		xrel := -(w / 2)
		yrel := -(h / 2)
		yrel -= yjoggle * scale


		// joggle top left to the centre of the centre unit
		midSize := configs[length][1]
		diff := float64(8 - midSize) / 2

		xrel += unitShifts[direction][0] * diff * fscale
		yrel += unitShifts[direction][1] * diff * fscale

		// Get diagonal bounding boxes centred
		xrel += boundingBoxJoggles[direction][0] * fscale
		yrel += boundingBoxJoggles[direction][1] * fscale

		fmt.Printf("  [ %d, 0, %d, %d, %d, %d ]\n", int(x), int(w), int(h), int(xrel), int(yrel))
	}

	fmt.Printf("}\n\n")
}
