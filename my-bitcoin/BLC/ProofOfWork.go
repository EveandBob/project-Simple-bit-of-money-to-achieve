package BLC

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

//工作量证明中目标hash要有16个零
const targetBit = 16

type ProofOfWork struct {
	Block  *Block   //当前要验证的区块
	targrt *big.Int //大数据存储
}

//数据凭借，返回字节数组
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevBlockHash,
			pow.Block.Data,
			IntToHex(int64(targetBit)),
			IntToHex(int64(nonce)),
			IntToHex(int64(pow.Block.Height)),
		},
		[]byte{})
	return data
}

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
