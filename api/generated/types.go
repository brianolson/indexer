// Package generated provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package generated

// Account defines model for Account.
type Account struct {

	// the account public key
	Address string `json:"address"`

	// \[algo\] total number of MicroAlgos in the account
	Amount uint64 `json:"amount"`

	// specifies the amount of MicroAlgos in the account, without the pending rewards.
	AmountWithoutPendingRewards uint64 `json:"amount-without-pending-rewards"`

	// \[asset\] assets held by this account.
	//
	// Note the raw object uses `map[int] -> AssetHolding` for this type.
	Assets *[]AssetHolding `json:"assets,omitempty"`

	// \[apar\] parameters of assets created by this account.
	//
	// Note: the raw account uses `map[int] -> Asset` for this type.
	CreatedAssets *[]Asset `json:"created-assets,omitempty"`
	Participation *struct {

		// \[sel\] Selection public key (if any) currently registered for this round.
		SelectionParticipationKey *string `json:"selection-participation-key,omitempty"`

		// \[voteFst\] First round for which this participation is valid.
		VoteFirstValid *uint64 `json:"vote-first-valid,omitempty"`

		// \[voteKD\] Number of subkeys in each batch of participation keys.
		VoteKeyDilution *uint64 `json:"vote-key-dilution,omitempty"`

		// \[voteLst\] Last round for which this participation is valid.
		VoteLastValid *uint64 `json:"vote-last-valid,omitempty"`

		// \[vote\] root participation public key (if any) currently registered for this round.
		VoteParticipationKey *string `json:"vote-participation-key,omitempty"`
	} `json:"participation,omitempty"`

	// amount of MicroAlgos of pending rewards in this account.
	PendingRewards uint64 `json:"pending-rewards"`

	// \[ebase\] used as part of the rewards computation. Only applicable to accounts which are participating.
	RewardBase *uint64 `json:"reward-base,omitempty"`

	// \[ern\] total rewards of MicroAlgos the account has received, including pending rewards.
	Rewards uint64 `json:"rewards"`

	// The round for which this information is relevant.
	Round uint64 `json:"round"`

	// \[onl\] delegation status of the account's MicroAlgos
	// * Offline - indicates that the associated account is delegated.
	// *  Online  - indicates that the associated account used as part of the delegation pool.
	// *   NotParticipating - indicates that the associated account is neither a delegator nor a delegate.
	Status string `json:"status"`

	// Indicates what type of signature is used by this account, must be one of:
	// * standard
	// * logic
	// * multisig
	Type *string `json:"type,omitempty"`
}

// Asset defines model for Asset.
type Asset struct {

	// unique asset identifier
	Index uint64 `json:"index"`

	// AssetParams specifies the parameters for an asset.
	//
	// \[apar\] when part of an AssetConfig transaction.
	//
	// Definition:
	// data/transactions/asset.go : AssetParams
	Params AssetParams `json:"params"`
}

// AssetHolding defines model for AssetHolding.
type AssetHolding struct {

	// \[a\] number of units held.
	Amount uint64 `json:"amount"`

	// Asset ID of the holding.
	AssetId string `json:"asset-id"`

	// Address that created this asset. This is the address where the parameters for this asset can be found, and also the address where unwanted asset units can be sent in the worst case.
	Creator string `json:"creator"`

	// \[f\] whether or not the holding is frozen.
	IsFrozen *bool `json:"is-frozen,omitempty"`
}

