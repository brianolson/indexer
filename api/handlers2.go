package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/algorand/go-algorand-sdk/types"

	"github.com/algorand/indexer/accounting"
	"github.com/algorand/indexer/api/generated"
	"github.com/algorand/indexer/idb"
)

// ServerImplementation implements the handler interface used by the generated route definitions.
type ServerImplementation struct {
	// EnableAddressSearchRoundRewind is allows configuring whether or not the
	// 'accounts' endpoint allows specifying a round number. This is done for
	// performance reasons, because requesting many accounts at a particular
	// round could put a lot of strain on the system (especially if the round
	// is from long ago).
	EnableAddressSearchRoundRewind bool

	db idb.IndexerDb
}

/////////////////////////////
// Handler implementations //
/////////////////////////////

// LookupAccountByID queries indexer for a given account.
// (GET /account/{account-id})
func (si *ServerImplementation) LookupAccountByID(ctx echo.Context, accountID string, params generated.LookupAccountByIDParams) error {
	addr, errors := decodeAddress(&accountID, "account-id", make([]string, 0))
	if len(errors) != 0 {
		return badRequest(ctx, errors[0])
	}

	options := idb.AccountQueryOptions{
		EqualToAddress:       addr[:],
		IncludeAssetHoldings: true,
		IncludeAssetParams:   true,
		Limit:                1,
	}

	accounts, err := si.fetchAccounts(ctx.Request().Context(), options, params.Round)

	if err != nil {
		return indexerError(ctx, fmt.Sprintf("Failed while searching for account: %v", err))
	}

	if len(accounts) == 0 {
		return badRequest(ctx, fmt.Sprintf("No accounts found for address: %s", accountID))
	}

	if len(accounts) > 1 {
		return badRequest(ctx, fmt.Sprintf("Multiple accounts found for address, this shouldn't have happened: %s", accountID))
	}

	round, err := si.db.GetMaxRound()
	if err != nil {
		return indexerError(ctx, err.Error())
	}

	return ctx.JSON(http.StatusOK, generated.AccountResponse{
		CurrentRound: round,
		Account:      accounts[0],
	})
}

// SearchAccounts returns accounts matching the provided parameters
// (GET /accounts)
func (si *ServerImplementation) SearchAccounts(ctx echo.Context, params generated.SearchAccountsParams) error {
	options := idb.AccountQueryOptions{
		IncludeAssetHoldings: true,
		IncludeAssetParams:   true,
		Limit:                uintOrDefault(params.Limit),
		HasAssetId:           uintOrDefault(params.AssetId),
	}

	// Set GT/LT on Algos or Asset depending on whether or not an assetID was specified
	if options.HasAssetId == 0 {
		options.AlgosGreaterThan = uintOrDefault(params.CurrencyGreaterThan)
		options.AlgosLessThan = uintOrDefault(params.CurrencyLessThan)
	} else {
		options.AssetGT = uintOrDefault(params.CurrencyGreaterThan)
		options.AssetLT = uintOrDefault(params.CurrencyLessThan)
	}

	var atRound *uint64

	if si.EnableAddressSearchRoundRewind {
		atRound = params.Round
	}

	if params.Next != nil {
		addr, err := types.DecodeAddress(*params.Next)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, "Unable to parse next.")
		}
		options.GreaterThanAddress = addr[:]
	}

	accounts, err := si.fetchAccounts(ctx.Request().Context(), options, atRound)

	if err != nil {
		return badRequest(ctx, fmt.Sprintf("Failed while searching for account: %v", err))
	}

	round, err := si.db.GetMaxRound()
	if err != nil {
		return indexerError(ctx, err.Error())
	}

	// Set the next token if we hit the results limit
	// TODO: set the limit to +1, so we know that there are actually more results?
	var next *string
	if params.Limit != nil && uint64(len(accounts)) >= *params.Limit {
		next = strPtr(accounts[len(accounts)-1].Address)
	}

	response := generated.AccountsResponse{
		CurrentRound: round,
		NextToken:    next,
		Accounts:     accounts,
	}

	return ctx.JSON(http.StatusOK, response)
}

