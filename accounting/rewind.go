package accounting

import (
	"bytes"
	"context"

	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	atypes "github.com/algorand/go-algorand-sdk/types"
	models "github.com/algorand/indexer/api/generated/v2"
	log "github.com/sirupsen/logrus"

	"github.com/algorand/indexer/idb"
	"github.com/algorand/indexer/types"
)

func assetUpdate(account *models.Account, assetid uint64, add, sub uint64) {
	if account.Assets == nil {
		account.Assets = new([]models.AssetHolding)
	}
	assets := *account.Assets
	for i, ah := range assets {
		if ah.AssetId == assetid {
			ah.Amount += add
			ah.Amount -= sub
			assets[i] = ah
			// found and updated asset, done
			return
		}
	}
	// add asset to list
	assets = append(assets, models.AssetHolding{
		Amount:  add - sub,
		AssetId: assetid,
		//Creator: base32 addr string of asset creator, TODO
		//IsFrozen: leave nil? // TODO: on close record frozen state for rewind
	})
	*account.Assets = assets
}

func applyReverseDelta(ls *models.ApplicationLocalState, delta []idb.StateDelta) {
	var tkvs []models.TealKeyValue
	if ls.KeyValue != nil {
		tkvs = *ls.KeyValue
	}
	for _, d := range delta {
		found := false
		for i, kv := range tkvs {
			if bytes.Equal([]byte(kv.Key), d.Key) {
				found = true
				switch d.Delta.Action {
				case types.SetBytesAction:
					tkvs[i].Value.Bytes = string(d.Delta.Bytes)
					tkvs[i].Value.Type = 1 // TODO: add constants somewhere undere generated/v2 ? Is this in the sdk?
				case types.SetUintAction:
					tkvs[i].Value.Uint = d.Delta.Uint
					tkvs[i].Value.Type = 2
				case types.DeleteAction:
					if i > len(tkvs)-1 {
						tkvs[i] = tkvs[len(tkvs)-1]
					}
					tkvs = tkvs[:len(tkvs)-1]
				}
				break
			}
		}
		if !found {
			nkv := models.TealKeyValue{Key: string(d.Key)}
			switch d.Delta.Action {
			case types.SetBytesAction:
				nkv.Value.Bytes = string(d.Delta.Bytes)
				nkv.Value.Type = 1 // TODO: add constants somewhere undere generated/v2 ? Is this in the sdk?
			case types.SetUintAction:
				nkv.Value.Uint = d.Delta.Uint
				nkv.Value.Type = 2
			}
			tkvs = append(tkvs, nkv)
		}
	}
	if len(tkvs) > 0 {
		ntkvs := models.TealKeyValueStore(tkvs)
		ls.KeyValue = &ntkvs
	} else {
		ls.KeyValue = nil
	}
}

func appRewind(account *models.Account, txnrow *idb.TxnRow, stxn *types.SignedTxnWithAD) error {
	thisaddr, err := atypes.DecodeAddress(account.Address)
	if err != nil {
		return err
	}
	// TODO: rewind app state
	//txnrow.TxnExtra.GlobalReverseDelta
	//txnrow.TxnExtra.LocalReverseDelta
	// TODO: if this account is the owner, apply global delta

	var ls models.ApplicationLocalState
	var lsi int
	existingLocalState := false
	lsSet := false
	// find the local app state for this txn
	for lsi, ls = range *account.AppsLocalState {
		if atypes.AppIndex(ls.Id) == stxn.Txn.ApplicationID {
			existingLocalState = true
			break
		}
	}
	// for each local delta, if it applies to _this_ account, apply it
	for _, ld := range txnrow.Extra.LocalReverseDelta {
		var addr atypes.Address
		if ld.AddressIndex == 0 {
			addr = stxn.Txn.Sender
		} else {
			addr = stxn.Txn.Accounts[ld.AddressIndex-1]
		}
		if addr == thisaddr {
			lsSet = true
			applyReverseDelta(&ls, ld.Delta)
			log.Info("TODO WRITEME appRewind local ", string(idb.JsonOneLine(ld)))
		}
	}
	if !lsSet {
		// nothing happened
	} else if existingLocalState {
		(*account.AppsLocalState)[lsi] = ls
	} else {
		(*account.AppsLocalState) = append((*account.AppsLocalState), ls)
	}
	log.Info("TODO WRITEME appRewind", string(idb.JsonOneLine(txnrow.Extra.GlobalReverseDelta))) //, string(idb.JsonOneLine(txnrow.Extra.LocalReverseDelta)))
	return nil
}