// AssetParams defines model for AssetParams.
type AssetParams struct {

	// \[c\] Address of account used to clawback holdings of this asset.  If empty, clawback is not permitted.
	Clawback *string `json:"clawback,omitempty"`

	// The address that created this asset. This is the address where the parameters for this asset can be found, and also the address where unwanted asset units can be sent in the worst case.
	Creator string `json:"creator"`

	// \[dc\] The number of digits to use after the decimal point when displaying this asset. If 0, the asset is not divisible. If 1, the base unit of the asset is in tenths. If 2, the base unit of the asset is in hundredths, and so on. This value must be between 0 and 19 (inclusive).
	Decimals uint64 `json:"decimals"`

	// \[df\] Whether holdings of this asset are frozen by default.
	DefaultFrozen *bool `json:"default-frozen,omitempty"`

	// \[f\] Address of account used to freeze holdings of this asset.  If empty, freezing is not permitted.
	Freeze *string `json:"freeze,omitempty"`

	// \[m\] Address of account used to manage the keys of this asset and to destroy it.
	Manager *string `json:"manager,omitempty"`

	// \[am\] A commitment to some unspecified asset metadata. The format of this metadata is up to the application.
	MetadataHash *string `json:"metadata-hash,omitempty"`

	// \[an\] Name of this asset, as supplied by the creator.
	Name *string `json:"name,omitempty"`

	// \[r\] Address of account holding reserve (non-minted) units of this asset.
	Reserve *string `json:"reserve,omitempty"`

	// \[t\] The total number of units of this asset.
	Total uint64 `json:"total"`

	// \[un\] Name of a unit of this asset, as supplied by the creator.
	UnitName *string `json:"unit-name,omitempty"`

	// \[au\] URL where more information about the asset can be retrieved.
	Url *string `json:"url,omitempty"`
}

// Block defines model for Block.
type Block struct {

	// \[gh\] hash to which this block belongs.
	GenesisHash string `json:"genesis-hash"`

	// \[gen\] ID to which this block belongs.
	GenesisId string `json:"genesis-id"`

	// Current block hash
	Hash string `json:"hash"`

	// Period on which the block was confirmed.
	Period uint64 `json:"period"`

	// \[prev\] Previous block hash.
	PreviousBlockHash string `json:"previous-block-hash"`

	// Address that proposed this block.
	Proposer string `json:"proposer"`

	// Fields relating to rewards,
	Rewards *struct {

		// \[fees\] accepts transaction fees, it can only spend to the incentive pool.
		FeeSink string `json:"fee-sink"`

		// \[rwcalr\] number of leftover MicroAlgos after the distribution of rewards-rate MicroAlgos for every reward unit in the next round.
		RewardsCalculationRound uint64 `json:"rewards-calculation-round"`

		// \[earn\] How many rewards, in MicroAlgos, have been distributed to each RewardUnit of MicroAlgos since genesis.
		RewardsLevel uint64 `json:"rewards-level"`

		// \[rwd\] accepts periodic injections from the fee-sink and continually redistributes them as rewards.
		RewardsPool string `json:"rewards-pool"`

		// \[rate\] Number of new MicroAlgos added to the participation stake from rewards at the next round.
		RewardsRate uint64 `json:"rewards-rate"`

		// \[frac\] Number of leftover MicroAlgos after the distribution of RewardsRate/rewardUnits MicroAlgos for every reward unit in the next round.
		RewardsResidue uint64 `json:"rewards-residue"`
	} `json:"rewards,omitempty"`

	// \[rnd\] Current round on which this block was appended to the chain.
	Round uint64 `json:"round"`

	// \[seed\] Sortition seed.
	Seed string `json:"seed"`

	// \[ts\] Block creation timestamp in seconds since eposh
	Timestamp uint64 `json:"timestamp"`

	// \[txns\] list of transactions corresponding to a given round.
	Transactions *[]Transaction `json:"transactions,omitempty"`

	// \[txn\] TransactionsRoot authenticates the set of transactions appearing in the block. More specifically, it's the root of a merkle tree whose leaves are the block's Txids, in lexicographic order. For the empty block, it's 0. Note that the TxnRoot does not authenticate the signatures on the transactions, only the transactions themselves. Two blocks with the same transactions but in a different order and with different signatures will have the same TxnRoot.
	TransactionsRoot string `json:"transactions-root"`

	// \[tc\] TxnCounter counts the number of transactions committed in the ledger, from the time at which support for this feature was introduced.
	//
	// Specifically, TxnCounter is the number of the next transaction that will be committed after this block.  It is 0 when no transactions have ever been committed (since TxnCounter started being supported).
	TxnCounter *uint64 `json:"txn-counter,omitempty"`

	// Fields relating to a protocol upgrade.
	UpgradeState *struct {

		// \[proto\] The current protocol version.
		CurrentProtocol string `json:"current-protocol"`

		// \[nextproto\] The next proposed protocol version.
		NextProtocol *string `json:"next-protocol,omitempty"`

		// \[nextyes\] Number of blocks which approved the protocol upgrade.
		NextProtocolApprovals *uint64 `json:"next-protocol-approvals,omitempty"`

		// \[nextswitch\] Round on which the protocol upgrade will take effect.
		NextProtocolSwitchOn *uint64 `json:"next-protocol-switch-on,omitempty"`

		// \[nextbefore\] Deadline round for this protocol upgrade (No votes will be consider after this round).
		NextProtocolVoteBefore *uint64 `json:"next-protocol-vote-before,omitempty"`
	} `json:"upgrade-state,omitempty"`

	// Fields relating to voting for a protocol upgrade.
	UpgradeVote *struct {

		// \[upgradeyes\] Indicates a yes vote for the current proposal.
		UpgradeApprove *bool `json:"upgrade-approve,omitempty"`

		// \[upgradedelay\] Indicates the time between acceptance and execution.
		UpgradeDelay *uint64 `json:"upgrade-delay,omitempty"`

		// \[upgradeprop\] Indicates a proposed upgrade.
		UpgradePropose *string `json:"upgrade-propose,omitempty"`
	} `json:"upgrade-vote,omitempty"`
}

