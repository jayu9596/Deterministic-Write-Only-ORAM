package oram

import (
	"crypto/rand"
	"io"
	"math"
	"math/big"

	oram "github.com/jaydeep/Path-Oram/server"
)

// Life is much easier with json:  You are
// going to want to use this so you can easily
// turn complex structures into strings etc...
// Likewise useful for debugging etc
// UUIDs are generated right based on the crypto RNG
// so lets make life easier and use those too...
//
// You need to add with "go get github.com/google/uuid"
// Want to import errors

//BlockSize : Size of a oram.Block
var BlockSize = 8

//Stash
var stash = make([]oram.Block, 0)

var maxSize = 1000

var maxPath = 1000

var maxBucketSize = 10

var maxSizeOfStash = 0

//PositionMap : PositionMap structure used to store the oram.Block information
// type PositionMap struct {
// 	Map []int
// }

//Pair : Pair data structure
type Pair struct {
	first  int
	second int
}

var pMap []int

//getPMap : returns the position Map
func getPMap() []int {
	return pMap
}

//GetStash : returns the stash
func GetStash() []oram.Block {
	return stash
}

//StashClear : Clear the stash
func StashClear() {
	stash = make([]oram.Block, 0)
}

//StashClear : Clear the stash
func GetMaxStashSize() int {
	return maxSizeOfStash
}

//StashSet : set the value of stash at particular index
func StashSet(index int, data oram.Block) {
	stash[index] = data
	if len(stash) > int(math.Log2(float64(maxSize))+1)*maxBucketSize {
		println("Overflowed!")
		println(len(stash))
	}
}

//StashDelete : delete entry in stash
func StashDelete(index int) {
	//TODO:delete element at this index
	stash = append(stash[:index], stash[index+1:]...)
}

//StashPush : Add entry in stash
func StashPush(data oram.Block) {
	stash = append(stash, data)
	if len(stash) > int(math.Log2(float64(maxSize))+1)*maxBucketSize {
		println("Overflowed!")
		println(len(stash))
	}
	if maxSizeOfStash < len(stash) {
		maxSizeOfStash = len(stash)
	}
}

//StashPop : delete entry in stash
func StashPop(index int) {
	stash = stash[1:]
}

//setMaxBucketSize : set setMaxBucketSize
func setMaxBucketSize(size int) {
	maxBucketSize = size
}

//RandomBytes : Helper function: Returns a byte slice of the specificed
// size filled with random data
func RandomBytes(bytes int) (data []byte) {
	data = make([]byte, bytes)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		panic(err)
	}
	return
}

func isBlockPresent(index int) int {
	i := 0
	for i < len(stash) {
		if stash[i].Pos == index {
			return i
		}
		i++
	}
	return -1
}

//ServerUploadData : First time this method is to be called
// This is used for uploading all data as single Binary Tree in Server
func ServerUploadData(plaintext []byte) (err error) {
	a := int(math.Ceil(float64(len(plaintext)) / (float64(BlockSize * maxBucketSize))))
	maxPath = (a + 1) / 2
	queue := make([]Pair, 0)
	pMap = make([]int, (a*maxBucketSize)+1)
	queue = append(queue, Pair{0, (a + 1) / 2})
	maxSize = a
	iter := 0
	for i := 0; i < a; i++ {
		var pair = queue[0]
		queue = queue[1:]
		blockArray := make([]oram.Block, 0)
		if pair.first != pair.second {
			queue = append(queue, Pair{pair.first, (pair.first + pair.second) / 2})
			queue = append(queue, Pair{(pair.first + pair.second) / 2, pair.second})
		}
		for j := 0; j < maxBucketSize; j++ {
			pMap[iter] = pair.first
			block := oram.Block{iter, plaintext[iter*BlockSize : (iter+1)*BlockSize]}
			blockArray = append(blockArray, block)
			iter++
		}
		bucket := oram.Bucket{maxBucketSize, blockArray}
		oram.Append(int(i), bucket)
	}
	return nil
}

//ReadPath : Reads blocks on this path
// Implementation details:  currNode will become the last node and after that we divide currNode by 2 till we reach root node
// and read all nodes on this path
func ReadPath(index int) (oram.Block, error) {
	path := pMap[index]
	id := isBlockPresent(index)
	retData := oram.Block{-1, RandomBytes(BlockSize)}
	if id != -1 {
		retData = stash[id]
	}

	//currNode will be the leaf node and then it will become it's parent till it reaches root
	currNode := int(path) + (maxSize / 2)

	flag := true
	for flag {
		if currNode == 0 {
			flag = false
		}
		j := 0
		for j < maxBucketSize {
			block, _ := oram.Read(int(currNode), j)
			//Read and Store the block in the stash
			//If Dummy Block then don't store in stash
			if block.Pos == -1 {
				currNode = (currNode - 1) / 2
				continue
			}
			idx := isBlockPresent(block.Pos)
			// _, ok := stash[block.Pos]
			if idx == -1 {
				StashPush(block)
				// StashSet(block.Pos, block)

				if block.Pos == index {
					retData = block
				}
			}
			j++
		}
		currNode = (currNode - 1) / 2
	}
	return retData, nil
}

//isIndexInPath : check if index will come this path or not
func isIndexInPath(index int, path2 int) bool {
	currNode := path2 + (maxSize / 2)
	ret := false
	flag := true
	for flag {
		if currNode == 0 {
			flag = false
		}
		if currNode == index {
			ret = true
			flag = false
		}
		currNode = (currNode - 1) / 2
	}
	return ret
}

//Writepath : write provded datablocks on the path
func Writepath(path int) (ciphertext []oram.Block, err error) {
	//currNode will be the leaf node and then it will become it's parent till it reaches root
	currNode := int(path) + (maxSize / 2)
	flag := true
	for flag {
		if currNode == 0 {
			flag = false
		}
		j := 0
		for j < maxBucketSize {
			//find and write a block which can fit in this node
			found := false
			iter := 0
			for ; iter < len(stash); iter++ {
				if isIndexInPath(currNode, pMap[stash[iter].Pos]) {
					oram.Write(currNode, j, stash[iter])
					StashDelete(iter)
					found = true
					break
				}
			}

			//If no such blocks found Write dummy block
			if !found {
				block := oram.Block{-1, RandomBytes(BlockSize)}
				oram.Write(currNode, j, block)
			}
			j++
		}
		currNode = (currNode - 1) / 2
	}
	return nil, nil
}

//GetMyDB : returns database
func GetMyDB() []oram.Bucket {
	return oram.GetDB()
}

//Access :should read/write data at index
func Access(index int, data []byte) ([]byte, error) {
	//Update data in stash
	if data != nil {
		idx := isBlockPresent(index)
		if idx == -1 {
			StashPush(oram.Block{index, data})
		} else {
			stash[idx] = oram.Block{index, data}
		}
	}

	//ReadPath
	retData, _ := ReadPath(index)
	path := pMap[index]

	//Assign new path for this index
	newPathBigInt, _ := rand.Int(rand.Reader, big.NewInt(int64(maxPath)))
	newPath := newPathBigInt.Int64() % int64(maxPath)
	// for newPath == int64(path) {
	// 	newPathBigInt, _ = rand.Int(rand.Reader, big.NewInt(int64(maxPath)))
	// 	newPath = newPathBigInt.Int64() % int64(maxPath)
	// }
	pMap[index] = int(newPath)

	//Read
	if data == nil {
		Writepath(path)
		return retData.Data, nil
	}

	//Write
	Writepath(path)
	return nil, nil
}
