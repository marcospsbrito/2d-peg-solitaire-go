package main

import (
	"fmt"
)


const DOWN,UP,LEFT,RIGHT = 1,2,3,4

type tabuleiro struct {
	pieces       [9][9]int
	originCoord  [2]int
	destinyCoord [2]int
	movement int
}


func (t *tabuleiro) start() {
	t.init()
	fmt.Println("The game started!")
	t.doMovement()
	t.printPieces()
}

func (t *tabuleiro)  init(){
	key :=1
	for i:=0; i<9;i++{
		for j:=0; j<9;j++{
			if t.coordsValid(i,j) {
				t.pieces[i][j] = key
				key++
			}else{
				t.pieces[i][j] = -1
			}
		}
	}
	t.pieces[4][4] = 0
}

func (t *tabuleiro) selectOriginSlot(){
	fmt.Println("Enter de coordenates of the piece.")
	t.originCoord = t.askUserCoordenate()

}
func (t *tabuleiro) selectDestinySlot(){
	t.movement = 0
	t.askMovement()
}
func (t *tabuleiro) askMovement(){
	for t.movement < 1 || t.movement > 4 {
		fmt.Print("Movements \n ============== \n")
		if t.canOriginMove(DOWN){
			fmt.Println("1) Down" )
		}
		if t.canOriginMove(UP){
			fmt.Println("2) Up" )
		}
		if t.canOriginMove(LEFT) {
			fmt.Println("3) Left" )
		}
		if t.canOriginMove(RIGHT){
			fmt.Println("4) Right" )
		}
		fmt.Print("Type the number of movement \n")
		fmt.Scanf("%d", &t.movement)
		if !t.canMove(t.originCoord[0],t.originCoord[1],t.movement) {
			t.movement = 0
		}
	}
	fmt.Printf("Selected movement %d \n",t.movement)
}

func (t tabuleiro) askUserCoordenate() ([2]int) {
	var coordX, coordY =-1,-1

	for !isCoordValid(coordX) {
		fmt.Print("Enter the column: ")
		fmt.Scanf("%d", &coordX)
	}

	for !isCoordValid(coordY) {
		fmt.Print("Enter the line: ")
		fmt.Scanf("%d", &coordY)
	}
	fmt.Printf("Selected (%d, %d)\n",coordX,coordY)
	return [2]int{coordX,coordY}
}

func isCoordValid(coord int) bool {
	return coord >= 0 && coord < 9
}

func (t *tabuleiro) doMovement() {

	for t.hasMovements(){
		t.printPieces()
		for !t.isOriginValid() {
			t.selectOriginSlot()
		}
		for !t.isDestinyValid() {
			t.selectDestinySlot()
		}
		x, y := t.originCoord[0], t.originCoord[1]
		slotToMove, slotToFree := [2]int{x, y},[2]int{x, y}
		switch move :=t.movement; move{
		case 1:
			slotToMove[0] = slotToMove[0]+2
			slotToFree[0] = slotToFree[0]+1
		case 2:
			slotToMove[0] = slotToMove[0]-2
			slotToFree[0] = slotToFree[0]-1
		case 3:
			slotToMove[1] = slotToMove[1]-2
			slotToFree[1] = slotToFree[1]-1
		case 4:
			slotToMove[1] = slotToMove[1]+2
			slotToFree[1] = slotToFree[1]+1
		default:
			fmt.Println("Invalid movement! - ", move)
		}
		t.pieces[slotToMove[0]][slotToMove[1]] = t.pieces[t.originCoord[0]][t.originCoord[1]]
		t.pieces[slotToFree[0]][slotToFree[1]] = 0
		t.pieces[t.originCoord[0]][t.originCoord[1]] = 0
		t.movement=0
		fmt.Println("Moviment done!")
	}

}
func (tabuleiro *tabuleiro) printPieces() {
	fmt.Println("\\\\\tY\t0\t1\t2\t3\t4\t5\t6\t7\t8\t|")
	fmt.Println("X \t\t__\t__\t__\t__\t__\t__\t__\t__\t__\t")
	for x, line := range tabuleiro.pieces{
		fmt.Print(x,"\t|\t")
		for _, piece := range line{
			if piece == -1 {
				fmt.Print("++\t")
			}else if piece == 0 {
				fmt.Print("[]\t")
			}else{
				fmt.Print(piece,"\t")
			}
		}
		fmt.Print("|\n")
	}
}

func (t *tabuleiro) hasMovements() bool{
	for i:=0; i<9;i++{
		for j:=0; j<9;j++{
			if(t.hasMovement(i,j)){
				return true
			}
		}
	}
	return false
}
func (t *tabuleiro) hasMovement(x int, y int) bool{
	return !t.slotFree([2]int{x, y}) && t.canMove(x, y,1) || t.canMove(x, y, 2) || t.canMove(x, y,3) || t.canMove(x , y,4)
}
func (t *tabuleiro) isOriginValid() bool{
	valid := t.hasMovement(t.originCoord[0],t.originCoord[1]) && t.getPiece(t.originCoord[0],t.originCoord[1]) > 0
	if(!valid) {
		fmt.Printf("Slot (%d, %d) is invalid.\n",t.originCoord[0],t.originCoord[1])
	}
	return valid;
}
func (t *tabuleiro) canOriginMove(direction int) bool{
	return t.canMove(t.originCoord[0],t.originCoord[1],direction)
}
func (t *tabuleiro) canMove(x int, y int, direction int) bool{
	moves := [4][2]int{{x+2,y},{x-2,y},{x,y-2},{x,y+2}}
	piecesToEat := [4][2]int{{x+1,y},{x-1,y},{x,y-1},{x,y+1}}
	if t.coordsValid(x,y) && !t.slotFree([2]int{x,y}) && t.slotFree(moves[direction-1]) && !t.slotFree(piecesToEat[direction-1]){
			return true
	}
	return false
}
func (t *tabuleiro) slotFree(coords [2]int) bool {
	if t.coordsValid(coords[0],coords[1]){
		return  t.getPiece(coords[0],coords[1])==0
	}
	return false
}
func (t *tabuleiro) coordsValid(i,j int) bool{
	return (i>2 && i<6 && j>=0 && j<9) || (j>2 && j<6 && i>=0 && i<9)
	//return (i>2 && i<6 && j>=3 && j<6) || (j>2 && j<6 && i>=3 && i<6)
}

func (t *tabuleiro) isDestinyValid() bool{
	coords := t.originCoord
	switch move := t.movement; move{
	case DOWN:
		coords[0] = coords[0]+2
	case UP:
		coords[0] = coords[0]-2
	case LEFT:
		coords[1] = coords[1]-2
	case RIGHT:
		coords[1] = coords[1]+2
	default:
		fmt.Println("Invalid movement! - ", move)
	}
	if t.coordsValid(coords[0],coords[1]) && t.getPiece(coords[0], coords[1]) == 0 {
		fmt.Println("Slot (%d, %d) is Free.")
		return true
	}
	fmt.Printf("Slot (%d, %d) is not Free.\n",coords[0],coords[1])
	return false
}

func (t *tabuleiro) getPiece(coordX, coordY int) int{
	return t.pieces[coordX][coordY]
}

func main() {
	t:=tabuleiro{}
	t.start()
}