// MiniAssetHolding defines model for MiniAssetHolding.
type MiniAssetHolding struct {
	Address  *string `json:"address,omitempty"`
	Amount   *uint64 `json:"amount,omitempty"`
	IsFrozen *bool   `json:"isFrozen,omitempty"`
}

// NodeStatus defines model for NodeStatus.
type NodeStatus struct {

	// CatchupTime in nanoseconds
	CatchupTime uint64 `json:"catchupTime"`

	// HasSyncedSinceStartup indicates whether a round has completed since startup
	HasSyncedSinceStartup bool `json:"hasSyncedSinceStartup"`

	// LastVersion indicates the last consensus version supported
	LastConsensusVersion string `json:"lastConsensusVersion"`

	// LastRound indicates the last round seen
	LastRound uint64 `json:"lastRound"`

	// NextVersion of consensus protocol to use
	NextConsensusVersion string `json:"nextConsensusVersion"`

	// NextVersionRound is the round at which the next consensus version will apply
	NextConsensusVersionRound uint64 `json:"nextConsensusVersionRound"`

	// NextVersionSupported indicates whether the next consensus version is supported by this node
	NextConsensusVersionSupported bool `json:"nextConsensusVersionSupported"`

	// StoppedAtUnsupportedRound indicates that the node does not support the new rounds and has stopped making progress
	StoppedAtUnsupportedRound bool `json:"stoppedAtUnsupportedRound"`

	// TimeSinceLastRound in nanoseconds
	TimeSinceLastRound uint64 `json:"timeSinceLastRound"`
}

// Supply defines model for Supply.
type Supply struct {

	// OnlineMoney
	OnlineMoney uint64 `json:"onlineMoney"`

	// Round
	Round uint64 `json:"round"`

	// TotalMoney
	TotalMoney uint64 `json:"totalMoney"`
}

