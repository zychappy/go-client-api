package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/go-ethereum/crypto/sha3"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/zychappy/go-client-api/common/base58"
	"github.com/zychappy/go-client-api/common/crypto"
	"github.com/zychappy/go-client-api/core"
	"github.com/zychappy/go-client-api/service"
	"github.com/zychappy/go-client-api/util"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func GetAddressByte(puk []byte) []byte {
	fmt.Printf("pubkey:%x\n", puk)
	sha3 := sha3.NewKeccak256()
	sha3.Write(puk)
	pub := sha3.Sum(nil)
	fmt.Printf("sha3:%x\n", pub)
	pub20 := pub[len(pub)-20:]
	pubmain := append([]byte{0x41}, pub20...)
	return pubmain
}

func int2Byte(x interface{}) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, x)
	return buf.Bytes()
}

func getInfo(reqBody string, url string) (resBody string, err error) {

	resp, err := http.Post(url,
		"application/x-www-form-urlencoded",
		strings.NewReader(reqBody))
	if err != nil {
		fmt.Println(err)
	}
	res, err := ioutil.ReadAll(resp.Body)
	return string(res), err
}
func main() {
	grpcAddress := flag.String("grpcAddress", "127.0.0.1:50051",
		"gRPC address: <IP:port> example: -grpcAddress localhost:50051")
	client := service.NewGrpcClient(*grpcAddress)
	client.Start()
	defer client.Conn.Close()
	ownerPriv := "D95611A9AF2A2A45359106222ED1AFED48853D9A44DEFF8DC7913F5CBA727366"
	// 拿到私钥对象
	//priv, _ := secp256k1.HexToECDSA(ownerPriv)
	priv, _ := crypto.GetPrivateKeyByHexString(ownerPriv)
	// owner公钥
	//ownerPub := priv.Public().Bytes()
	ownerPub := append(priv.X.Bytes(), priv.Y.Bytes()...)
	// 根据公钥获取原始地址,主网41开头
	ownerAddress := GetAddressByte(ownerPub)
	fmt.Println("owner address:", ownerAddress)
	//接收方
	toAddr := "TGehVcNhud84JDCGrNHKVz9jEAVKUpbuiv"
	// 进行base58解码
	to := base58.DecodeCheck(toAddr)
	fmt.Println("to addr", hex.EncodeToString(to))

	tfc := core.TransferContract{
		OwnerAddress: ownerAddress,
		ToAddress:    to,
		Amount:       100000000,
	}
	tfcBytes, _ := proto.Marshal(&tfc)
	aa := any.Any{
		TypeUrl: "type.googleapis.com/protocol.TransferContract",
		Value:   tfcBytes,
	}

	tsc := &core.Transaction_Contract{
		Type:      core.Transaction_Contract_TransferContract,
		Parameter: &aa,
	}

	tscs := []*core.Transaction_Contract{tsc}

	block := client.GetNowBlock()

	fmt.Println(block.BlockHeader.RawData)
	blockHash := util.GetBlockHash(*block)
	//参考源码 setRefBlockHash(ByteString.copyFrom(ByteArray.subArray(blockHash, 8, 16))
	refBlockHash := blockHash[8:16]
	//rawByte := []byte(rawData)
	//获取当前区块时间戳
	timestamp := block.BlockHeader.RawData.Timestamp
	//设置过期时间
	var expiration int64
	expiration = timestamp + 10*60*60*1000
	blockHeight := block.BlockHeader.RawData.Number
	blockNum := int2Byte(blockHeight)
	//参考源码setRefBlockBytes(ByteString.copyFrom(ByteArray.subArray(refBlockNum, 6, 8))
	refBlockBytes := blockNum[6:8]

	tca := core.TransactionRaw{
		RefBlockBytes: refBlockBytes,
		RefBlockHash:  refBlockHash,
		Expiration:    expiration,
		Contract:      tscs,
		Timestamp:     time.Now().UnixNano() / 1000000,
	}
	tc := core.Transaction{
		RawData: &tca,
	}
	fmt.Println(tc)
	result := client.BroadcastTransaction(priv, &tc)
	fmt.Println(result)
}
