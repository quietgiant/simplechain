package main

import (
	"os"
	"fmt"
	"bufio"
	"time"
	"math/rand"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"strconv"
	"github.com/briandowns/spinner"
)

// VARIABLES

// a block which holds one piece of data
type block struct {
	index int
	timestamp int
	data string
	hash string
	previousHash string
}

var blockchain []block // blockchain
var target int // target for miners
var difficulty string // difficulty for proof of work

// MAIN

func main() {
	blockchain = append(blockchain, getGenesisBlock())
	initProofOfWork()
	for {
		printMenu()
	}
}

// INTERFACE

func printMenu() {
	menu := "Simplechain functions:\n"
	menu += "Type 'print' to print current blockchain\n"
	menu += "Type 'mine' to mine new block\n"
	menu += "Type 'quit' to quit program\n> "
	fmt.Print(menu)
	var input string
    _, err := fmt.Scanln(&input)
    fmt.Println()
    if err != nil {
         fmt.Println("Invalid input. Please try again.\n")
         return
    }
    switch strings.ToLower(input) {
	    case "print":
			s := "Blockchain as of: " + time.Now().String() + "\n\n"
			fmt.Println(s + blockchainToString() + "\n")
	    case "mine":
	    	if (!mineBlock()) { // attempt to mine a new block
	    		fmt.Println("Error occured mining new block.\n")
	    	}
	    case "quit":
	    	fmt.Println("Quitting...")
	        os.Exit(0) // successful exit
	    default:
	        fmt.Println("Invalid command. Please try again.\n")
    }
}

// BLOCKCHAIN FUNCTIONS

func generateNewBlock(data string) block {
	prevBlock := getCurrentBlock()
	newIndex := prevBlock.index + 1
	newTime := int(time.Now().Unix())
	newData := data
	newHash := calculateHash(newIndex, newTime, newData, prevBlock.hash)
	return block{newIndex, newTime, newData, newHash, prevBlock.hash}	
}

func addBlock(newBlock block) bool {
	if (validateNewBlock(getCurrentBlock(), newBlock)) { // ensure block is valid
		if (tryProofOfWork(newBlock)) { // force proof of work in order to mine new block
			blockchain = append(blockchain, newBlock)
			fmt.Printf("Block successfully added. (height: %d)\n\n", newBlock.index)
		}
	} else {
		fmt.Println("Error occured adding new block.")
		return false
	}
	return true
}

func mineBlock() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter data for block\n> ")
	data, _, _ := reader.ReadLine()
	if (!addBlock(generateNewBlock(string(data[:])))) { // attempt to create a new block
		return false // mine failed
	}
	return true // mine successful
}

func tryProofOfWork(newBlock block) bool {
	var test string
	t := strconv.Itoa(target)
	nonce := 0 // generate nonce to test against target hashes
	icon := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	icon.Prefix = "Computing proof of work "
	go icon.Start()
	for {
		test = t + strconv.Itoa(nonce)
		h := sha256.New()
		h.Write([]byte(test))
		digest := hex.EncodeToString(h.Sum(nil))
		if (strings.HasPrefix(digest, difficulty)) {
			icon.Stop()
			fmt.Println("Problem solved!\nNonce: " + strconv.Itoa(nonce))
			setDifficulty()
			return true
		}
		nonce += 1
	}
	icon.Stop()
	return false
}

func validateNewBlock(prevBlock block, newBlock block) bool {
	if (prevBlock.index + 1 != newBlock.index) {
		fmt.Println("ERROR: Block index invalid...")
		fmt.Printf("%d : %d\n", prevBlock.index + 1, newBlock.index)
		fmt.Println("New block rejected.")
		return false
	} else if (prevBlock.hash != newBlock.previousHash) {
		fmt.Println("ERROR: Previous block hash invalid...")
		fmt.Printf("%s : %s\n", prevBlock.hash, newBlock.previousHash)
		fmt.Println("New block rejected.")
		return false
	} else if (calculateBlockHash(newBlock) != newBlock.hash) {
		fmt.Println("ERROR: Current block hash invalid...")
		fmt.Printf("%s : %s\n", calculateBlockHash(newBlock), newBlock.hash)
		fmt.Println("New block rejected.")
		return false
	}
	return true
}

func validateChain(chainToValidate []block) bool {
	if (chainToValidate[0] != getGenesisBlock()) {
		return false
	}

	var tempChain = make([]block, 0)
	tempChain = append(tempChain, getGenesisBlock())
	
	for i := 1; i < len(chainToValidate); i++ {
		if (validateNewBlock(tempChain[i - 1], chainToValidate[i])) {
			tempChain = append(tempChain, chainToValidate[i])
		} else {
			return false
		}
	}
	return true // chain is valid
}

func replaceChain(incomingChain []block) {
	if (validateChain(incomingChain) && len(incomingChain) > len(blockchain)) {
		fmt.Println("Incoming blockchain is valid. Replacing current blockchain...")
		blockchain = incomingChain
		fmt.Println("Blockchain successfully replaced. Fork resolved.")
	} else {
		fmt.Print("Incoming blockchain invalid. Rejected.")
	}
}

func calculateBlockHash(b block) string {
	return calculateHash(b.index, b.timestamp, b.data, b.previousHash)
}

func calculateHash(index int, timestamp int, data string, previousHash string) string {
	payload := strconv.Itoa(index) + strconv.Itoa(timestamp) + data + previousHash
	h := sha256.New()
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}

func initProofOfWork() {
	rand.Seed(time.Now().UTC().UnixNano())
	target = rand.Int()
	setDifficulty()
}

func getCurrentBlock() block {
	return blockchain[len(blockchain) - 1]
}

func getGenesisBlock() block {
	genData := "genesis block"
	h := sha256.New()
	h.Write([]byte(genData))
	genHash := hex.EncodeToString(h.Sum(nil))
	return block{0, 0, genData, genHash, "0"}
}

func setDifficulty() {
	difficulty = "0"
	max := 5
	n := rand.Intn(max) + 1 // dynamically adjust difficulty
	for i := 0; i < n; i++ {
		difficulty += "0"
	}
}

func blockchainToString() string {
	var s string
	for _, block := range blockchain {
		s += blockToString(block)
		s += "\n"
	}
	return s
}

func blockToString(b block) string {
	s := "Index: " + strconv.Itoa(b.index)
	s += "\nTimestamp: " + strconv.Itoa(b.timestamp)
	s += "\nData: " + b.data
	s += "\nHash: " + b.hash
	s += "\nPrevious hash: " + b.previousHash + "\n"
	return s
}
