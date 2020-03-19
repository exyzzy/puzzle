package main

import (
	"fmt"
	"math"
	"time"
)

//Unique patterns to match, each bisected side
const (
	lBigBee     = iota //Left of Big Bee
	rBigBee            //Right of Big Bee
	fBigBee            //Front of Big Bee
	bBigBee            //Back of Big Bee
	fLittleBees        //Front of Little Bees
	bLittleBees        //Back of Little Bees
	fStripeBfly        //Front of Stripe Butterfly
	bStripeBfly        //Back of Stripe Butterfly
	lBlueBfly          //Left of Blue Butterfly
	rBlueBfly          //Right of Blue Butterfly
	fPinkBfly          //Front of Pink Butterfly
	bPinkBfly          //Back of Pink Butterfly
)

//Each pattern can extend onto the border
//so need to initialize the outside puzzle border
var border = [12]int{bBigBee, bStripeBfly, fLittleBees, bStripeBfly, lBlueBfly, lBigBee, fPinkBfly, bLittleBees, bPinkBfly, rBlueBfly, fBigBee, fPinkBfly}

//Traingular puzzle pieces
var pieces = [16][3]int{
	{fLittleBees, fPinkBfly, lBlueBfly},     //0
	{fBigBee, lBlueBfly, bPinkBfly},         //1
	{bPinkBfly, fBigBee, lBigBee},           //2
	{bBigBee, fPinkBfly, rBigBee},           //3
	{bStripeBfly, bStripeBfly, rBigBee},     //4
	{fBigBee, bPinkBfly, bPinkBfly},         //5
	{fStripeBfly, rBigBee, bBigBee},         //6
	{lBlueBfly, rBlueBfly, lBigBee},         //7
	{fStripeBfly, rBlueBfly, lBigBee},       //8
	{bPinkBfly, fPinkBfly, lBlueBfly},       //9
	{fStripeBfly, rBlueBfly, lBlueBfly},     //10
	{fStripeBfly, fStripeBfly, rBlueBfly},   //11
	{rBigBee, bLittleBees, rBlueBfly},       //12
	{fPinkBfly, bPinkBfly, fLittleBees},     //13
	{fPinkBfly, bBigBee, fLittleBees},       //14
	{bLittleBees, bLittleBees, bStripeBfly}, //15
}

//Each edge of a tile refers to the edge list
//and identifies which side faces inward to the tile
type Tile struct {
	edgeRef [3]int
	sideRef [3]int
}

var board = [16]Tile{
	{[3]int{11, 0, 12}, [3]int{1, 1, 1}},  //0
	{[3]int{10, 13, 15}, [3]int{1, 0, 1}}, //1
	{[3]int{13, 12, 14}, [3]int{1, 0, 1}}, //2
	{[3]int{14, 1, 16}, [3]int{0, 1, 1}},  //3
	{[3]int{9, 17, 21}, [3]int{1, 0, 1}},  //4
	{[3]int{17, 15, 18}, [3]int{1, 0, 1}}, //5
	{[3]int{18, 19, 22}, [3]int{0, 1, 1}}, //6
	{[3]int{19, 16, 20}, [3]int{0, 0, 1}}, //7
	{[3]int{20, 2, 23}, [3]int{0, 1, 1}},  //8
	{[3]int{8, 24, 7}, [3]int{1, 0, 1}},   //9
	{[3]int{24, 21, 25}, [3]int{1, 0, 1}}, //10
	{[3]int{25, 26, 6}, [3]int{0, 1, 1}},  //11
	{[3]int{26, 22, 27}, [3]int{0, 0, 1}}, //12
	{[3]int{27, 28, 5}, [3]int{0, 1, 1}},  //13
	{[3]int{28, 23, 29}, [3]int{0, 0, 1}}, //14
	{[3]int{29, 3, 4}, [3]int{0, 1, 1}},   //15
}

//Keep track of pieces as they are used
type PieceStatus struct {
	tileRef int
}

var pieceStatus [16]PieceStatus

//Keep track of tiles as they are used
type TileStatus struct {
	pieceRef int
	rotation int
}

var tileStatus [16]TileStatus