// LookupAccountTransactions looks up transactions associated with a particular account.
// (GET /account/{account-id}/transactions)
func (si *ServerImplementation) LookupAccountTransactions(ctx echo.Context, accountID string, params generated.LookupAccountTransactionsParams) error {
	// Check that a valid account was provided
	_, errors := decodeAddress(strPtr(accountID), "account-id", make([]string, 0))
	if len(errors) != 0 {
		return badRequest(ctx, errors[0])
	}

	searchParams := generated.SearchForTransactionsParams{
		Address: strPtr(accountID),
		// not applicable to this endpoint
		//AddressRole:         params.AddressRole,
		//ExcludeCloseTo:      params.ExcludeCloseTo,
		AssetId:             params.AssetId,
		Limit:               params.Limit,
		Next:                params.Next,
		NotePrefix:          params.NotePrefix,
		TxType:              params.TxType,
		SigType:             params.SigType,
		TxId:                params.TxId,
		Round:               params.Round,
		MinRound:            params.MinRound,
		MaxRound:            params.MaxRound,
		BeforeTime:          params.BeforeTime,
		AfterTime:           params.AfterTime,
		CurrencyGreaterThan: params.CurrencyGreaterThan,
		CurrencyLessThan:    params.CurrencyLessThan,
	}

	return si.SearchForTransactions(ctx, searchParams)
}

// LookupAssetByID looks up a particular asset
// (GET /asset/{asset-id})
func (si *ServerImplementation) LookupAssetByID(ctx echo.Context, assetID uint64) error {
	search := generated.SearchForAssetsParams{
		AssetId: uint64Ptr(assetID),
		Limit:   uint64Ptr(1),
	}
	options, err := assetParamsToAssetQuery(search)
	if err != nil {
		return badRequest(ctx, err.Error())
	}

	assets, err := si.fetchAssets(ctx.Request().Context(), options)
	if err != nil {
		return indexerError(ctx, err.Error())
	}

	if len(assets) == 0 {
		return badRequest(ctx, fmt.Sprintf("No assets found for id: %d", assetID))
	}

	if len(assets) > 1 {
		return badRequest(ctx, fmt.Sprintf("Multiple assets found for id, this shouldn't have happened: assetid=%d", assetID))
	}

	round, err := si.db.GetMaxRound()
	if err != nil {
		return indexerError(ctx, err.Error())
	}

	return ctx.JSON(http.StatusOK, generated.AssetResponse{
		Asset:        assets[0],
		CurrentRound: round,
	})
}

// LookupAssetBalances looks up balances for a particular asset
// (GET /asset/{asset-id}/balances)
func (si *ServerImplementation) LookupAssetBalances(ctx echo.Context, assetID uint64, params generated.LookupAssetBalancesParams) error {
	query := idb.AssetBalanceQuery{
		AssetId:   assetID,
		MinAmount: uintOrDefault(params.CurrencyGreaterThan),
		MaxAmount: uintOrDefault(params.CurrencyLessThan),
		Limit:     uintOrDefault(params.Limit),
	}

	if params.Next != nil {
		addr, err := types.DecodeAddress(*params.Next)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, "Unable to parse next.")
		}
		query.PrevAddress = addr[:]
	}

	balances, err := si.fetchAssetBalances(ctx.Request().Context(), query)
	if err != nil {
		indexerError(ctx, err.Error())
	}

	round, err := si.db.GetMaxRound()
	if err != nil {
		return indexerError(ctx, err.Error())
	}

	// Set the next token if we hit the results limit
	// TODO: set the limit to +1, so we know that there are actually more results?
	var next *string
	if params.Limit != nil && uint64(len(balances)) >= *params.Limit {
		next = strPtr(balances[len(balances)-1].Address)
	}

	return ctx.JSON(http.StatusOK, generated.AssetBalancesResponse{
		CurrentRound: round,
		NextToken:    next,
		Balances:     balances,
	})
}

// LookupAssetTransactions looks up transactions associated with a particular asset
// (GET /asset/{asset-id}/transactions)
func (si *ServerImplementation) LookupAssetTransactions(ctx echo.Context, assetID uint64, params generated.LookupAssetTransactionsParams) error {
	searchParams := generated.SearchForTransactionsParams{
		AssetId:             uint64Ptr(assetID),
		Limit:               params.Limit,
		Next:                params.Next,
		NotePrefix:          params.NotePrefix,
		TxType:              params.TxType,
		SigType:             params.SigType,
		TxId:                params.TxId,
		Round:               params.Round,
		MinRound:            params.MinRound,
		MaxRound:            params.MaxRound,
		BeforeTime:          params.BeforeTime,
		AfterTime:           params.AfterTime,
		CurrencyGreaterThan: params.CurrencyGreaterThan,
		CurrencyLessThan:    params.CurrencyLessThan,
		Address:             params.AddressRole,
		AddressRole:         params.AddressRole,
		ExcludeCloseTo:      params.ExcludeCloseTo,
	}

	return si.SearchForTransactions(ctx, searchParams)
}