// Transaction defines model for Transaction.
type Transaction struct {

	// Fields for asset allocation, re-configuration, and destruction.
	//
	//
	// A zero value for asset-id indicates asset creation.
	// A zero value for the params indicates asset destruction.
	//
	// Definition:
	// data/transactions/asset.go : AssetConfigTxnFields
	AssetConfigTransaction *struct {

		// \[xaid\] ID of the asset being configured or empty if creating.
		AssetId *uint64 `json:"asset-id,omitempty"`

		// AssetParams specifies the parameters for an asset.
		//
		// \[apar\] when part of an AssetConfig transaction.
		//
		// Definition:
		// data/transactions/asset.go : AssetParams
		Params *AssetParams `json:"params,omitempty"`
	} `json:"asset-config-transaction,omitempty"`

	// Fields for an asset freeze transaction.
	//
	// Definition:
	// data/transactions/asset.go : AssetFreezeTxnFields
	AssetFreezeTransaction *struct {

		// \[fadd\] Address of the account whose asset is being frozen or thawed.
		Address string `json:"address"`

		// \[faid\] ID of the asset being frozen or thawed.
		AssetId uint64 `json:"asset-id"`

		// \[afrz\] The new freeze status.
		NewFreezeStatus bool `json:"new-freeze-status"`
	} `json:"asset-freeze-transaction,omitempty"`

	// Fields for an asset transfer transaction.
	//
	// Definition:
	// data/transactions/asset.go : AssetTransferTxnFields
	AssetTransferTransaction *struct {

		// \[aamt\] Amount of asset to transfer. A zero amount transferred to self allocates that asset in the account's Assets map.
		Amount uint64 `json:"amount"`

		// \[xaid\] ID of the asset being transferred.
		AssetId uint64 `json:"asset-id"`

		// \[aclose\] Indicates that the asset should be removed from the account's Assets map, and specifies where the remaining asset holdings should be transferred.  It's always valid to transfer remaining asset holdings to the creator account.
		CloseTo *string `json:"close-to,omitempty"`

		// \[arcv\] Recipient address of the transfer.
		Receiver string `json:"receiver"`

		// \[asnd\] Sender of the transfer.  If this is not a zero value, the real transaction sender must be the Clawback address from the AssetParams.  If this is the zero value, the asset is sent from the transaction's Sender.
		Sender *string `json:"sender,omitempty"`
	} `json:"asset-transfer-transaction,omitempty"`

	// \[rc\] rewards applied to close-remainder-to account.
	CloseRewards *uint64 `json:"close-rewards,omitempty"`

	// \[ca\] closing amount for transaction.
	ClosingAmount *uint64 `json:"closing-amount,omitempty"`

	// Round when the transaction was confirmed.
	ConfirmedRound *uint64 `json:"confirmed-round,omitempty"`

	// Specifies an asset index (ID) if an asset was created with this transaction.
	CreatedAssetIndex *uint64 `json:"created-asset-index,omitempty"`

	// \[fee\] Transaction fee.
	Fee *uint64 `json:"fee,omitempty"`

	// \[fv\] First valid round for this transaction.
	FirstValid *uint64 `json:"first-valid,omitempty"`

	// \[gh\] Hash of genesis block.
	GenesisHash *string `json:"genesis-hash,omitempty"`

	// \[gen\] genesis block ID.
	GenesisId *string `json:"genesis-id,omitempty"`

	// \[grp\] Base64 encoded byte array of a sha512/256 digest. When present indicates that this transaction is part of a transaction group and the value is the sha512/256 hash of the transactions in that group.
	Group *string `json:"group,omitempty"`

	// \[hgh\]
	HasGenesisHash *bool `json:"has-genesis-hash,omitempty"`

	// \[hgi\]
	HasGenesisId *bool `json:"has-genesis-id,omitempty"`

	// Transaction ID
	Id *string `json:"id,omitempty"`

	// Fields for a keyreg transaction.
	//
	// Definition:
	// data/transactions/keyreg.go : KeyregTxnFields
	KeyregTransaction *struct {

		// \[nonpart\] Mark the account as participating or non-participating.
		NonParticipation *bool `json:"non-participation,omitempty"`

		// \[selkey\] Public key used with the Verified Random Function (VRF) result during committee selection.
		SelectionParticipationKey *string `json:"selection-participation-key,omitempty"`

		// \[votefst\] First round this participation key is valid.
		VoteFirstValid *uint64 `json:"vote-first-valid,omitempty"`

		// \[votekd\] Number of subkeys in each batch of participation keys.
		VoteKeyDilution *uint64 `json:"vote-key-dilution,omitempty"`

		// \[votelst\] Last round this participation key is valid.
		VoteLastValid *uint64 `json:"vote-last-valid,omitempty"`

		// \[votekey\] Participation public key used in key registration transactions.
		VoteParticipationKey *string `json:"vote-participation-key,omitempty"`
	} `json:"keyreg-transaction,omitempty"`

	// \[lv\] Last valid round for this transaction.
	LastValid *uint64 `json:"last-valid,omitempty"`

	// \[lx\] Base64 encoded 32-byte array. Lease enforces mutual exclusion of transactions.  If this field is nonzero, then once the transaction is confirmed, it acquires the lease identified by the (Sender, Lease) pair of the transaction until the LastValid round passes.  While this transaction possesses the lease, no other transaction specifying this lease can be confirmed.
	Lease *string `json:"lease,omitempty"`

	// \[note\] Free form data.
	Note *string `json:"note,omitempty"`

	// Fields for a payment transaction.
	//
	// Definition:
	// data/transactions/payment.go : PaymentTxnFields
	PaymentTransaction *struct {

		// \[amt\] number of MicroAlgos intended to be transferred.
		Amount uint64 `json:"amount"`

		// Number of MicroAlgos that were sent to the close-remainder-to address when closing the sender account.
		CloseAmount *uint64 `json:"close-amount,omitempty"`

		// \[close\] when set, indicates that the sending account should be closed and all remaining funds be transferred to this address.
		CloseRemainderTo *string `json:"close-remainder-to,omitempty"`

		// \[rcv\] receiver's address.
		Receiver string `json:"receiver"`
	} `json:"payment-transaction,omitempty"`

	// Part of algod API only.
	//
	// Indicates the transaction was evicted from this node's transaction pool (if non-empty).  A non-empty PoolError does not guarantee that the transaction will never be committed; other nodes may not have evicted the transaction and may attempt to commit it in the future.
	PoolError *string `json:"pool-error,omitempty"`

	// \[rr\] rewards applied to receiver account.
	ReceiverRewards *uint64 `json:"receiver-rewards,omitempty"`

	// \[snd\] Sender's address.
	Sender *string `json:"sender,omitempty"`

	// \[rs\] rewards applied to sender account.
	SenderRewards *uint64 `json:"sender-rewards,omitempty"`

	// Validation signature associated with some data. Only one of the signatures should be provided.
	Signature *struct {

		// \[lsig\] Programatic transaction signature.
		//
		// Definition:
		// data/transactions/logicsig.go
		Logicsig *struct {

			// \[arg\] Logic arguments, base64 encoded.
			Args *[]string `json:"args,omitempty"`

			// \[l\] Program signed by a signature or multi signature, or hashed to be the address of ana ccount. Base64 encoded TEAL program.
			Logic *string `json:"logic,omitempty"`

			// \[msig\] structure holding multiple subsignatures.
			//
			// Definition:
			// crypto/multisig.go : MultisigSig
			MultisigSignature *struct {

				// \[subsig\] holds pairs of public key and signatures.
				Subsignature *[]struct {

					// \[pk\]
					PublicKey *string `json:"public-key,omitempty"`

					// \[s\]
					Signature *string `json:"signature,omitempty"`
				} `json:"subsignature,omitempty"`

				// \[thr\]
				Threshold *uint64 `json:"threshold,omitempty"`

				// \[v\]
				Version *uint64 `json:"version,omitempty"`
			} `json:"multisig-signature,omitempty"`

			// \[sig\] ed25519 signature.
			Signature *string `json:"signature,omitempty"`
		} `json:"logicsig,omitempty"`

		// \[msig\] structure holding multiple subsignatures.
		//
		// Definition:
		// crypto/multisig.go : MultisigSig
		Multisig *struct {

			// \[subsig\] holds pairs of public key and signatures.
			Subsignature *[]struct {

				// \[pk\]
				PublicKey *string `json:"public-key,omitempty"`

				// \[s\]
				Signature *string `json:"signature,omitempty"`
			} `json:"subsignature,omitempty"`

			// \[thr\]
			Threshold *uint64 `json:"threshold,omitempty"`

			// \[v\]
			Version *uint64 `json:"version,omitempty"`
		} `json:"multisig,omitempty"`

		// \[sig\] Standard ed25519 signature.
		Sig *string `json:"sig,omitempty"`
	} `json:"signature,omitempty"`

	// \[type\] Indicates what type of transaction this is. Different types have different fields.
	//
	// Valid types, and where their fields are stored:
	// * \[pay\] payment-transaction
	// * \[keyreg\] keyreg-transaction
	// * \[acfg\] asset-config-transaction
	// * \[axfer\] asset-transfer-transaction
	// * \[afrz\] asset-freeze-transaction
	Type *string `json:"type,omitempty"`
}

