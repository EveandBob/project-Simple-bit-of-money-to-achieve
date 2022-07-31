package BLC

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"strconv"
	"time"
)

type Block struct {
	//1.区块高度
	Height int64
	//2.上一个区块的Hash
	PrevBlockHash []byte
	//3.交易数据
	Data []byte
	//4.时间戳
	Timestamp int64
	//5.Hash
	Hash []byte
	//6.Nonce
	Nonce int64
}

func (block *Block) SetHash() {
	//1.Height转化为字节数组
	heightBytes := IntToHex(block.Height)

	//2.timestmp转化为字节数组
	timeString := strconv.FormatInt(block.Timestamp, 2)
	timeBytes := []byte(timeString)

	//3.拼接所有的属性
	blockBytes := bytes.Join([][]byte{heightBytes, block.PrevBlockHash, block.Data, timeBytes, block.Hash}, []byte{})

	//4.计算hash值
	hash := sha256.Sum256(blockBytes)

	block.Hash = hash[:]

}

//1.创建新区块
func NewBlock(data string, height int64, prevBlockHash []byte) *Block {
	//创建区块
	block := &Block{Height: height,
		PrevBlockHash: prevBlockHash,
		Data:          []byte(data),
		Timestamp:     time.Now().Unix(),
		Hash:          nil,
		Nonce:         0}
	//调用工作量证明的方法并且返回有效的Hash和Nonce
	pow := NewProofOfWork(block)

	//执行工作量证明
	hash, nonce := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

//2.生成创世区块
func CreateGenesisBlock(data string) *Block {
	return NewBlock("Genenis Block", 1, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}

//将区块序列化为字节数组
func (block *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(block); err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

//反序列化
func DeserializeBlock(blockBytes []byte) *Block {
	var deBlock Block
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	if err := decoder.Decode(&deBlock); err != nil {
		log.Panic(err)
	}
	return &deBlock
}
