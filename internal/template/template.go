package template

import (
	"fmt"
	"github.com/mattkimber/temporarily-late/internal/manifest"
)

const yjoggle = 4
const road_y_joggle = 3

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

var unitShifts = map[int][]float64{
	0:   {0, -2},
	45:  {2, -1},
	90:  {4, 0},
	135: {2, 1},
	180: {0, 2},
	225: {-2, 1},
	270: {-4, 0},
	315: {-2, -1},
}

var boundingBoxJoggles = map[int][]float64{
	0:   {1, 0}, // vertical, heading N
	45:  {-3, 1}, // diagonal, heading NE
	90:  {0, -1}, // horizontal, heading E
	135: {3, 0.5},  // diagonal, heading SE
	180: {1, 0},
	225: {-3, 0.5},
	270: {0, -1},
	315: {3, 0.5},
}

var roadBoundingBoxJoggles = map[int][]float64{
	0:   {1, 1}, // vertical, heading N
	45:  {-2.5, 2}, // diagonal, heading NE
	90:  {0, 1}, // horizontal, heading E
	135: {2.5, 0.5}, // diagonal, heading SE
	180: {0, 1},
	225: {-2, 3},
	270: {0, 1},
	315: {3, 2},
}

var uphillJoggles = map[int]float64 {
	45: -1,
	135: 0,
	225: 0,
	315: -1,
}