// TransactionParams defines model for TransactionParams.
type TransactionParams struct {

	// ConsensusVersion indicates the consensus protocol version
	// as of LastRound.
	ConsensusVersion string `json:"consensusVersion"`

	// Fee is the suggested transaction fee
	// Fee is in units of micro-Algos per byte.
	// Fee may fall to zero but transactions must still have a fee of
	// at least MinTxnFee for the current network protocol.
	Fee uint64 `json:"fee"`

	// Genesis ID
	GenesisID string `json:"genesisID"`

	// Genesis hash
	Genesishashb64 string `json:"genesishashb64"`

	// LastRound indicates the last round seen
	LastRound uint64 `json:"lastRound"`

	// The minimum transaction fee (not per byte) required for the
	// txn to validate for the current network protocol.
	MinFee *uint64 `json:"minFee,omitempty"`
}

// Version defines model for Version.
type Version struct {

	// the current algod build version information.
	Build struct {
		Branch      string `json:"branch"`
		BuildNumber uint64 `json:"build_number"`
		Channel     string `json:"channel"`
		CommitHash  string `json:"commit_hash"`
		Major       uint64 `json:"major"`
		Minor       uint64 `json:"minor"`
	} `json:"build"`
	GenesisHashB64 string   `json:"genesis_hash_b64"`
	GenesisId      string   `json:"genesis_id"`
	Versions       []string `json:"versions"`
}

