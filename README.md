# 项目名称
用go 语言实现一个简单的区块链项目

# 工作量说明
这个项目是我花费时间最长的项目，从学习go语言和bolt数据库到基本实现耗时4天

# 项目功能介绍
1.该代码实现了区块链的基本构型

2.该项目能够进行挖矿

3.该项目可以将区块的内容存储在boltdb数据库中

4.该项目实现了对区块内容的格式化输出

# 项目的算法的代码
1.创建创世区块(并插入数据库中)
```go
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
```
2.添加区块(并插入数据库中)
```go
// 添加区块
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
```
3.进行工作量证明
```go
func (proofOfWork *ProofOfWork) IsVailue() bool {
	//1.proofOfWork.Block,Hash
	//2.proofOfWork.Target
	var hashInt big.Int
	hashInt.SetBytes(proofOfWork.Block.Hash)
	if proofOfWork.targrt.Cmp(&hashInt) == 1 {
		return true
	}
	return false
}

// 进行工作量证明算法
func (proofOfWork *ProofOfWork) Run() ([]byte, int64) {
	//1.将Block的属性拼接成字节数组
	//2.生成hash
	//3.判断hash的有效性，如果满足条件，跳出循环
	nonce := 0
	var hashInt big.Int // 存储新生成的Hash
	var hash [32]byte

	for {
		//准备数据
		dataBytes := proofOfWork.prepareData(nonce)
		//对数据做Hash
		hash = sha256.Sum256(dataBytes)
		//将hash存储到hashInt
		hashInt.SetBytes(hash[:])

		//判断hash是否有效
		if proofOfWork.targrt.Cmp(&hashInt) == 1 {
			fmt.Printf("\r%x\n", hash)
			break

		}

		nonce = nonce + 1
	}
	return hash[:], int64(nonce)
}

//创建新的工作量证明对象
func NewProofOfWork(block *Block) *ProofOfWork {
	//1.创建一个初始值为1的target
	target := big.NewInt(1)
	target = target.Lsh(target, 256-targetBit)
	return &ProofOfWork{block, target}
}
```
4.对区块内容格式化输出
```go
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
```
# 最终效果
我以下列代码对区块链进行操作
```go
//创世区块
	Blockchain := BLC.CreateBlockChainWithGenesisBlock()
	Blockchain.AddBlockToBlockchain("send 2 bitcoin to satoshi")
	Blockchain.AddBlockToBlockchain("send 3 bitcoin to satoshi")
	Blockchain.PrintChain()
```

# 结果如下
![Screenshot 2022-07-31 145722](https://user-images.githubusercontent.com/104854836/182014070-b9fc83f3-9701-4623-977a-b1616a064aa7.jpg)


# 心得
时间有限不能实现全部的比特币算法的内容，但是在假期的剩余时间我会继续学习并争取完成，随之会上传到github中
