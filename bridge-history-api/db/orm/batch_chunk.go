package orm

import (
	"database/sql"
	"fmt"

	"github.com/ethereum/go-ethereum/log"
	"github.com/jmoiron/sqlx"
)

type bridgeBatchOrm struct {
	db *sqlx.DB
}

type BridgeBatch struct {
	ID               uint64 `json:"id" db:"id"`
	BatchIndex       uint64 `json:"batch_index" db:"batch_index"`
	BatchHash        string `json:"batch_hash" db:"batch_hash"`
	Height           uint64 `json:"height" db:"height"`
	StartBlockNumber uint64 `json:"start_block_number" db:"start_block_number"`
	EndBlockNumber   uint64 `json:"end_block_number" db:"end_block_number"`
}

// NewBridgeBatchOrm create an NewBridgeBatchOrm instance
func NewBridgeBatchOrm(db *sqlx.DB) BridgeBatchOrm {
	return &bridgeBatchOrm{db: db}
}

func (b *bridgeBatchOrm) BatchInsertBridgeBatchDBTx(dbTx *sqlx.Tx, messages []*BridgeBatch) error {
	if len(messages) == 0 {
		return nil
	}
	var err error
	messageMaps := make([]map[string]interface{}, len(messages))
	for i, msg := range messages {
		messageMaps[i] = map[string]interface{}{
			"height":             msg.Height,
			"batch_index":        msg.BatchIndex,
			"batch_hash":         msg.BatchHash,
			"start_block_number": msg.StartBlockNumber,
			"end_block_number":   msg.EndBlockNumber,
		}
		var exists bool
		err := dbTx.QueryRow(`SELECT EXISTS(SELECT 1 FROM bridge_batch WHERE batch_index = $1 AND NOT is_deleted)`, msg.BatchIndex).Scan(&exists)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("BatchInsertBridgeBatchDBTx: batch index %v already exists at height %v", msg.BatchIndex, msg.Height)
		}

		_, err = dbTx.NamedExec(`insert into bridge_batch(height, batch_index, batch_hash, start_block_number, end_block_number) values(:height, :batch_index, :batch_hash, :start_block_number, :end_block_number);`, messageMaps[i])
		if err != nil {
			log.Error("BatchInsertBridgeBatchDBTx: failed to insert batch event msgs", "height", msg.Height, "batch_index", msg.BatchIndex, "batch_hash", msg.BatchHash, "start_block_number", msg.StartBlockNumber, "end_block_number", msg.EndBlockNumber)
			break
		}
	}
	return err
}

func (b *bridgeBatchOrm) GetLatestBridgeBatch() (*BridgeBatch, error) {
	result := &BridgeBatch{}
	row := b.db.QueryRowx(`SELECT id, batch_index, height, batch_hash, start_block_number, end_block_number FROM bridge_batch ORDER BY batch_index DESC LIMIT 1;`)
	if err := row.StructScan(result); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (b *bridgeBatchOrm) GetBridgeBatchByIndex(index uint64) (*BridgeBatch, error) {
	result := &BridgeBatch{}
	row := b.db.QueryRowx(`SELECT id, batch_index, batch_hash, height, start_block_number, end_block_number FROM bridge_batch WHERE batch_index = $1;`, index)
	if err := row.StructScan(result); err != nil {
		return nil, err
	}
	return result, nil
}