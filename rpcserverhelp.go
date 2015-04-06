// Copyright (c) 2015 Conformal Systems LLC.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"sort"
	"strings"
	"sync"

	"github.com/btcsuite/btcd/btcjson"
)

// helpDescsEnUS defines the English descriptions used for the help strings.
var helpDescsEnUS = map[string]string{
	// DebugLevelCmd help.
	"debuglevel--synopsis": "Dynamically changes the debug logging level.\n" +
		"The levelspec can either a debug level or of the form:\n" +
		"<subsystem>=<level>,<subsystem2>=<level2>,...\n" +
		"The valid debug levels are trace, debug, info, warn, error, and critical.\n" +
		"The valid subsystems are AMGR, ADXR, BCDB, BMGR, BTCD, CHAN, DISC, PEER, RPCS, SCRP, SRVR, and TXMP.\n" +
		"Finally the keyword 'show' will return a list of the available subsystems.",
	"debuglevel-levelspec":   "The debug level(s) to use or the keyword 'show'",
	"debuglevel--condition0": "levelspec!=show",
	"debuglevel--condition1": "levelspec=show",
	"debuglevel--result0":    "The string 'Done.'",
	"debuglevel--result1":    "The list of subsystems",

	// AddNodeCmd help.
	"addnode--synopsis": "Attempts to add or remove a persistent peer.",
	"addnode-addr":      "IP address and port of the peer to operate on",
	"addnode-subcmd":    "'add' to add a persistent peer, 'remove' to remove a persistent peer, or 'onetry' to try a single connection to a peer",

	// NodeCmd help.
	"node--synopsis":     "Attempts to add or remove a peer.",
	"node-subcmd":        "'disconnect' to remove all matching non-persistent peers, 'remove' to remove a persistent peer, or 'connect' to connect to a peer",
	"node-target":        "Either the IP address and port of the peer to operate on, or a valid peer ID.",
	"node-connectsubcmd": "'perm' to make the connected peer a permanent one, 'temp' to try a single connect to a peer",

	// TransactionInput help.
	"transactioninput-txid": "The hash of the input transaction",
	"transactioninput-vout": "The specific output of the input transaction to redeem",

	// CreateRawTransactionCmd help.
	"createrawtransaction--synopsis": "Returns a new transaction spending the provided inputs and sending to the provided addresses.\n" +
		"The transaction inputs are not signed in the created transaction.\n" +
		"The signrawtransaction RPC command provided by wallet must be used to sign the resulting transaction.",
	"createrawtransaction-inputs":         "The inputs to the transaction",
	"createrawtransaction-amounts":        "JSON object with the destination addresses as keys and amounts as values",
	"createrawtransaction-amounts--key":   "address",
	"createrawtransaction-amounts--value": "n.nnn",
	"createrawtransaction-amounts--desc":  "The destination address as the key and the amount in BTC as the value",
	"createrawtransaction--result0":       "Hex-encoded bytes of the serialized transaction",

	// ScriptSig help.
	"scriptsig-asm": "Disassembly of the script",
	"scriptsig-hex": "Hex-encoded bytes of the script",

	// Vin help.
	"vin-coinbase":  "The hex-encoded bytes of the signature script (coinbase txns only)",
	"vin-txid":      "The hash of the origin transaction (non-coinbase txns only)",
	"vin-vout":      "The index of the output being redeemed from the origin transaction (non-coinbase txns only)",
	"vin-scriptSig": "The signature script used to redeem the origin transaction as a JSON object (non-coinbase txns only)",
	"vin-sequence":  "The script sequence number",

	// ScriptPubKeyResult help.
	"scriptpubkeyresult-asm":       "Disassembly of the script",
	"scriptpubkeyresult-hex":       "Hex-encoded bytes of the script",
	"scriptpubkeyresult-reqSigs":   "The number of required signatures",
	"scriptpubkeyresult-type":      "The type of the script (e.g. 'pubkeyhash')",
	"scriptpubkeyresult-addresses": "The bitcoin addresses associated with this script",

	// Vout help.
	"vout-value":        "The amount in BTC",
	"vout-n":            "The index of this transaction output",
	"vout-scriptPubKey": "The public key script used to pay coins as a JSON object",

	// TxRawDecodeResult help.
	"txrawdecoderesult-txid":     "The hash of the transaction",
	"txrawdecoderesult-version":  "The transaction version",
	"txrawdecoderesult-locktime": "The transaction lock time",
	"txrawdecoderesult-vin":      "The transaction inputs as JSON objects",
	"txrawdecoderesult-vout":     "The transaction outputs as JSON objects",

	// DecodeRawTransactionCmd help.
	"decoderawtransaction--synopsis": "Returns a JSON object representing the provided serialized, hex-encoded transaction.",
	"decoderawtransaction-hextx":     "Serialized, hex-encoded transaction",

	// DecodeScriptResult help.
	"decodescriptresult-asm":       "Disassembly of the script",
	"decodescriptresult-reqSigs":   "The number of required signatures",
	"decodescriptresult-type":      "The type of the script (e.g. 'pubkeyhash')",
	"decodescriptresult-addresses": "The bitcoin addresses associated with this script",
	"decodescriptresult-p2sh":      "The script hash for use in pay-to-script-hash transactions",

	// DecodeScriptCmd help.
	"decodescript--synopsis": "Returns a JSON object with information about the provided hex-encoded script.",
	"decodescript-hexscript": "Hex-encoded script",

	// GetAddedNodeInfoResultAddr help.
	"getaddednodeinforesultaddr-address":   "The ip address for this DNS entry",
	"getaddednodeinforesultaddr-connected": "The connection 'direction' (inbound/outbound/false)",

	// GetAddedNodeInfoResult help.
	"getaddednodeinforesult-addednode": "The ip address or domain of the added peer",
	"getaddednodeinforesult-connected": "Whether or not the peer is currently connected",
	"getaddednodeinforesult-addresses": "DNS lookup and connection information about the peer",

	// GetAddedNodeInfo help.
	"getaddednodeinfo--synopsis":   "Returns information about manually added (persistent) peers.",
	"getaddednodeinfo-dns":         "Specifies whether the returned data is a JSON object including DNS and connection information, or just a list of added peers",
	"getaddednodeinfo-node":        "Only return information about this specific peer instead of all added peers",
	"getaddednodeinfo--condition0": "dns=false",
	"getaddednodeinfo--condition1": "dns=true",
	"getaddednodeinfo--result0":    "List of added peers",

	// GetBestBlockResult help.
	"getbestblockresult-hash":   "Hex-encoded bytes of the best block hash",
	"getbestblockresult-height": "Height of the best block",

	// GetBestBlockCmd help.
	"getbestblock--synopsis": "Get block height and hash of best block in the main chain.",
	"getbestblock--result0":  "Get block height and hash of best block in the main chain.",

	// GetBestBlockHashCmd help.
	"getbestblockhash--synopsis": "Returns the hash of the of the best (most recent) block in the longest block chain.",
	"getbestblockhash--result0":  "The hex-encoded block hash",

	// GetBlockCmd help.
	"getblock--synopsis":   "Returns information about a block given its hash.",
	"getblock-hash":        "The hash of the block",
	"getblock-verbose":     "Specifies the block is returned as a JSON object instead of hex-encoded string",
	"getblock-verbosetx":   "Specifies that each transaction is returned as a JSON object and only applies if the verbose flag is true (btcd extension)",
	"getblock--condition0": "verbose=false",
	"getblock--condition1": "verbose=true",
	"getblock--result0":    "Hex-encoded bytes of the serialized block",

	// TxRawResult help.
	"txrawresult-hex":           "Hex-encoded transaction",
	"txrawresult-txid":          "The hash of the transaction",
	"txrawresult-version":       "The transaction version",
	"txrawresult-locktime":      "The transaction lock time",
	"txrawresult-vin":           "The transaction inputs as JSON objects",
	"txrawresult-vout":          "The transaction outputs as JSON objects",
	"txrawresult-blockhash":     "Hash of the block the transaction is part of",
	"txrawresult-confirmations": "Number of confirmations of the block",
	"txrawresult-time":          "Transaction time in seconds since 1 Jan 1970 GMT",
	"txrawresult-blocktime":     "Block time in seconds since the 1 Jan 1970 GMT",

	// GetBlockVerboseResult help.
	"getblockverboseresult-hash":              "The hash of the block (same as provided)",
	"getblockverboseresult-confirmations":     "The number of confirmations",
	"getblockverboseresult-size":              "The size of the block",
	"getblockverboseresult-height":            "The height of the block in the block chain",
	"getblockverboseresult-version":           "The block version",
	"getblockverboseresult-merkleroot":        "Root hash of the merkle tree",
	"getblockverboseresult-tx":                "The transaction hashes (only when verbosetx=false)",
	"getblockverboseresult-rawtx":             "The transactions as JSON objects (only when verbosetx=true)",
	"getblockverboseresult-time":              "The block time in seconds since 1 Jan 1970 GMT",
	"getblockverboseresult-nonce":             "The block nonce",
	"getblockverboseresult-bits":              "The bits which represent the block difficulty",
	"getblockverboseresult-difficulty":        "The proof-of-work difficulty as a multiple of the minimum difficulty",
	"getblockverboseresult-previousblockhash": "The hash of the previous block",
	"getblockverboseresult-nextblockhash":     "The hash of the next block",

	// GetBlockCountCmd help.
	"getblockcount--synopsis": "Returns the number of blocks in the longest block chain.",
	"getblockcount--result0":  "The current block count",

	// GetBlockHashCmd help.
	"getblockhash--synopsis": "Returns hash of the block in best block chain at the given height.",
	"getblockhash-index":     "The block height",
	"getblockhash--result0":  "The block hash",

	// TemplateRequest help.
	"templaterequest-mode":         "This is 'template', 'proposal', or omitted",
	"templaterequest-capabilities": "List of capabilities",
	"templaterequest-longpollid":   "The long poll ID of a job to monitor for expiration; required and valid only for long poll requests ",
	"templaterequest-sigoplimit":   "Number of signature operations allowed in blocks (this parameter is ignored)",
	"templaterequest-sizelimit":    "Number of bytes allowed in blocks (this parameter is ignored)",
	"templaterequest-maxversion":   "Highest supported block version number (this parameter is ignored)",
	"templaterequest-target":       "The desired target for the block template (this parameter is ignored)",
	"templaterequest-data":         "Hex-encoded block data (only for mode=proposal)",
	"templaterequest-workid":       "The server provided workid if provided in block template (not applicable)",

	// GetBlockTemplateResultTx help.
	"getblocktemplateresulttx-data":    "Hex-encoded transaction data (byte-for-byte)",
	"getblocktemplateresulttx-hash":    "Hex-encoded transaction hash (little endian if treated as a 256-bit number)",
	"getblocktemplateresulttx-depends": "Other transactions before this one (by 1-based index in the 'transactions'  list) that must be present in the final block if this one is",
	"getblocktemplateresulttx-fee":     "Difference in value between transaction inputs and outputs (in Satoshi)",
	"getblocktemplateresulttx-sigops":  "Total number of signature operations as counted for purposes of block limits",

	// GetBlockTemplateResultAux help.
	"getblocktemplateresultaux-flags": "Hex-encoded byte-for-byte data to include in the coinbase signature script",

	// GetBlockTemplateResult help.
	"getblocktemplateresult-bits":              "Hex-encoded compressed difficulty",
	"getblocktemplateresult-curtime":           "Current time as seen by the server (recommended for block time); must fall within mintime/maxtime rules",
	"getblocktemplateresult-height":            "Height of the block to be solved",
	"getblocktemplateresult-previousblockhash": "Hex-encoded big-endian hash of the previous block",
	"getblocktemplateresult-sigoplimit":        "Number of sigops allowed in blocks ",
	"getblocktemplateresult-sizelimit":         "Number of bytes allowed in blocks",
	"getblocktemplateresult-transactions":      "Array of transactions as JSON objects",
	"getblocktemplateresult-version":           "The block version",
	"getblocktemplateresult-coinbaseaux":       "Data that should be included in the coinbase signature script",
	"getblocktemplateresult-coinbasetxn":       "Information about the coinbase transaction",
	"getblocktemplateresult-coinbasevalue":     "Total amount available for the coinbase in Satoshi",
	"getblocktemplateresult-workid":            "This value must be returned with result if provided (not provided)",
	"getblocktemplateresult-longpollid":        "Identifier for long poll request which allows monitoring for expiration",
	"getblocktemplateresult-longpolluri":       "An alternate URI to use for long poll requests if provided (not provided)",
	"getblocktemplateresult-submitold":         "Not applicable",
	"getblocktemplateresult-target":            "Hex-encoded big-endian number which valid results must be less than",
	"getblocktemplateresult-expires":           "Maximum number of seconds (starting from when the server sent the response) this work is valid for",
	"getblocktemplateresult-maxtime":           "Maximum allowed time",
	"getblocktemplateresult-mintime":           "Minimum allowed time",
	"getblocktemplateresult-mutable":           "List of mutations the server explicitly allows",
	"getblocktemplateresult-noncerange":        "Two concatenated hex-encoded big-endian 32-bit integers which represent the valid ranges of nonces the miner may scan",
	"getblocktemplateresult-capabilities":      "List of server capabilities including 'proposal' to indicate support for block proposals",
	"getblocktemplateresult-reject-reason":     "Reason the proposal was invalid as-is (only applies to proposal responses)",

	// GetBlockTemplateCmd help.
	"getblocktemplate--synopsis": "Returns a JSON object with information necessary to construct a block to mine or accepts a proposal to validate.\n" +
		"See BIP0022 and BIP0023 for the full specification.",
	"getblocktemplate-request":     "Request object which controls the mode and several parameters",
	"getblocktemplate--condition0": "mode=template",
	"getblocktemplate--condition1": "mode=proposal, rejected",
	"getblocktemplate--condition2": "mode=proposal, accepted",
	"getblocktemplate--result1":    "An error string which represents why the proposal was rejected or nothing if accepted",

	// GetConnectionCountCmd help.
	"getconnectioncount--synopsis": "Returns the number of active connections to other peers.",
	"getconnectioncount--result0":  "The number of connections",

	// GetCurrentNetCmd help.
	"getcurrentnet--synopsis": "Get bitcoin network the server is running on.",
	"getcurrentnet--result0":  "The network identifer",

	// GetDifficultyCmd help.
	"getdifficulty--synopsis": "Returns the proof-of-work difficulty as a multiple of the minimum difficulty.",
	"getdifficulty--result0":  "The difficulty",

	// GetGenerateCmd help.
	"getgenerate--synopsis": "Returns if the server is set to generate coins (mine) or not.",
	"getgenerate--result0":  "True if mining, false if not",

	// GetHashesPerSecCmd help.
	"gethashespersec--synopsis": "Returns a recent hashes per second performance measurement while generating coins (mining).",
	"gethashespersec--result0":  "The number of hashes per second",

	// InfoChainResult help.
	"infochainresult-version":         "The version of the server",
	"infochainresult-protocolversion": "The latest supported protocol version",
	"infochainresult-blocks":          "The number of blocks processed",
	"infochainresult-timeoffset":      "The time offset",
	"infochainresult-connections":     "The number of connected peers",
	"infochainresult-proxy":           "The proxy used by the server",
	"infochainresult-difficulty":      "The current target difficulty",
	"infochainresult-testnet":         "Whether or not server is using testnet",
	"infochainresult-relayfee":        "The minimum relay fee for non-free transactions in BTC/KB",
	"infochainresult-errors":          "Any current errors",

	// InfoWalletResult help.
	"infowalletresult-version":         "The version of the server",
	"infowalletresult-protocolversion": "The latest supported protocol version",
	"infowalletresult-walletversion":   "The version of the wallet server",
	"infowalletresult-balance":         "The total bitcoin balance of the wallet",
	"infowalletresult-blocks":          "The number of blocks processed",
	"infowalletresult-timeoffset":      "The time offset",
	"infowalletresult-connections":     "The number of connected peers",
	"infowalletresult-proxy":           "The proxy used by the server",
	"infowalletresult-difficulty":      "The current target difficulty",
	"infowalletresult-testnet":         "Whether or not server is using testnet",
	"infowalletresult-keypoololdest":   "Seconds since 1 Jan 1970 GMT of the oldest pre-generated key in the key pool",
	"infowalletresult-keypoolsize":     "The number of new keys that are pre-generated",
	"infowalletresult-unlocked_until":  "The timestamp in seconds since 1 Jan 1970 GMT that the wallet is unlocked for transfers, or 0 if the wallet is locked",
	"infowalletresult-paytxfee":        "The transaction fee set in BTC/KB",
	"infowalletresult-relayfee":        "The minimum relay fee for non-free transactions in BTC/KB",
	"infowalletresult-errors":          "Any current errors",

	// GetInfoCmd help.
	"getinfo--synopsis": "Returns a JSON object containing various state info.",

	// GetMiningInfoResult help.
	"getmininginforesult-blocks":           "Height of the latest best block",
	"getmininginforesult-currentblocksize": "Size of the latest best block",
	"getmininginforesult-currentblocktx":   "Number of transactions in the latest best block",
	"getmininginforesult-difficulty":       "Current target difficulty",
	"getmininginforesult-errors":           "Any current errors",
	"getmininginforesult-generate":         "Whether or not server is set to generate coins",
	"getmininginforesult-genproclimit":     "Number of processors to use for coin generation (-1 when disabled)",
	"getmininginforesult-hashespersec":     "Recent hashes per second performance measurement while generating coins",
	"getmininginforesult-networkhashps":    "Estimated network hashes per second for the most recent blocks",
	"getmininginforesult-pooledtx":         "Number of transactions in the memory pool",
	"getmininginforesult-testnet":          "Whether or not server is using testnet",

	// GetMiningInfoCmd help.
	"getmininginfo--synopsis": "Returns a JSON object containing mining-related information.",

	// GetNetworkHashPSCmd help.
	"getnetworkhashps--synopsis": "Returns the estimated network hashes per second for the block heights provided by the parameters.",
	"getnetworkhashps-blocks":    "The number of blocks, or -1 for blocks since last difficulty change",
	"getnetworkhashps-height":    "Perform estimate ending with this height or -1 for current best chain block height",
	"getnetworkhashps--result0":  "Estimated hashes per second",

	// GetNetTotalsCmd help.
	"getnettotals--synopsis": "Returns a JSON object containing network traffic statistics.",

	// GetNetTotalsResult help.
	"getnettotalsresult-totalbytesrecv": "Total bytes received",
	"getnettotalsresult-totalbytessent": "Total bytes sent",
	"getnettotalsresult-timemillis":     "Number of milliseconds since 1 Jan 1970 GMT",

	// GetPeerInfoResult help.
	"getpeerinforesult-id":             "A unique node ID",
	"getpeerinforesult-addr":           "The ip address and port of the peer",
	"getpeerinforesult-addrlocal":      "Local address",
	"getpeerinforesult-services":       "Services bitmask which represents the services supported by the peer",
	"getpeerinforesult-lastsend":       "Time the last message was received in seconds since 1 Jan 1970 GMT",
	"getpeerinforesult-lastrecv":       "Time the last message was sent in seconds since 1 Jan 1970 GMT",
	"getpeerinforesult-bytessent":      "Total bytes sent",
	"getpeerinforesult-bytesrecv":      "Total bytes received",
	"getpeerinforesult-conntime":       "Time the connection was made in seconds since 1 Jan 1970 GMT",
	"getpeerinforesult-timeoffset":     "The time offset of the peer",
	"getpeerinforesult-pingtime":       "Number of microseconds the last ping took",
	"getpeerinforesult-pingwait":       "Number of microseconds a queued ping has been waiting for a response",
	"getpeerinforesult-version":        "The protocol version of the peer",
	"getpeerinforesult-subver":         "The user agent of the peer",
	"getpeerinforesult-inbound":        "Whether or not the peer is an inbound connection",
	"getpeerinforesult-startingheight": "The latest block height the peer knew about when the connection was established",
	"getpeerinforesult-currentheight":  "The current height of the peer",
	"getpeerinforesult-banscore":       "The ban score",
	"getpeerinforesult-syncnode":       "Whether or not the peer is the sync peer",

	// GetPeerInfoCmd help.
	"getpeerinfo--synopsis": "Returns data about each connected network peer as an array of json objects.",

	// GetRawMempoolVerboseResult help.
	"getrawmempoolverboseresult-size":             "Transaction size in bytes",
	"getrawmempoolverboseresult-fee":              "Transaction fee in bitcoins",
	"getrawmempoolverboseresult-time":             "Local time transaction entered pool in seconds since 1 Jan 1970 GMT",
	"getrawmempoolverboseresult-height":           "Block height when transaction entered the pool",
	"getrawmempoolverboseresult-startingpriority": "Priority when transaction entered the pool",
	"getrawmempoolverboseresult-currentpriority":  "Current priority",
	"getrawmempoolverboseresult-depends":          "Unconfirmed transactions used as inputs for this transaction",

	// GetRawMempoolCmd help.
	"getrawmempool--synopsis":   "Returns information about all of the transactions currently in the memory pool.",
	"getrawmempool-verbose":     "Returns JSON object when true or an array of transaction hashes when false",
	"getrawmempool--condition0": "verbose=false",
	"getrawmempool--condition1": "verbose=true",
	"getrawmempool--result0":    "Array of transaction hashes",

	// GetRawTransactionCmd help.
	"getrawtransaction--synopsis":   "Returns information about a transaction given its hash.",
	"getrawtransaction-txid":        "The hash of the transaction",
	"getrawtransaction-verbose":     "Specifies the transaction is returned as a JSON object instead of a hex-encoded string",
	"getrawtransaction--condition0": "verbose=false",
	"getrawtransaction--condition1": "verbose=true",
	"getrawtransaction--result0":    "Hex-encoded bytes of the serialized transaction",

	// GetTxOutResult help.
	"gettxoutresult-bestblock":     "The block hash that contains the transaction output",
	"gettxoutresult-confirmations": "The number of confirmations",
	"gettxoutresult-value":         "The transaction amount in BTC",
	"gettxoutresult-scriptPubKey":  "The public key script used to pay coins as a JSON object",
	"gettxoutresult-version":       "The transaction version",
	"gettxoutresult-coinbase":      "Whether or not the transaction is a coinbase",

	// GetTxOutCmd help.
	"gettxout--synopsis":      "Returns information about an unspent transaction output..",
	"gettxout-txid":           "The hash of the transaction",
	"gettxout-vout":           "The index of the output",
	"gettxout-includemempool": "Include the mempool when true",

	// GetWorkResult help.
	"getworkresult-data":     "Hex-encoded block data",
	"getworkresult-hash1":    "(DEPRECATED) Hex-encoded formatted hash buffer",
	"getworkresult-midstate": "(DEPRECATED) Hex-encoded precomputed hash state after hashing first half of the data",
	"getworkresult-target":   "Hex-encoded little-endian hash target",

	// GetWorkCmd help.
	"getwork--synopsis":   "(DEPRECATED - Use getblocktemplate instead) Returns formatted hash data to work on or checks and submits solved data.",
	"getwork-data":        "Hex-encoded data to check",
	"getwork--condition0": "no data provided",
	"getwork--condition1": "data provided",
	"getwork--result1":    "Whether or not the solved data is valid and was added to the chain",

	// HelpCmd help.
	"help--synopsis":   "Returns a list of all commands or help for a specified command.",
	"help-command":     "The command to retrieve help for",
	"help--condition0": "no command provided",
	"help--condition1": "command specified",
	"help--result0":    "List of commands",
	"help--result1":    "Help for specified command",

	// PingCmd help.
	"ping--synopsis": "Queues a ping to be sent to each connected peer.\n" +
		"Ping times are provided by getpeerinfo via the pingtime and pingwait fields.",

	// SearchRawTransactionsCmd help.
	"searchrawtransactions--synopsis": "Returns raw data for transactions involving the passed address.\n" +
		"Returned transactions are pulled from both the database, and transactions currently in the mempool.\n" +
		"Transactions pulled from the mempool will have the 'confirmations' field set to 0.\n" +
		"Usage of this RPC requires the optional --addrindex flag to be activated, otherwise all responses will simply return with an error stating the address index has not yet been built.\n" +
		"Similarly, until the address index has caught up with the current best height, all requests will return an error response in order to avoid serving stale data.",
	"searchrawtransactions-address":     "The Bitcoin address to search for",
	"searchrawtransactions-verbose":     "Specifies the transaction is returned as a JSON object instead of hex-encoded string",
	"searchrawtransactions-skip":        "The number of leading transactions to leave out of the final response",
	"searchrawtransactions-count":       "The maximum number of transactions to return",
	"searchrawtransactions--condition0": "verbose=0",
	"searchrawtransactions--condition1": "verbose=1",
	"searchrawtransactions--result0":    "Hex-encoded serialized transaction",

	// SendRawTransactionCmd help.
	"sendrawtransaction--synopsis":     "Submits the serialized, hex-encoded transaction to the local peer and relays it to the network.",
	"sendrawtransaction-hextx":         "Serialized, hex-encoded signed transaction",
	"sendrawtransaction-allowhighfees": "Whether or not to allow insanely high fees (btcd does not yet implement this parameter, so it has no effect)",
	"sendrawtransaction--result0":      "The hash of the transaction",

	// SetGenerateCmd help.
	"setgenerate--synopsis":    "Set the server to generate coins (mine) or not.",
	"setgenerate-generate":     "Use true to enable generation, false to disable it",
	"setgenerate-genproclimit": "The number of processors (cores) to limit generation to or -1 for default",

	// StopCmd help.
	"stop--synopsis": "Shutdown btcd.",
	"stop--result0":  "The string 'btcd stopping.'",

	// SubmitBlockOptions help.
	"submitblockoptions-workid": "This parameter is currently ignored",

	// SubmitBlockCmd help.
	"submitblock--synopsis":   "Attempts to submit a new serialized, hex-encoded block to the network.",
	"submitblock-hexblock":    "Serialized, hex-encoded block",
	"submitblock-options":     "This parameter is currently ignored",
	"submitblock--condition0": "Block successfully submitted",
	"submitblock--condition1": "Block rejected",
	"submitblock--result1":    "The reason the block was rejected",

	// ValidateAddressResult help.
	"validateaddresschainresult-isvalid": "Whether or not the address is valid",
	"validateaddresschainresult-address": "The bitcoin address (only when isvalid is true)",

	// ValidateAddressCmd help.
	"validateaddress--synopsis": "Verify an address is valid.",
	"validateaddress-address":   "Bitcoin address to validate",

	// VerifyChainCmd help.
	"verifychain--synopsis": "Verifies the block chain database.\n" +
		"The actual checks performed by the checklevel parameter are implementation specific.\n" +
		"For btcd this is:\n" +
		"checklevel=0 - Look up each block and ensure it can be loaded from the database.\n" +
		"checklevel=1 - Perform basic context-free sanity checks on each block.",
	"verifychain-checklevel": "How thorough the block verification is",
	"verifychain-checkdepth": "The number of blocks to check",
	"verifychain--result0":   "Whether or not the chain verified",

	// VerifyMessageCmd help.
	"verifymessage--synopsis": "Verify a signed message.",
	"verifymessage-address":   "The bitcoin address to use for the signature",
	"verifymessage-signature": "The base-64 encoded signature provided by the signer",
	"verifymessage-message":   "The signed message",
	"verifymessage--result0":  "Whether or not the signature verified",

	// -------- Websocket-specific help --------

	// NotifyBlocksCmd help.
	"notifyblocks--synopsis": "Request notifications for whenever a block is connected or disconnected from the main (best) chain.",

	// NotifyNewTransactionsCmd help.
	"notifynewtransactions--synopsis": "Send either a txaccepted or a txacceptedverbose notification when a new transaction is accepted into the mempool.",
	"notifynewtransactions-verbose":   "Specifies which type of notification to receive. If verbose is true, then the caller receives txacceptedverbose, otherwise the caller receives txaccepted",

	// NotifyReceivedCmd help.
	"notifyreceived--synopsis": "Send a recvtx notification when a transaction added to mempool or appears in a newly-attached block contains a txout pkScript sending to any of the passed addresses.\n" +
		"Matching outpoints are automatically registered for redeemingtx notifications.",
	"notifyreceived-addresses": "List of address to receive notifications about",

	// OutPoint help.
	"outpoint-hash":  "The hex-encoded bytes of the outpoint hash",
	"outpoint-index": "The index of the outpoint",

	// NotifySpentCmd help.
	"notifyspent--synopsis": "Send a redeemingtx notification when a transaction spending an outpoint appears in mempool (if relayed to this btcd instance) and when such a transaction first appears in a newly-attached block.",
	"notifyspent-outpoints": "List of transaction outpoints to monitor.",

	// Rescan help.
	"rescan--synopsis": "Rescan block chain for transactions to addresses.\n" +
		"When the endblock parameter is omitted, the rescan continues through the best block in the main chain.\n" +
		"Rescan results are sent as recvtx and redeemingtx notifications.\n" +
		"This call returns once the rescan completes.",
	"rescan-beginblock": "Hash of the first block to begin rescanning",
	"rescan-addresses":  "List of addresses to include in the rescan",
	"rescan-outpoints":  "List of transaction outpoints to include in the rescan",
	"rescan-endblock":   "Hash of final block to rescan",
}

