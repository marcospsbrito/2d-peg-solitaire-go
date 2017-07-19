package main

import 	"fmt"

const (
	DOWN  = 1
	UP    = 2
	LEFT  = 3
	RIGHT = 4
	EMPTY = "E"
	DENIED = "D"
	NOT_AVALIABLE = "NA"
)

type slot struct {
	x,y int
	slotType string
	piece *piece
}

func (s slot) isEmpty() bool{
	return s.slotType == EMPTY
}

func (s slot) hasPiece() bool{
	return s.slotType == NOT_AVALIABLE
}
func (s *slot) setPiece(piece *piece){
	s.piece = piece
	s.piece.setSlot(s)
	s.slotType=NOT_AVALIABLE
}
func (piece *piece) setSlot(slot *slot) {
	piece.s = slot
}

func (s *slot) free() {
	if s.slotType == NOT_AVALIABLE{
		*s = slot{s.x, s.y, EMPTY, nil}
	}
}

type piece struct {
	key int
	s *slot
}

type tabuleiro struct {
	slots [9][9]*slot
	piecesMap map[int]*piece
	pieceSelected int
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
	t.piecesMap = make(map[int]*piece)
	for i:=0; i<9;i++{
		for j:=0; j<9;j++{
				slot := new(slot)
			if t.coordsValid(i,j) {
				piece := piece{key, slot}
				slot.x,slot.y,slot.slotType,slot.piece = i,j,NOT_AVALIABLE,&piece
				t.piecesMap[key] = &piece
				key++
			}else{
				slot.x,slot.y,slot.slotType,slot.piece = i,j,DENIED, nil
			}
			t.slots[i][j] = slot
		}
	}
	t.slots[4][4].free()
}

func (t *tabuleiro) selectOriginSlot(){
	t.askPieceNumber()
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
		fmt.Print("Type the movement number: ")
		fmt.Scanf("%d", &t.movement)
		if !t.canMove(t.pieceSelected, t.movement) {
			t.movement = 0
		}
	}
	fmt.Printf("Selected movement %d \n", t.movement)
}

func (t *tabuleiro) askPieceNumber() {
	t.pieceSelected = -1

	for valid := false; !valid; valid = t.isPieceValid(t.pieceSelected) {
		fmt.Print("Enter the Piece Number: ")
		fmt.Scanf("%d", &t.pieceSelected)
	}

	fmt.Printf("Piece selected (%d)\n",t.pieceSelected)
}
func (t tabuleiro) isPieceValid(pieceNumber int) bool {
	if val, ok := t.piecesMap[pieceNumber]; ok{
		if &val != nil && !val.s.isEmpty(){
			return true
		}
	}
	fmt.Println("Invalid pieceNumber ", pieceNumber)
	return false
}

func (t *tabuleiro) doMovement() {

	for t.hasMovements(){
		t.printPieces()
		originValid,destinyValid := false,false
		for !originValid {
			t.selectOriginSlot()
			originValid = t.isOriginValid()
		}
		for !destinyValid{
			t.selectDestinySlot()
			destinyValid = t.isDestinyValid()
		}
		t.moveSelectedPiece()
		t.movement=0
		fmt.Println("Moviment done!")
	}

}
func (t *tabuleiro) moveSelectedPiece() {
	s := t.getPiece(t.pieceSelected)
	switch move := t.movement; move{
	case UP:
		t.move(s, t.slots[s.x-2][s.y])
	case DOWN:
		t.move(s, t.slots[s.x+2][s.y])
	case LEFT:
		t.move(s, t.slots[s.x][s.y-2])
	case RIGHT:
		t.move(s, t.slots[s.x][s.y+2])
	default:
		fmt.Println("Invalid movement! - ", move)
	}
}
func (t *tabuleiro) move(origin *slot, destiny *slot) {
	slotToFree := t.slots[origin.x+((destiny.x-origin.x)/2)][origin.y+((destiny.y-origin.y)/2)]
	destiny.setPiece(origin.piece)
	slotToFree.free()
	origin.free()
}
func (t *tabuleiro) printPieces() {
	fmt.Println("\\\\\tY\t0\t1\t2\t3\t4\t5\t6\t7\t8\t|")
	fmt.Println("X \t\t__\t__\t__\t__\t__\t__\t__\t__\t__\t")
	for x, line := range t.slots{
		fmt.Print(x,"\t|\t")
		for _, slot := range line{
			switch slot.slotType {
			case DENIED:
				fmt.Print("++\t")
			case NOT_AVALIABLE:
				fmt.Print(slot.piece.key,"\t")
			case EMPTY:
				fmt.Print("[]\t")
			}
		}
		fmt.Print("|\n")
	}
}

func (t *tabuleiro) hasMovements() bool{
	for i:=0; i<9;i++{
		for j:=0; j<9;j++{
			if t.hasMovement(i,j) {
				return true
			}
		}
	}
	return false
}
func (t *tabuleiro) hasMovement(x int, y int) bool{
	if t.slots[x][y].piece==nil{
		return false
	}
	pk :=t.slots[x][y].piece.key
	return !t.slotFree([2]int{x, y}) && t.canMove(pk,1) ||
		t.canMove(pk, 2) || t.canMove(pk,3) || t.canMove(pk,4)
}
func (t *tabuleiro) isOriginValid() bool{
	slot := t.getPiece(t.pieceSelected)
	valid := t.hasMovement(slot.x, slot.y) && slot.hasPiece()
	if !valid  {
		fmt.Printf("Piece (%d) can't be moved.\n", t.pieceSelected)
	}
	return valid
}
func (t *tabuleiro) canOriginMove(direction int) bool{
	return t.canMove(t.pieceSelected, direction)
}
func (t *tabuleiro) canMove(pieceNumber int, direction int) bool{
	if direction <1 || direction > 4 || t.getPiece(pieceNumber) == nil {
		fmt.Println("no piece found - ", pieceNumber)
		return false
	}
	s := t.getPiece(pieceNumber)
	moves := [4][2]int{{s.x+2,s.y},{s.x-2,s.y},{s.x,s.y-2},{s.x,s.y+2}}
	piecesToEat := [4][2]int{{s.x+1,s.y},{s.x-1,s.y},{s.x,s.y-1},{s.x,s.y+1}}
	if t.coordsValid(s.x,s.y) && s.hasPiece() && t.slotFree(moves[direction-1]) && !t.slotFree(piecesToEat[direction-1]){
			return true
	}
	return false
}
func (t *tabuleiro) slotFree(coords [2]int) bool {
	if t.coordsValid(coords[0],coords[1]){
		return  t.slots[coords[0]][coords[1]].isEmpty()
	}
	return false
}
func (t *tabuleiro) coordsValid(i,j int) bool{
	return (i>2 && i<6 && j>=0 && j<9) || (j>2 && j<6 && i>=0 && i<9)
}

func (t *tabuleiro) isDestinyValid() bool{
	coords := [2]int{t.getPiece(t.pieceSelected).x,t.getPiece(t.pieceSelected).y}
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
	if t.coordsValid(coords[0],coords[1]) && t.slots[coords[0]][coords[1]].isEmpty() {
		return true
	}
	return false
}

func (t *tabuleiro) getPiece(pieceNumber int) *slot{
	return t.piecesMap[pieceNumber].s
}

func main() {
	t:=tabuleiro{}
	t.start()
}