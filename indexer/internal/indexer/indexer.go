package indexer

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"w3todo-indexer/internal/config"
	"w3todo-indexer/internal/db"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	TransferEvent    = common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
	TodoCreatedEvent = common.HexToHash("0x3b1896bf181bf210836ea70a0a32a0bf2ed762ead0f6abb599c894e477434ff8")
	RewardPaidEvent  = common.HexToHash("0xe2403640ba68fed3a2f88b7557551d1993f84b99bb10ff833f0cf8db0c5e0486")
)

type Indexer struct {
	client   *ethclient.Client
	db       *db.Postgres
	token    common.Address
	todoList common.Address
}

func New(ctx context.Context, cfg config.Config, database *db.Postgres) (*Indexer, error) {
	client, err := ethclient.Dial(cfg.RPC)
	if err != nil {
		return nil, fmt.Errorf("dial rpc: %w", err)
	}
	return &Indexer{
		client:   client,
		db:       database,
		token:    common.HexToAddress(cfg.Token),
		todoList: common.HexToAddress(cfg.TodoList),
	}, nil
}

func (idx *Indexer) Start(ctx context.Context) {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{idx.token, idx.todoList},
		Topics:    [][]common.Hash{{TransferEvent, TodoCreatedEvent, RewardPaidEvent}},
	}

	logs := make(chan types.Log)
	sub, err := idx.client.SubscribeFilterLogs(ctx, query, logs)
	if err != nil {
		log.Printf("Ошибка подписки: %v", err)
		return
	}
	defer sub.Unsubscribe()

	fmt.Println("Индексатор слушает события...")

	for {
		select {
		case <-ctx.Done():
			return
		case err := <-sub.Err():
			log.Printf("Ошибка подписки: %v", err)
			return
		case vLog := <-logs:
			idx.handleEvent(vLog)
		}
	}
}

func (idx *Indexer) handleEvent(vLog types.Log) {
	switch vLog.Topics[0] {
	case TransferEvent:
		from := common.BytesToAddress(vLog.Topics[1].Bytes())
		to := common.BytesToAddress(vLog.Topics[2].Bytes())
		value := new(big.Int).SetBytes(vLog.Data)

		if err := idx.db.SaveTransfer(from.Hex(), to.Hex(), value, vLog.TxHash.Hex(), vLog.BlockNumber); err != nil {
			log.Printf("Ошибка сохранения трансфера: %v", err)
		}
		fmt.Printf("[Transfer] %s -> %s : %s\n", from.Hex(), to.Hex(), value.String())

	case TodoCreatedEvent:
		id := new(big.Int).SetBytes(vLog.Topics[1].Bytes())
		owner := common.BytesToAddress(vLog.Topics[2].Bytes())
		text := parseString(vLog.Data)

		if err := idx.db.SaveTodo(id.Uint64(), text, owner.Hex(), vLog.TxHash.Hex(), vLog.BlockNumber); err != nil {
			log.Printf("Ошибка сохранения todo: %v", err)
		}
		fmt.Printf("[Todo] #%d: %s (by %s)\n", id.Uint64(), text, owner.Hex())

	case RewardPaidEvent:
		user := common.BytesToAddress(vLog.Topics[1].Bytes())
		amount := new(big.Int).SetBytes(vLog.Data)

		if err := idx.db.SaveReward(user.Hex(), amount, vLog.TxHash.Hex(), vLog.BlockNumber); err != nil {
			log.Printf("Ошибка сохранения награды: %v", err)
		}
		fmt.Printf("[Reward] %s получил %s токенов\n", user.Hex(), amount.String())
	}
}

func parseString(data []byte) string {
	if len(data) < 64 {
		return ""
	}
	raw := data[64:]
	for i, b := range raw {
		if b == 0 {
			return string(raw[:i])
		}
	}
	return string(raw)
}
