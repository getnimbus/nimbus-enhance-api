package entity

import (
	"encoding/hex"
	"encoding/json"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type BlockResult struct {
	Block           *types.Block                       `json:"block"`
	Transactions    []*TransactionResult               `json:"transactions"`
	transactionMaps map[common.Hash]*TransactionResult `json:"-"`
}

func (b *BlockResult) WithTransactions(transactionMaps map[common.Hash]*TransactionResult) *BlockResult {
	b.transactionMaps = transactionMaps
	b.toTransactions()
	return b
}

func (b *BlockResult) toTransactions() []*TransactionResult {
	b.Transactions = make([]*TransactionResult, 0, len(b.transactionMaps))
	for _, transaction := range b.transactionMaps {
		b.Transactions = append(b.Transactions, transaction)
	}
	return b.Transactions
}

type TransactionResult struct {
	Index       int                `json:"index"`
	ChainId     string             `json:"chain_id"`
	Transaction *types.Transaction `json:"transaction"`
	Message     *types.Message     `json:"message"`
	Receipt     *types.Receipt     `json:"receipt"`
}

type BlockRawLogs struct {
	BlockNumber  string               `json:"blockNumber"`
	BlockDate    string               `json:"blockDate"`
	Timestamp    string               `json:"timestamp"`
	GasLimit     string               `json:"gasLimit"`
	GasUsed      string               `json:"gasUsed"`
	LogsBloom    string               `json:"logsBloom"`
	Miner        string               `json:"miner"`
	Difficulty   string               `json:"difficulty"`
	Nonce        string               `json:"nonce"`
	Root         string               `json:"root"`
	ParentHash   string               `json:"parentHash"`
	Hash         string               `json:"hash"`
	ReceiptHash  string               `json:"receiptHash"`
	UncleHash    string               `json:"uncleHash"`
	ExtraData    string               `json:"extraData"`
	BaseFee      string               `json:"baseFee"`
	NumberOfTx   int                  `json:"numberOfTx"`
	Transactions []TransactionRawLogs `json:"transactions"`
}

func (b *BlockRawLogs) FromEntity(block *types.Block) *BlockRawLogs {
	b.BlockNumber = block.Number().String()
	b.Timestamp = strconv.FormatUint(block.Time(), 10)
	b.BlockDate = time.Unix(int64(block.Time()), 0).Format("2006-01-02")
	b.GasLimit = strconv.FormatUint(block.GasLimit(), 10)
	b.GasUsed = strconv.FormatUint(block.GasUsed(), 10)
	b.Miner = block.Coinbase().Hex()
	b.Difficulty = block.Difficulty().String()
	b.Nonce = strconv.FormatUint(block.Nonce(), 10)
	b.Root = block.Root().Hex()
	b.ParentHash = block.ParentHash().Hex()
	b.Hash = block.Hash().Hex()
	b.ReceiptHash = block.ReceiptHash().Hex()
	b.UncleHash = block.UncleHash().Hex()
	b.ExtraData = string(block.Extra())
	b.BaseFee = block.BaseFee().String()
	b.NumberOfTx = len(block.Transactions())
	logsBloom, err := block.Bloom().MarshalText()
	if err == nil {
		b.LogsBloom = string(logsBloom)
	}
	return b
}

func (b *BlockRawLogs) WithTransactions(data []*TransactionResult) *BlockRawLogs {
	transactionRawLogs := make([]TransactionRawLogs, 0, len(data))
	for _, item := range data {
		tx := new(TransactionRawLogs).FromEntity(item)
		transactionRawLogs = append(transactionRawLogs, *tx)
	}
	b.Transactions = transactionRawLogs
	return b
}

type TransactionRawLogs struct {
	Index      int            `json:"index"`
	Hash       string         `json:"hash"`
	From       string         `json:"from"`
	To         string         `json:"to"`
	Value      string         `json:"value"`
	GasPrice   string         `json:"gasPrice"`
	GasFeeCap  string         `json:"gasFeeCap"`
	GasTipCap  string         `json:"gasTipCap"`
	Gas        string         `json:"gas"`
	Nonce      string         `json:"nonce"`
	Data       string         `json:"data"`
	AccessList string         `json:"accessList"`
	ChainId    string         `json:"chainId"`
	V          string         `json:"v"`
	R          string         `json:"r"`
	S          string         `json:"s"`
	Receipt    ReceiptRawLogs `json:"receipt"`
}

func (tx *TransactionRawLogs) FromEntity(res *TransactionResult) *TransactionRawLogs {
	tx.Index = res.Index
	tx.Hash = res.Transaction.Hash().Hex()
	tx.From = res.Message.From().Hex()
	if res.Message.To() != nil {
		tx.To = res.Message.To().Hex()
	}
	tx.Value = res.Message.Value().String()
	tx.GasPrice = res.Message.GasPrice().String()
	tx.GasFeeCap = res.Message.GasFeeCap().String()
	tx.GasTipCap = res.Message.GasTipCap().String()
	tx.Gas = strconv.FormatUint(res.Message.Gas(), 10)
	tx.Nonce = strconv.FormatUint(res.Message.Nonce(), 10)
	tx.Data = hex.EncodeToString(res.Message.Data())
	accessList, err := json.Marshal(res.Message.AccessList())
	if err == nil {
		tx.AccessList = string(accessList)
	}
	tx.ChainId = res.ChainId
	v, r, s := res.Transaction.RawSignatureValues()
	tx.V = v.String()
	tx.R = r.String()
	tx.S = s.String()
	if res.Receipt != nil {
		receipt := new(ReceiptRawLogs).FromEntity(res.Receipt)
		tx.Receipt = *receipt
	}
	return tx
}

type ReceiptRawLogs struct {
	Type              string    `json:"type"`
	PostState         string    `json:"postState"`
	Status            string    `json:"status"`
	CumulativeGasUsed string    `json:"cumulativeGasUsed"`
	LogsBloom         string    `json:"logsBloom"`
	Logs              []LogData `json:"logs"`
	TxHash            string    `json:"txHash"`
	ContractAddress   string    `json:"contractAddress"`
	GasUsed           string    `json:"gasUsed"`
	BlockHash         string    `json:"blockHash"`
	BlockNumber       string    `json:"blockNumber"`
	TransactionIndex  int       `json:"transactionIndex"`
}

func (r *ReceiptRawLogs) FromEntity(receipt *types.Receipt) *ReceiptRawLogs {
	r.Type = strconv.Itoa(int(receipt.Type))
	r.PostState = string(receipt.PostState)
	r.Status = strconv.FormatUint(receipt.Status, 10)
	r.CumulativeGasUsed = strconv.FormatUint(receipt.CumulativeGasUsed, 10)
	logsBloom, err := json.Marshal(receipt.Bloom)
	if err == nil {
		r.LogsBloom = string(logsBloom)
	}
	if len(receipt.Logs) > 0 {
		r.Logs = make([]LogData, 0, len(receipt.Logs))
		for _, log := range receipt.Logs {
			r.Logs = append(r.Logs, *new(LogData).FromEntity(log))
		}
	}
	r.TxHash = receipt.TxHash.Hex()
	r.ContractAddress = receipt.ContractAddress.Hex()
	r.GasUsed = strconv.FormatUint(receipt.GasUsed, 10)
	r.BlockHash = receipt.BlockHash.Hex()
	r.BlockNumber = receipt.BlockNumber.String()
	r.TransactionIndex = int(receipt.TransactionIndex)
	return r
}

type LogData struct {
	Address          string   `json:"address"`
	Topics           []string `json:"topics"`
	Data             string   `json:"data"`
	BlockNumber      string   `json:"blockNumber"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
	BlockHash        string   `json:"blockHash"`
	LogIndex         string   `json:"logIndex"`
	Removed          bool     `json:"removed"`
}

func (d *LogData) FromEntity(log *types.Log) *LogData {
	data, err := log.MarshalJSON()
	if err == nil {
		json.Unmarshal(data, d)
	}
	return d
}