// rpcResultTypes specifies the result types that each RPC command can return.
// This information is used to generate the help.  Each result type must be a
// pointer to the type (or nil to indicate no return value).
var rpcResultTypes = map[string][]interface{}{
	"addnode":               nil,
	"createrawtransaction":  []interface{}{(*string)(nil)},
	"debuglevel":            []interface{}{(*string)(nil), (*string)(nil)},
	"decoderawtransaction":  []interface{}{(*btcjson.TxRawDecodeResult)(nil)},
	"decodescript":          []interface{}{(*btcjson.DecodeScriptResult)(nil)},
	"getaddednodeinfo":      []interface{}{(*[]string)(nil), (*[]btcjson.GetAddedNodeInfoResult)(nil)},
	"getbestblock":          []interface{}{(*btcjson.GetBestBlockResult)(nil)},
	"getbestblockhash":      []interface{}{(*string)(nil)},
	"getblock":              []interface{}{(*string)(nil), (*btcjson.GetBlockVerboseResult)(nil)},
	"getblockcount":         []interface{}{(*int64)(nil)},
	"getblockhash":          []interface{}{(*string)(nil)},
	"getblocktemplate":      []interface{}{(*btcjson.GetBlockTemplateResult)(nil), (*string)(nil), nil},
	"getconnectioncount":    []interface{}{(*int32)(nil)},
	"getcurrentnet":         []interface{}{(*uint32)(nil)},
	"getdifficulty":         []interface{}{(*float64)(nil)},
	"getgenerate":           []interface{}{(*bool)(nil)},
	"gethashespersec":       []interface{}{(*float64)(nil)},
	"getinfo":               []interface{}{(*btcjson.InfoChainResult)(nil)},
	"getmininginfo":         []interface{}{(*btcjson.GetMiningInfoResult)(nil)},
	"getnettotals":          []interface{}{(*btcjson.GetNetTotalsResult)(nil)},
	"getnetworkhashps":      []interface{}{(*int64)(nil)},
	"getpeerinfo":           []interface{}{(*[]btcjson.GetPeerInfoResult)(nil)},
	"getrawmempool":         []interface{}{(*[]string)(nil), (*btcjson.GetRawMempoolVerboseResult)(nil)},
	"getrawtransaction":     []interface{}{(*string)(nil), (*btcjson.TxRawResult)(nil)},
	"gettxout":              []interface{}{(*btcjson.GetTxOutResult)(nil)},
	"getwork":               []interface{}{(*btcjson.GetWorkResult)(nil), (*bool)(nil)},
	"node":                  nil,
	"help":                  []interface{}{(*string)(nil), (*string)(nil)},
	"ping":                  nil,
	"searchrawtransactions": []interface{}{(*string)(nil), (*[]btcjson.TxRawResult)(nil)},
	"sendrawtransaction":    []interface{}{(*string)(nil)},
	"setgenerate":           nil,
	"stop":                  []interface{}{(*string)(nil)},
	"submitblock":           []interface{}{nil, (*string)(nil)},
	"validateaddress":       []interface{}{(*btcjson.ValidateAddressChainResult)(nil)},
	"verifychain":           []interface{}{(*bool)(nil)},
	"verifymessage":         []interface{}{(*bool)(nil)},

	// Websocket commands.
	"notifyblocks":          nil,
	"notifynewtransactions": nil,
	"notifyreceived":        nil,
	"notifyspent":           nil,
	"rescan":                nil,
}