// SearchForAssets returns assets matching the provided parameters
// (GET /assets)
func (si *ServerImplementation) SearchForAssets(ctx echo.Context, params generated.SearchForAssetsParams) error {
	options, err := assetParamsToAssetQuery(params)
	if err != nil {
		return badRequest(ctx, err.Error())
	}

	assets, err := si.fetchAssets(ctx.Request().Context(), options)
	if err != nil {
		return indexerError(ctx, err.Error())
	}

	round, err := si.db.GetMaxRound()
	if err != nil {
		return indexerError(ctx, err.Error())
	}

	// Set the next token if we hit the results limit
	// TODO: set the limit to +1, so we know that there are actually more results?
	var next *string
	if params.Limit != nil && uint64(len(assets)) >= *params.Limit {
		next = strPtr(strconv.FormatUint(assets[len(assets)-1].Index, 10))
	}

	return ctx.JSON(http.StatusOK, generated.AssetsResponse{
		CurrentRound: round,
		NextToken:    next,
		Assets:       assets,
	})
}

// LookupBlock returns the block for a given round number
// (GET /block/{round-number})
func (si *ServerImplementation) LookupBlock(ctx echo.Context, roundNumber uint64) error {
	blk, err := si.fetchBlock(roundNumber)
	if err != nil {
		return indexerError(ctx, err.Error())
	}

	// Lookup transactions
	filter := idb.TransactionFilter{Round: uint64Ptr(roundNumber)}
	txns, _, err := si.fetchTransactions(ctx.Request().Context(), filter)
	if err != nil {
		return indexerError(ctx, fmt.Sprintf("error while looking up for transactions for round '%d': %v", roundNumber, err))
	}

	blk.Transactions = &txns
	return ctx.JSON(http.StatusOK, generated.BlockResponse(blk))
}

// SearchForTransactions returns transactions matching the provided parameters
// (GET /transactions)
func (si *ServerImplementation) SearchForTransactions(ctx echo.Context, params generated.SearchForTransactionsParams) error {
	filter, err := transactionParamsToTransactionFilter(params)
	if err != nil {
		return badRequest(ctx, err.Error())
	}

	// Fetch the transactions
	txns, next, err := si.fetchTransactions(ctx.Request().Context(), filter)

	if err != nil {
		return indexerError(ctx, fmt.Sprintf("error while searching for transactions: %v", err))
	}

	round, err := si.db.GetMaxRound()
	if err != nil {
		return indexerError(ctx, err.Error())
	}

	response := generated.TransactionsResponse{
		CurrentRound: round,
		NextToken:    strPtr(next),
		Transactions: txns,
	}

	return ctx.JSON(http.StatusOK, response)
}

///////////////////
// Error Helpers //
///////////////////

// badRequest is a simple helper to return a 400 error.
func badRequest(ctx echo.Context, err string) error {
	return ctx.JSON(http.StatusBadRequest, generated.Error{
		Error: err,
	})
}

// indexerRequest is a simple helper to return a 500 error.
func indexerError(ctx echo.Context, err string) error {
	return ctx.JSON(http.StatusInternalServerError, generated.Error{
		Error: err,
	})
}

///////////////////////
// IndexerDb helpers //
///////////////////////

// fetchAssets fetches all results and converts them into generated.Asset objects
func (si *ServerImplementation) fetchAssets(ctx context.Context, options idb.AssetsQuery) ([]generated.Asset, error) {
	assetchan := si.db.Assets(ctx, options)
	assets := make([]generated.Asset, 0)
	for row := range assetchan {
		if row.Error != nil {
			return nil, row.Error
		}

		creator := types.Address{}
		if len(row.Creator) != len(creator) {
			return nil, fmt.Errorf("found an invalid creator address")
		}
		copy(creator[:], row.Creator[:])

		asset := generated.Asset{
			Index: row.AssetId,
			Params: generated.AssetParams{
				Creator:       creator.String(),
				Name:          strPtr(row.Params.AssetName),
				UnitName:      strPtr(row.Params.UnitName),
				Url:           strPtr(row.Params.URL),
				Total:         row.Params.Total,
				Decimals:      uint64(row.Params.Decimals),
				DefaultFrozen: boolPtr(row.Params.DefaultFrozen),
				MetadataHash:  bytePtr(row.Params.MetadataHash[:]),
				Clawback:      strPtr(row.Params.Clawback.String()),
				Reserve:       strPtr(row.Params.Reserve.String()),
				Freeze:        strPtr(row.Params.Freeze.String()),
				Manager:       strPtr(row.Params.Manager.String()),
			},
		}

		assets = append(assets, asset)
	}
	return assets, nil
}