func AccountAtRound(account models.Account, round uint64, db idb.IndexerDb) (acct models.Account, err error) {
	acct = account
	addr, err := atypes.DecodeAddress(account.Address)
	if err != nil {
		return
	}
	tf := idb.TransactionFilter{
		Address:  addr[:],
		MinRound: round + 1,
		MaxRound: account.Round,
	}
	txns := db.Transactions(context.Background(), tf)
	txcount := 0
	for txnrow := range txns {
		if txnrow.Error != nil {
			err = txnrow.Error
			return
		}
		txcount++
		var stxn types.SignedTxnWithAD
		err = msgpack.Decode(txnrow.TxnBytes, &stxn)
		if err != nil {
			return
		}
		if addr == stxn.Txn.Sender {
			acct.AmountWithoutPendingRewards += uint64(stxn.Txn.Fee)
			acct.AmountWithoutPendingRewards -= uint64(stxn.SenderRewards)
		}
		switch stxn.Txn.Type {
		case atypes.PaymentTx:
			if addr == stxn.Txn.Sender {
				acct.AmountWithoutPendingRewards += uint64(stxn.Txn.Amount)
			}
			if addr == stxn.Txn.Receiver {
				acct.AmountWithoutPendingRewards -= uint64(stxn.Txn.Amount)
				acct.AmountWithoutPendingRewards -= uint64(stxn.ReceiverRewards)
			}
			if addr == stxn.Txn.CloseRemainderTo {
				// unwind receiving a close-to
				acct.AmountWithoutPendingRewards -= uint64(stxn.ClosingAmount)
				acct.AmountWithoutPendingRewards -= uint64(stxn.CloseRewards)
			} else if !stxn.Txn.CloseRemainderTo.IsZero() {
				// unwind sending a close-to
				acct.AmountWithoutPendingRewards += uint64(stxn.ClosingAmount)
				acct.AmountWithoutPendingRewards += uint64(stxn.CloseRewards)
			}
		case atypes.KeyRegistrationTx:
			// TODO: keyreg does not rewind. workaround: query for txns on an account with typeenum=2 to find previous values it was set to.
		case atypes.AssetConfigTx:
			if stxn.Txn.ConfigAsset == 0 {
				// create asset, unwind the application of the value
				assetUpdate(&acct, txnrow.AssetId, 0, stxn.Txn.AssetParams.Total)
			}
		case atypes.AssetTransferTx:
			if addr == stxn.Txn.AssetSender || addr == stxn.Txn.Sender {
				assetUpdate(&acct, uint64(stxn.Txn.XferAsset), stxn.Txn.AssetAmount+txnrow.Extra.AssetCloseAmount, 0)
			}
			if addr == stxn.Txn.AssetReceiver {
				assetUpdate(&acct, uint64(stxn.Txn.XferAsset), 0, stxn.Txn.AssetAmount)
			}
			if addr == stxn.Txn.AssetCloseTo {
				assetUpdate(&acct, uint64(stxn.Txn.XferAsset), 0, txnrow.Extra.AssetCloseAmount)
			}
		case atypes.AssetFreezeTx:
			// TODO: mark an asset of the account as frozen or not?
		case atypes.ApplicationCallTx:
			err = appRewind(&acct, &txnrow, &stxn)
			if err != nil {
				return
			}
		default:
			panic("unknown txn type")
		}
	}

	if txcount > 0 {
		// If we found any txns above, we need to find one
		// more so we can know what the previous RewardsBase
		// of the account was so we can get the accurate
		// pending rewards at the target round.
		//
		// (If there weren't any txns above, the recorded
		// RewardsBase is current from whatever previous txn
		// happened to this account.)

		tf.MaxRound = round
		tf.MinRound = 0
		tf.Limit = 1
		txns = db.Transactions(context.Background(), tf)
		for txnrow := range txns {
			if txnrow.Error != nil {
				err = txnrow.Error
				return
			}
			var stxn types.SignedTxnWithAD
			err = msgpack.Decode(txnrow.TxnBytes, &stxn)
			if err != nil {
				return
			}
			var baseBlock types.Block
			baseBlock, err = db.GetBlock(txnrow.Round)
			if err != nil {
				return
			}
			prevRewardsBase := baseBlock.RewardsLevel
			var blockheader types.Block
			blockheader, err = db.GetBlock(round)
			if err != nil {
				return
			}
			var proto types.ConsensusParams
			proto, err = db.GetProto(string(blockheader.CurrentProtocol))
			if err != nil {
				return
			}
			rewardsUnits := acct.AmountWithoutPendingRewards / proto.RewardUnit
			rewardsDelta := blockheader.RewardsLevel - prevRewardsBase
			acct.PendingRewards = rewardsDelta * rewardsUnits
			acct.Amount = acct.PendingRewards + acct.AmountWithoutPendingRewards
			acct.Round = round
			return
		}

		// There were no prior transactions, it must have been empty before, zero out things
		acct.PendingRewards = 0
		acct.Amount = acct.AmountWithoutPendingRewards
	}

	acct.Round = round
	return
}