// AccountId defines model for account-id.
type AccountId string

// AssetId defines model for asset-id.
type AssetId uint64

// Gt defines model for gt.
type Gt uint64

// Limit defines model for limit.
type Limit uint64

// Lt defines model for lt.
type Lt uint64

// MaxRound defines model for max-round.
type MaxRound uint64

// MaxTs defines model for max-ts.
type MaxTs uint64

// MinRound defines model for min-round.
type MinRound uint64

// MinTs defines model for min-ts.
type MinTs uint64

// Offset defines model for offset.
type Offset uint64

// Round defines model for round.
type Round uint64

// RoundNumber defines model for round-number.
type RoundNumber uint64

// AccountResponse defines model for AccountResponse.
type AccountResponse Account

// AccountsResponse defines model for AccountsResponse.
type AccountsResponse struct {
	Accounts []Account `json:"accounts"`

	// Round at which the results are valid. This should be the most recent round, so that you can tell how old a cached result is. This field doesn't take into account any max round filter, you'll need to remember that.
	Round uint64 `json:"round"`

	// Total number of results, used to verify whether or not the transaction array is truncated.
	Total uint64 `json:"total"`
}

// AssetBalancesResponse defines model for AssetBalancesResponse.
type AssetBalancesResponse struct {

	// A simplified version of AssetHolding
	Balances MiniAssetHolding `json:"balances"`

	// Round at which the results are valid.
	Round uint64 `json:"round"`

	// The total number of results from thsi query, used to determine whether or not the balances array is truncated.
	Total uint64 `json:"total"`
}