// helpCacher provides a concurrent safe type that provides help and usage for
// the RPC server commands and caches the results for future calls.
type helpCacher struct {
	sync.Mutex
	usage      string
	methodHelp map[string]string
}

// rpcMethodHelp returns an RPC help string for the provided method.
//
// This function is safe for concurrent access.
func (c *helpCacher) rpcMethodHelp(method string) (string, error) {
	c.Lock()
	defer c.Unlock()

	// Return the cached method help if it exists.
	if help, exists := c.methodHelp[method]; exists {
		return help, nil
	}

	// Look up the result types for the method.
	resultTypes, ok := rpcResultTypes[method]
	if !ok {
		return "", errors.New("no result types specified for method " +
			method)
	}

	// Generate, cache, and return the help.
	help, err := btcjson.GenerateHelp(method, helpDescsEnUS, resultTypes...)
	if err != nil {
		return "", err
	}
	c.methodHelp[method] = help
	return help, nil
}

// rpcUsage returns one-line usage for all support RPC commands.
//
// This function is safe for concurrent access.
func (c *helpCacher) rpcUsage(includeWebsockets bool) (string, error) {
	c.Lock()
	defer c.Unlock()

	// Return the cached usage if it is available.
	if c.usage != "" {
		return c.usage, nil
	}

	// Generate a list of one-line usage for every command.
	usageTexts := make([]string, 0, len(rpcHandlers))
	for k := range rpcHandlers {
		usage, err := btcjson.MethodUsageText(k)
		if err != nil {
			return "", err
		}
		usageTexts = append(usageTexts, usage)
	}

	// Include websockets commands if requested.
	if includeWebsockets {
		for k := range wsHandlers {
			usage, err := btcjson.MethodUsageText(k)
			if err != nil {
				return "", err
			}
			usageTexts = append(usageTexts, usage)
		}
	}

	sort.Sort(sort.StringSlice(usageTexts))
	c.usage = strings.Join(usageTexts, "\n")
	return c.usage, nil
}

// newHelpCacher returns a new instance of a help cacher which provides help and
// usage for the RPC server commands and caches the results for future calls.
func newHelpCacher() *helpCacher {
	return &helpCacher{
		methodHelp: make(map[string]string),
	}
}