// fetchAssetBalances fetches all balances from a query and converts them into
// generated.MiniAssetHolding objects
func (si *ServerImplementation) fetchAssetBalances(ctx context.Context, options idb.AssetBalanceQuery) ([]generated.MiniAssetHolding, error) {
	assetbalchan := si.db.AssetBalances(ctx, options)
	balances := make([]generated.MiniAssetHolding, 0)
	for row := range assetbalchan {
		if row.Error != nil {
			return nil, row.Error
		}

		addr := types.Address{}
		if len(row.Address) != len(addr) {
			return nil, fmt.Errorf("found an invalid creator address")
		}
		copy(addr[:], row.Address[:])

		bal := generated.MiniAssetHolding{
			Address:  addr.String(),
			Amount:   row.Amount,
			IsFrozen: row.Frozen,
		}

		balances = append(balances, bal)
	}

	return balances, nil
}

// fetchBlock looks up a block and converts it into a generated.Block object
func (si *ServerImplementation) fetchBlock(round uint64) (generated.Block, error) {
	blk, err := si.db.GetBlock(round)
	if err != nil {
		return generated.Block{}, fmt.Errorf("error while looking up for block for round '%d': %v", round, err)
	}

	rewards := generated.BlockRewards{
		FeeSink:                 "",
		RewardsCalculationRound: uint64(blk.RewardsRecalculationRound),
		RewardsLevel:            blk.RewardsLevel,
		RewardsPool:             blk.RewardsPool.String(),
		RewardsRate:             blk.RewardsRate,
		RewardsResidue:          blk.RewardsResidue,
	}

	upgradeState := generated.BlockUpgradeState{
		CurrentProtocol:        string(blk.CurrentProtocol),
		NextProtocol:           strPtr(string(blk.NextProtocol)),
		NextProtocolApprovals:  uint64Ptr(blk.NextProtocolApprovals),
		NextProtocolSwitchOn:   uint64Ptr(uint64(blk.NextProtocolSwitchOn)),
		NextProtocolVoteBefore: uint64Ptr(uint64(blk.NextProtocolVoteBefore)),
	}

	upgradeVote := generated.BlockUpgradeVote{
		UpgradeApprove: boolPtr(blk.UpgradeApprove),
		UpgradeDelay:   uint64Ptr(uint64(blk.UpgradeDelay)),
		UpgradePropose: strPtr(string(blk.UpgradePropose)),
	}

	ret := generated.Block{
		GenesisHash:       blk.GenesisHash[:],
		GenesisId:         blk.GenesisID,
		PreviousBlockHash: blk.Branch[:],
		Rewards:           &rewards,
		Round:             uint64(blk.Round),
		Seed:              blk.Seed[:],
		Timestamp:         uint64(blk.TimeStamp),
		Transactions:      nil,
		TransactionsRoot:  blk.TxnRoot[:],
		TxnCounter:        uint64Ptr(blk.TxnCounter),
		UpgradeState:      &upgradeState,
		UpgradeVote:       &upgradeVote,
	}

	return ret, nil
}

// fetchAccounts queries for accounts and converts them into generated.Account
// objects, optionally rewinding their value back to a particular round.
func (si *ServerImplementation) fetchAccounts(ctx context.Context, options idb.AccountQueryOptions, atRound *uint64) ([]generated.Account, error) {
	accountchan := si.db.GetAccounts(ctx, options)

	accounts := make([]generated.Account, 0)
	for row := range accountchan {
		if row.Error != nil {
			return nil, row.Error
		}

		// Compute for a given round if requested.
		var account generated.Account
		if atRound != nil {
			acct, err := accounting.AccountAtRound(row.Account, *atRound, si.db)
			if err != nil {
				return nil, fmt.Errorf("problem computing account at round: %v", err)
			}
			account = acct
		} else {
			account = row.Account
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

// fetchTransactions is used to query the backend for transactions, and compute the next token
func (si *ServerImplementation) fetchTransactions(ctx context.Context, filter idb.TransactionFilter) ([]generated.Transaction, string, error) {
	results := make([]generated.Transaction, 0)
	txchan := si.db.Transactions(ctx, filter)
	nextToken := ""
	for txrow := range txchan {
		tx, err := txnRowToTransaction(txrow)
		if err != nil {
			return nil, "", err
		}
		results = append(results, tx)
		nextToken = txrow.Next()
	}

	return results, nextToken, nil
}