// AssetResponse defines model for AssetResponse.
type AssetResponse Asset

// AssetsResponse defines model for AssetsResponse.
type AssetsResponse struct {
	Assets []Asset `json:"assets"`

	// Round at which the results are valid. This should be the most recent round, so that you can tell how old a cached result is. This field doesn't take into account any max round filter, you'll need to remember that.
	Round uint64 `json:"round"`

	// Total number of results, used to verify whether or not the transaction array is truncated.
	Total uint64 `json:"total"`
}

// BlockResponse defines model for BlockResponse.
type BlockResponse Block

// BlockTimesResponse defines model for BlockTimesResponse.
type BlockTimesResponse struct {
	Rounds *[]struct {
		Round *uint64 `json:"round,omitempty"`

		// Time when block was confirmed.
		Timestamp *uint64 `json:"timestamp,omitempty"`
	} `json:"rounds,omitempty"`
}

// Error defines model for Error.
type Error struct {
	Error *string `json:"error,omitempty"`
}

// TransactionsResponse defines model for TransactionsResponse.
type TransactionsResponse struct {

	// Round at which the results are valid. This should be the most recent round, so that you can tell how old a cached result is. This field doesn't take into account any max round filter, you'll need to remember that.
	Round *uint64 `json:"round,omitempty"`

	// Total number of results, used to verify whether or not the transaction array is truncated.
	Total        *uint64        `json:"total,omitempty"`
	Transactions *[]Transaction `json:"transactions,omitempty"`
}

// LookupAccountByIDParams defines parameters for LookupAccountByID.
type LookupAccountByIDParams struct {

	// Include results for the specified round.
	Round *uint64 `json:"round,omitempty"`
}

// LookupAccountTransactionsParams defines parameters for LookupAccountTransactions.
type LookupAccountTransactionsParams struct {

	// Include results at or after the specified min-round.
	MinRound *uint64 `json:"min-round,omitempty"`

	// Include results at or before the specified max-round.
	MaxRound *uint64 `json:"max-round,omitempty"`

	// Include results at or before the given max timestamp.
	MaxTs *uint64 `json:"max-ts,omitempty"`

	// Include results at or after the given max timestamp.
	MinTs *uint64 `json:"min-ts,omitempty"`
	Asset *uint64 `json:"asset,omitempty"`

	// Used in conjunction with limit to page through results.
	Offset *uint64 `json:"offset,omitempty"`

	// Maximum number of results to return.
	Limit *uint64 `json:"limit,omitempty"`

	// Results should have an amount greater than this value.
	Gt *uint64 `json:"gt,omitempty"`

	// Results should have an amount less than this value.
	Lt *uint64 `json:"lt,omitempty"`
}

// SearchAccountsParams defines parameters for SearchAccounts.
type SearchAccountsParams struct {

	// Include accounts holding the specified asset
	AssetId     *string `json:"asset-id,omitempty"`
	AssetParams *string `json:"assetParams,omitempty"`

	// Maximum number of results to return.
	Limit *uint64 `json:"limit,omitempty"`

	// Used in conjunction with limit to page through results.
	Offset *uint64 `json:"offset,omitempty"`

	// Results should have an amount greater than this value.
	Gt *uint64 `json:"gt,omitempty"`

	// Results should have an amount less than this value.
	Lt *uint64 `json:"lt,omitempty"`
}

