package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"sort"
	"strconv"
	"time"
)

const (
	voteNodeNum      = 200
	superNodeNum     = 20
	mineSuperNodeNum = 5
)

type block struct {
	//Hash of the previous block
	preHash string
	//Hash of this block
	hash string
	//Timestamp
	timestamp string
	//Block Content
	data string
	//Block Height
	height int
	//Address of miner of block
	address string
}

//blockchain
var blockchain []block

//Normal Nodes
type node struct {
	//Number of tokens
	votes int
	//Node Address
	address string
}

//Validator Node
type superNode struct {
	node
}

//voting node pool
var voteNodesPool []node

//campaign node pool
var starNodesPool []superNode

//Store the super node pool that can be mined
var superStarNodesPool []superNode

//generate new block
func generateNewBlock(oldBlock block, data string, address string) block {
	newBlock := block{}
	newBlock.preHash = oldBlock.hash
	newBlock.data = data
	newBlock.timestamp = time.Now().Format("2006-01-02 15:04:05")
	newBlock.height = oldBlock.height + 1
	newBlock.address = address
	newBlock.getHash()
	return newBlock
}

//hash itself
func (b *block) getHash() {
	sumString := b.preHash + b.timestamp + b.data + b.address + strconv.Itoa(b.height)
	hash := sha256.Sum256([]byte(sumString))
	b.hash = hex.EncodeToString(hash[:])
}

//vote
func voting() {
	for _, v := range voteNodesPool {
		rInt, err := rand.Int(rand.Reader, big.NewInt(superNodeNum+1))
		if err != nil {
			log.Panic(err)
		}
		starNodesPool[int(rInt.Int64())].votes += v.votes
	}
}

//Sort mining nodes
func sortMineNodes() {
	sort.Slice(starNodesPool, func(i, j int) bool {
		return starNodesPool[i].votes > starNodesPool[j].votes
	})
	superStarNodesPool = starNodesPool[:mineSuperNodeNum]
}

//initialization
func init() {
	//Initialize voting nodes
	for i := 0; i <= voteNodeNum; i++ {
		rInt, err := rand.Int(rand.Reader, big.NewInt(10000))
		if err != nil {
			log.Panic(err)
		}
		voteNodesPool = append(voteNodesPool, node{int(rInt.Int64()), "voting node" + strconv.Itoa(i)})
	}
	//Initialize campaign node
	for i := 0; i <= superNodeNum; i++ {
		starNodesPool = append(starNodesPool, superNode{node{0, "super node" + strconv.Itoa(i)}})
	}
}

func main() {
	fmt.Println("initialization", voteNodeNum, "voting nodes...")
	fmt.Println(voteNodesPool)
	fmt.Println("currently existing", superNodeNum, "campaign nodes")
	fmt.Println(starNodesPool)
	fmt.Println("Voting nodes start voting...")
	voting()
	fmt.Println("End the voting and check the votes received by the campaign nodes...")
	fmt.Println(starNodesPool)
	fmt.Println("Sort the campaign nodes by the number of votes they get", mineSuperNodeNum, "Name, elected as super node")
	sortMineNodes()
	fmt.Println(superStarNodesPool)
	fmt.Println("start mining...")
	genesisBlock := block{"0000000000000000000000000000000000000000000000000000000000000000", "", time.Now().Format("2022-07-20 13:54:05"), "I am the genesis block", 1, "000000000"}
	genesisBlock.getHash()
	blockchain = append(blockchain, genesisBlock)
	fmt.Println(blockchain[0])
	i, j := 0, 0
	for {
		time.Sleep(time.Second)
		newBlock := generateNewBlock(blockchain[i], "I am block content", superStarNodesPool[j].address)
		blockchain = append(blockchain, newBlock)
		fmt.Println(blockchain[i+1])
		i++
		j++
		j = j % len(superStarNodesPool) //Super nodes obtain the right to produce blocks in a round-robin fashion
	}
}
