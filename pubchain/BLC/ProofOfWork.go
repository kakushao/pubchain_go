package BLC

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
)

const targetBit = 20

type ProofOfWork struct {
	Block  *Block
	target *big.Int // 大数据存储
}

// 数据拼接，返回字节数组
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevBlockHash,
			pow.Block.HashTransactions(),
			IntToHex(pow.Block.Timestamp),
			IntToHex(int64(targetBit)),
			IntToHex(int64(nonce)),
			IntToHex(int64(pow.Block.Height)),
		},
		[]byte{},
	)

	return data
}

func (proofOfWork *ProofOfWork) Run() ([]byte, int64) {

	nonce := 0

	var hashInt big.Int // 存储我们新生成的hash
	var hash [32]byte

	for {
		dataBytes := proofOfWork.prepareData(nonce)

		hash = sha256.Sum256(dataBytes)
		fmt.Printf("\r%x", hash)

		hashInt.SetBytes(hash[:])

		if proofOfWork.target.Cmp(&hashInt) == 1 {
			break
		}

		nonce = nonce + 1
	}

	return hash[:], int64(nonce)
}

// 创建新的工作量证明对象
func NewProofOfWork(block *Block) *ProofOfWork {

	target := big.NewInt(1)

	target = target.Lsh(target, 256-targetBit)

	return &ProofOfWork{block, target}
}