// LookupAssetBalancesParams defines parameters for LookupAssetBalances.
type LookupAssetBalancesParams struct {

	// Maximum number of results to return.
	Limit *uint64 `json:"limit,omitempty"`

	// Used in conjunction with limit to page through results.
	Offset *uint64 `json:"offset,omitempty"`

	// Include results for the specified round.
	Round *uint64 `json:"round,omitempty"`

	// Results should have an amount greater than this value.
	Gt *uint64 `json:"gt,omitempty"`

	// Results should have an amount less than this value.
	Lt *uint64 `json:"lt,omitempty"`
}

// LookupAssetTransactionsParams defines parameters for LookupAssetTransactions.
type LookupAssetTransactionsParams struct {

	// Maximum number of results to return.
	Limit *uint64 `json:"limit,omitempty"`

	// Include results at or after the specified min-round.
	MinRound *uint64 `json:"min-round,omitempty"`

	// Include results at or before the specified max-round.
	MaxRound *uint64 `json:"max-round,omitempty"`

	// Include results at or after the given max timestamp.
	MinTs *uint64 `json:"min-ts,omitempty"`

	// Include results at or before the given max timestamp.
	MaxTs *uint64 `json:"max-ts,omitempty"`

	// Used in conjunction with limit to page through results.
	Offset *uint64 `json:"offset,omitempty"`

	// Results should have an amount greater than this value.
	Gt *uint64 `json:"gt,omitempty"`

	// Results should have an amount less than this value.
	Lt *uint64 `json:"lt,omitempty"`
}

// SearchForAssetsParams defines parameters for SearchForAssets.
type SearchForAssetsParams struct {

	// Results should have an amount greater than this value.
	Gt *uint64 `json:"gt,omitempty"`

	// Maximum number of results to return.
	Limit *uint64 `json:"limit,omitempty"`

	// Filter just assets with the given creator address.
	Creator *string `json:"creator,omitempty"`

	// Filter just assets with the given name.
	Name *string `json:"name,omitempty"`

	// Filter just assets with the given unit.
	Unit *string `json:"unit,omitempty"`

	// Used in conjunction with limit to page through results.
	Offset *uint64 `json:"offset,omitempty"`
}

// SearchForTransactionsParams defines parameters for SearchForTransactions.
type SearchForTransactionsParams struct {

	// Specifies a prefix which must be contained in the note field.
	Noteprefix *string `json:"noteprefix,omitempty"`

	// The transaction type, one of:
	// * pay - Payment
	// * keyreg - Key Registration
	// * acfg - Asset Configuration
	// * axfer - Asset Transfer
	// * afrz - Asset Freeze
	Type *string `json:"type,omitempty"`

	// Type of signature used to sign the transaction, must be one of:
	// * standard
	// * multisig
	// * logicsig
	Sigtype *string `json:"sigtype,omitempty"`

	// Lookup the specific transaction by ID.
	Txid *string `json:"txid,omitempty"`

	// Include results for the specified round.
	Round *uint64 `json:"round,omitempty"`

	// Used in conjunction with limit to page through results.
	Offset *uint64 `json:"offset,omitempty"`

	// Include results at or after the specified min-round.
	MinRound *uint64 `json:"min-round,omitempty"`

	// Include results at or before the specified max-round.
	MaxRound *uint64 `json:"max-round,omitempty"`

	// Asset transactions related to a given asset ID.
	AssetId *string `json:"asset-id,omitempty"`

	// Encoding format returned by this endpoint. Default is json.
	Format *string `json:"format,omitempty"`

	// Maximum number of results to return.
	Limit *uint64 `json:"limit,omitempty"`

	// Include results at or before the given max timestamp.
	MaxTs *uint64 `json:"max-ts,omitempty"`

	// Include results at or after the given max timestamp.
	MinTs *uint64 `json:"min-ts,omitempty"`

	// Results should have an amount greater than this value.
	Gt *uint64 `json:"gt,omitempty"`

	// Results should have an amount less than this value.
	Lt *uint64 `json:"lt,omitempty"`
}
