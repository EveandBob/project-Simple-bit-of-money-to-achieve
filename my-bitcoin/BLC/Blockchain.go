package BLC

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"time"
)

//数据库的名字
const dbName = "blockchain.db"

//添加表的名字
const blockTableName = "blocks"

//区块链
type Blockchain struct {
	Tip []byte //最新区块的Hash
	DB  *bolt.DB
}

//增加区块到区块链里面
//func (blc *Blockchain) AddBlockToBloickchain(data string, height int64, preHash []byte) {
//创建新区块
//	newBlock := NewBlock(data, height, preHash)
//往链里添加区块
//	blc.Blocks = append(blc.Blocks, newBlock)
//}

//使用迭代器
func (blockchain *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{blockchain.Tip, blockchain.DB}
}

//遍历输出所有区块的信息
func (blc *Blockchain) PrintChain() {
	blockchainIterator := blc.Iterator()
	for {
		block := blockchainIterator.Next()
		//打印幸喜
		fmt.Printf("Height: %d\n", block.Height)
		fmt.Printf("PrevBlockHash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Timestamp: %s\n", time.Unix(block.Timestamp, 0).Format("2006-01.02 03:04:05 PM"))
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Println()
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)
		if big.NewInt(0).Cmp(&hashInt) == 0 {
			break
		}
	}
}

//1.创建带有创世区块的区块链
func CreateBlockChainWithGenesisBlock() *Blockchain {
	//创建或者打开数据库
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	var blockHash []byte
	//更新数据库
	err = db.Update(func(tx *bolt.Tx) error {

		b, err := tx.CreateBucket([]byte(blockTableName))
		if err != nil {
			log.Panic(err)
		}
		if b != nil {
			//创建创世区块
			genesisBlock := CreateGenesisBlock("Genesis Data")
			err := b.Put(genesisBlock.Hash, genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}

			//存储最新的区块的hash
			err = b.Put([]byte("l"), genesisBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
			blockHash = genesisBlock.Hash
		}
		return nil
	})

	//创建创世区块
	return &Blockchain{blockHash, db}
}

// add new block to blockchain
func (blockchain *Blockchain) AddBlockToBlockchain(data string) {

	// update database
	err := blockchain.DB.Update(func(tx *bolt.Tx) error {
		// 1.get table
		bucket := tx.Bucket([]byte(blockTableName))

		// 2.create new block
		if bucket != nil {
			// 3.get latest block
			blockBytes := bucket.Get(blockchain.Tip)
			// deserialize
			block := DeserializeBlock(blockBytes)

			// 4.store new block
			newBlock := NewBlock(data, block.Height+1, block.Hash)
			err := bucket.Put(newBlock.Hash, newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}

			// 4.update "l"
			err = bucket.Put([]byte("l"), newBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
			// 5.update Tip
			blockchain.Tip = newBlock.Hash
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}
