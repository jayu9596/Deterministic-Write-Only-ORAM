package oram

// Life is much easier with json:  You are
// going to want to use this so you can easily
// turn complex structures into strings etc...

// Likewise useful for debugging etc

// UUIDs are generated right based on the crypto RNG
// so lets make life easier and use those too...
//
// You need to add with "go get github.com/google/uuid"

//Block : Block structure used to store the block in binary tree
type Block struct {
	Pos  int
	Data []byte
}

//Bucket : Bucket structure used to store the number of blocks
type Bucket struct {
	Size   int
	Blocks []Block
}

// var mydb = make([]Block, 1000)
var mydb = []Bucket{}

//getDB : returns database
func GetDB() []Bucket {
	return mydb
}

//Read : reads data at index
func Read(index int, subIndex int) (data Block, err error) {
	// fmt.Println(index)
	return mydb[index].Blocks[subIndex], nil
}

//Write :should overwrite data at index
func Write(index int, subIndex int, block Block) (err error) {
	// fmt.Println(index)
	mydb[index].Blocks[subIndex] = block
	return nil
}

//Append :should Append data
func Append(index int, bucket Bucket) (err error) {
	// block := Block{index, ciphertext, false}
	mydb = append(mydb, bucket)
	return nil
}