//Keep track of edges as they are used
type Edge struct {
	side [2]int
}

var edgeStatus [30]Edge

//=== Helper functions ===

//Clear all status, and set the border
func initBoard() {
	fmt.Println("Initialize Board")
	for i, _ := range edgeStatus {
		edgeStatus[i].side[0] = -1
		edgeStatus[i].side[1] = -1
	}
	for i, _ := range tileStatus {
		tileStatus[i].pieceRef = -1
		tileStatus[i].rotation = -1
	}
	for i, _ := range pieceStatus {
		pieceStatus[i].tileRef = -1
	}
	for i, _ := range border {
		edgeStatus[i].side[0] = border[i]
	}
}

//Print current tiles on the board
func printBoard() {
	fmt.Println("======================")
	fmt.Println("Board: ")
	for i, v := range tileStatus {
		if v.pieceRef == -1 {
			fmt.Printf("Tile: %d: No Tile\n", i)
		} else {
			fmt.Printf("Tile: %d: = piece: %d\n", i, v.pieceRef)
		}
	}
	fmt.Println("======================")
}

//Check 2 edges for a bisected match of same pattern
func edgeMatch(edge1, edge2 int) bool {
	if ((edge1 / 2) == (edge2 / 2)) && (edge1 != edge2) {
		return true
	}
	return false
}

//Check if all edges match (solution)
func allEdgesMatch() bool {
	for i, _ := range edgeStatus {

		if ((edgeStatus[i].side[0] == -1) || (edgeStatus[i].side[1] == -1)) || !edgeMatch(edgeStatus[i].side[0], edgeStatus[i].side[1]) {
			return false
		}
	}
	return true
}

//Place a rotated piece at tile location, update appropriate status
func placePiece(piece int, tile int, rot int) {
	// fmt.Println("Placing piece: ", piece, " on tile: ", tile, "rot: ", rot)
	tileStatus[tile].pieceRef = piece
	tileStatus[tile].rotation = rot
	pieceStatus[piece].tileRef = tile
	rotPiece := rotate(pieces[piece][0:], rot)
	for i, v := range rotPiece {
		edgeStatus[board[tile].edgeRef[i]].side[board[tile].sideRef[i]] = v
	}
}

//Remove a piece from board, clear appropriate status
func removePiece(tile int) {
	// fmt.Println(">>Removing ", tile)
	for i := 0; i < 3; i++ {
		edgeStatus[board[tile].edgeRef[i]].side[board[tile].sideRef[i]] = -1
	}
	pieceStatus[tileStatus[tile].pieceRef].tileRef = -1
	tileStatus[tile].pieceRef = -1
	tileStatus[tile].rotation = -1
}

//Rotate a piece that is already placed on the board
func rotatePiece(tile int) {
	if tileStatus[tile].pieceRef == -1 {
		return
	}
	tileStatus[tile].rotation = (tileStatus[tile].rotation + 1) % 3
	for i, v := range pieces[tileStatus[tile].pieceRef] {
		offset := (i + tileStatus[tile].rotation) % 3
		edgeStatus[board[tile].edgeRef[offset]].side[board[tile].sideRef[offset]] = v
	}
}

//Tests a piece for a match at a tile location in all 3 rotations
//Returns rot >=0 if match or -1 if no match
func checkMatch(piece int, tile int) int {
	for rot := 0; rot < 3; rot++ {
		match := true
		rotPiece := rotate(pieces[piece][0:], rot)
		for edge := 0; edge < 3; edge++ {
			adj := tileEdgeAdjacent(tile, edge)
			if (adj != -1) && (!edgeMatch(rotPiece[edge], adj)) {
				match = false
				break
			}
		}
		if match {
			return rot
		}
	}
	return -1
}

//Return the edge pattern on the outside ede of a tile
func tileEdgeAdjacent(tile int, edge int) int {
	otherSide := (board[tile].sideRef[edge] + 1) % 2
	return edgeStatus[board[tile].edgeRef[edge]].side[otherSide]
}

//Return factorial of a number
func factorial(n uint64) (result uint64) {
	if n > 0 {
		result = n * factorial(n-1)
		return result
	}
	return 1
}