var downhillJoggles = map[int]float64 {
	45: 0,
	135: -1,
	225: -1,
	315: 0,
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

var configs = map[int][]int{
	3:  {1, 1, 1},
	4:  {1, 1, 2},
	5:  {2, 1, 2},
	6:  {2, 2, 2},
	7:  {2, 3, 2},
	8:  {2, 4, 2},
	9:  {2, 5, 2},
	10: {4, 2, 4},
	11: {4, 3, 4},
	12: {4, 4, 4},
	13: {4, 5, 4},
	14: {4, 6, 4},
	15: {4, 7, 4},
	16: {4, 8, 4},
}

func WriteTemplates(m manifest.Manifest, addAnimatedFlag bool) {
	//scales := []int{1,2,4}
	lengths := []int{4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

	scales := []float64{1, 2}

	for _, scale := range scales {
		for _, length := range lengths {
			WriteTemplate(scale, length, m, addAnimatedFlag)
		}
	}
}

func WriteRoadTemplates(m manifest.Manifest, addAnimatedFlag bool) {
	lengths := []int{1,2,3,4,5,6,7,8}
	scales := []float64{1,2}

	for _, scale := range scales {
		for _, length := range lengths {
			WriteRoadTemplate(scale, length, m, addAnimatedFlag)
		}
	}
}

func WriteRoadTemplate(scale float64, length int, m manifest.Manifest, addAnimatedFlag bool) {
	spritemap := make(map[int]manifest.Sprite)
	tx := 0

	for _, spr := range m.Sprites {
		spr.X = tx
		spritemap[int(spr.Angle)] = spr
		tx += spr.Width + 8 // GoRender sprite spacing
	}

	produceFlatRoadTemplate("", scale, length, spritemap, addAnimatedFlag)
}

func produceFlatRoadTemplate(name string,scale float64, length int, spritemap map[int]manifest.Sprite, addAnimatedFlag bool) {
	suffix := ""
	if addAnimatedFlag {
		suffix = "_anim"
	}
	fmt.Printf("template template_auto%s_%d%s_%dx() {\n", name, length, suffix, int(scale))

	// Basic template
	for i := 0; i < 360; i += 45 {
		direction := i
		spr := spritemap[(i+360)%360]

		x := float64(spr.X) * scale
		w := float64(spr.Width) * scale
		h := float64(spr.Height) * scale

		fscale := scale

		xrel, yrel := getRoadRels(w, h, scale, length, direction, fscale)

		fmt.Printf("  [ %d, 0, %d, %d, %d, %d ]\n", int(x), int(w), int(h), int(xrel), int(yrel))
	}

	fmt.Printf("}\n\n")
}


func getRoadRels(w float64, h float64, scale float64, length int, direction int, fscale float64) (xrel float64, yrel float64) {
	// Set xrel and yrel to the middle of the object
	xrel = -(w / 2)
	yrel = -(h / 2)
	yrel -= road_y_joggle * scale

	// joggle top left to the centre of the unit
	diff := float64(8-length) / 2

	// Special handling for L4 vehicles
	// not needed for road?
	//if length == 4 {
	//	diff -= 0.5
	//}

	xrel += unitShifts[direction][0] * diff * fscale
	yrel += unitShifts[direction][1] * diff * fscale

	// Get diagonal bounding boxes centred
	xrel += roadBoundingBoxJoggles[direction][0] * fscale
	yrel += roadBoundingBoxJoggles[direction][1] * fscale
	return xrel, yrel
}


func WriteTemplate(scale float64, length int, m manifest.Manifest, addAnimatedFlag bool) {
	spritemap := make(map[int]manifest.Sprite)
	tx := 0

	// This is a bit of a hack for slice length files
	if m.SliceLength > 0 {
		m.Sprites = m.Sprites[0:8]
	}

	for _, spr := range m.Sprites {
		spr.X = tx
		spritemap[int(spr.Angle)] = spr
		tx += spr.Width + 8 // GoRender sprite spacing
	}

	if m.IsHill {
		produceHillTemplate("_up", scale, length, uphillJoggles, spritemap, addAnimatedFlag)
		produceHillTemplate("_down", scale, length, downhillJoggles, spritemap, addAnimatedFlag)
	} else if m.SliceLength > 0 && length > 8 {
		produceFlatTemplate("_front", 0, scale, length, 0, 0, -1, spritemap, addAnimatedFlag)
		produceFlatTemplate("_mid", 0, scale, length, 0, tx, 0, spritemap, addAnimatedFlag)
		produceFlatTemplate("_rear", 0, scale, length, 0, tx*2, 1, spritemap, addAnimatedFlag)
	} else if m.SliceLength == 0 {
		produceFlatTemplate("", 0, scale, length, 0, 0, 0, spritemap, addAnimatedFlag)
		produceFlatTemplate("_turn_1", 15, scale, length, 0, 0,0, spritemap, addAnimatedFlag)
		produceFlatTemplate("_turn_2", 30, scale, length, -45, 0, 0, spritemap, addAnimatedFlag)
	}
}

func produceFlatTemplate(name string, angleOffset int, scale float64, length int, shiftAngle int, offsetWithinFile int, unitOffset int, spritemap map[int]manifest.Sprite, addAnimatedFlag bool) {
	suffix := ""
	if addAnimatedFlag {
		suffix = "_anim"
	}
	fmt.Printf("template template_auto%s_%d%s_%dx() {\n", name, length,  suffix, int(scale))

	// Basic template
	for i := angleOffset; i < 360; i += 45 {
		direction := i - angleOffset
		spr := spritemap[(i+360+shiftAngle)%360]

		x := float64(spr.X + offsetWithinFile) * scale
		w := float64(spr.Width) * scale
		h := float64(spr.Height) * scale

		fscale := scale

		xrel, yrel := getRels(w, h, scale, length, direction, fscale, unitOffset, offsetWithinFile)

		flag := ""
		if addAnimatedFlag {
			flag = ", ANIM"
		}
		fmt.Printf("  [ %d, 0, %d, %d, %d, %d%s ]\n", int(x), int(w), int(h), int(xrel), int(yrel), flag)
	}

	fmt.Printf("}\n\n")
}

func produceHillTemplate(name string, scale float64, length int, hillJoggles map[int]float64, spritemap map[int]manifest.Sprite, addAnimatedFlag bool) {
	suffix := ""
	if addAnimatedFlag {
		suffix = "_anim"
	}
	fmt.Printf("template template_auto%s_%d%s_%dx() {\n", name, length, suffix, int(scale))

	// Basic template
	for i := 0; i < 360; i += 45 {
		if i == 45 || i == 135 || i == 225 || i == 315 {
			direction := i
			spr := spritemap[i]

			x := float64(spr.X) * scale
			w := float64(spr.Width) * scale
			h := float64(spr.Height) * scale

			fscale := scale

			xrel, yrel := getRels(w, h, scale, length, direction, fscale, 0, 0)

			yrel += 0.25 * float64(length) * scale
			yrel += hillJoggles[i] * scale

			flag := ""
			if addAnimatedFlag {
				flag = ", ANIM"
			}
			fmt.Printf("  [ %d, 0, %d, %d, %d, %d%s ]\n", int(x), int(w), int(h), int(xrel), int(yrel), flag)
		} else {
			fmt.Printf("  [ 0, 0, 1, 1, 0, 0 ]\n")
		}
	}

	fmt.Printf("}\n\n")
}

func getRels(w float64, h float64, scale float64, length int, direction int, fscale float64, unitOffset int, offsetWithinFile int) (xrel float64, yrel float64) {
	// Set xrel and yrel to the middle of the object
	xrel = -(w / 2)
	yrel = -(h / 2)
	yrel -= yjoggle * scale

	// joggle top left to the centre of the centre unit
	midSize := configs[length][1]
	diff := float64(8-midSize) / 2

	// Joggle backward (or forward) for front and rear sections.
	if unitOffset == -1 {
		sectionSize := configs[length][0]
		sectionDiff := float64(sectionSize)
		diff -= sectionDiff
	}

	if unitOffset == 1 {
		sectionSize := configs[length][1]
		sectionDiff := float64(sectionSize)
		diff += sectionDiff
	}

	// Special handling for L4 vehicles
	if length == 4 {
		diff -= 0.5
	}

	xrel += unitShifts[direction][0] * diff * fscale
	yrel += unitShifts[direction][1] * diff * fscale

	// Get diagonal bounding boxes centred
	xrel += boundingBoxJoggles[direction][0] * fscale
	yrel += boundingBoxJoggles[direction][1] * fscale
	return xrel, yrel
}