//Return rotated slice of ints
func rotate(nums []int, k int) []int {
	if k < 0 || len(nums) == 0 {
		return nums
	}
	r := len(nums) - k%len(nums)
	nums = append(nums[r:], nums[:r]...)
	return nums
}

//=== Brute Force (one) ===

//Place all tiles on the board and then permute all rotations
//for a time check of one placement (not full brute force)
func bruteForceOne() {
	//Just one placement to do a time check on the rotations (n! more to go)
	fmt.Println("One tile placement")
	for i, _ := range pieces {
		placePiece(i, i, 0)
	}
	permuteRot()
}

//perform all 3^16 rotations of the tiles in one placement, checking each board for full match
func permuteRot() {
	count := int64(math.Pow(3, 16))
	fmt.Println("Checking ", count, " permutations")
	for i := int64(0); i < count; i++ {
		carry := true
		for k := 0; k < 16; k++ {
			if carry {
				tileStatus[k].rotation++
				if tileStatus[k].rotation == 3 {
					tileStatus[k].rotation = 0
					rotatePiece(k)
				} else {
					carry = false
				}
			} else {
				break
			}
		}
		if allEdgesMatch() {
			fmt.Println(">>SOLUTION<<")
			printBoard()
		}
	}
}

//=== Recursive Backtrack ===

//general backtrack pattern:
// func backtrack(item) {
// 	if item < 0 {
// 		printSolution()
// 		return
//     }
//     for choice := len(choices) -1; choice >=0; choice-- {
//         if goodChoice(choices[choice]) {
//             saveItemToSolution(items[item])
//             bactrack(item - 1)
//             deleteItemFromSolution(items[item])
//         }
//     }
// }

func backtrack(tile int) {
	if tile < 0 {
		fmt.Println(">>SOLUTION<<")
		printBoard()
		return
	}
	for piece := len(pieces) - 1; piece >= 0; piece-- {
		if pieceStatus[piece].tileRef != -1 {
			continue
		}
		rot := checkMatch(piece, tile)
		if rot >= 0 { //goodChoice
			placePiece(piece, tile, rot)
			backtrack(tile - 1)
			removePiece(tile)
		}
	}
}

//=== Non-Recursive Backtrack ===

func nrBacktrack(tile int) {
	for stackPtr := tile; ; {
		if stackPtr < 0 {
			fmt.Println(">>SOLUTION<<")
			printBoard()
		}
		for piece := len(pieces) - 1; ; piece-- {
			if piece < 0 {
				stackPtr++           //pop stack
				if stackPtr > tile { //all done
					return
				}
				piece = tileStatus[stackPtr].pieceRef
				removePiece(stackPtr)
				continue
			}
			if pieceStatus[piece].tileRef != -1 {
				continue
			}
			rot := checkMatch(piece, stackPtr)
			if rot >= 0 { //goodChoice
				placePiece(piece, stackPtr, rot) //push stack
				stackPtr--
				break
			}
		}
	}
}

func main() {
	fmt.Println("==================PUZZLE: Brute Force (one) ==================")
	start := time.Now()
	initBoard()
	bruteForceOne()
	end := time.Now()
	scaledDuration := int64(end.Sub(start) / time.Millisecond)
	fmt.Printf("Calculation finished in %d milliseconds\n", scaledDuration)
	f := factorial(16)
	years := (float64(f) * float64(scaledDuration)) / (1000.0 * 60 * 60 * 24 * 365)
	fmt.Printf("Doing that 16! times would take: %.0f years\n\n", years)
	fmt.Println("==================PUZZLE: Recursive Backtrack ==================")
	start = time.Now()
	initBoard()
	backtrack(15)
	end = time.Now()
	scaledDuration = int64(end.Sub(start) / time.Microsecond)
	fmt.Printf("Calculation finished in %d microseconds\n\n", scaledDuration)
	fmt.Println("==================PUZZLE: Non-Recursive Backtrack ==================")
	start = time.Now()
	initBoard()
	nrBacktrack(15)
	end = time.Now()
	scaledDuration = int64(end.Sub(start) / time.Microsecond)
	fmt.Printf("Calculation finished in %d microseconds\n", scaledDuration)
}
