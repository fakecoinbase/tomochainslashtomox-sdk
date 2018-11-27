package config

import (
	"github.com/ethereum/go-ethereum/crypto"
)

type EthereumConfig struct {
	NetworkID       string `valid:"required,int" toml:"network_id"`
	MasterPublicKey string `valid:"required" toml:"master_public_key"`
	// Minimum value of transaction accepted by Bifrost in ETH.
	// Everything below will be ignored.
	MinimumValueEth string `valid:"required" toml:"minimum_value_eth"`
	// TokenPrice is a price of one token in ETH
	TokenPrice string `valid:"required" toml:"token_price"`
	// Host only
	RpcServer string `valid:"required" toml:"rpc_server"`
}

type TomochainConfig struct {
	// TokenAssetCode is asset code of token that will be purchased using ETH.
	TokenAssetCode string `valid:"required" toml:"token_asset_code"`
	// IssuerPublicKey is public key of the assets issuer.
	IssuerPublicKey string `valid:"required,tomochain_accountid" toml:"issuer_public_key"`
	// DistributionPublicKey is public key of the distribution account.
	// Distribution account can be the same account as issuer account however it's recommended
	// to use a separate account.
	// Distribution account is also used to fund new accounts, this is via smart contract.
	DistributionPublicKey string `valid:"required,tomochain_accountid" toml:"distribution_public_key"`
	// SignerPrivateKey is:
	// * Distribution's secret key if only one instance of Bifrost is deployed.
	// * Channel's secret key of Distribution account if more than one instance of Bifrost is deployed.
	// https://www.tomochain.org/developers/guides/channels.html
	// Signer's sequence number will be consumed in transaction's sequence number.
	SignerPrivateKey string `valid:"required,tomochain_seed" toml:"signer_secret_key"`
	// StartingBalance is the starting amount of XLM for newly created accounts.
	// Default value is 41. Increase it if you need Data records / other custom entities on new account.
	StartingBalance string `valid:"optional,tomochain_amount" toml:"starting_balance"`
	// LockUnixTimestamp defines unix timestamp when user account will be unlocked.
	LockUnixTimestamp uint64 `valid:"optional" toml:"lock_unix_timestamp"`
}

type Config struct {
	Ethereum  *EthereumConfig  `valid:"optional" json:"ethereum"`
	Tomochain *TomochainConfig `valid:"optional" json:"tomochain"`
}

func (c *Config) SignerPublicKey() string {
	if c.Tomochain == nil {
		return ""
	}
	// from private key to sign smart contract
	privkey, err := crypto.LoadECDSA(c.Tomochain.SignerPrivateKey)
	if err != nil {
		return ""
	}
	address := crypto.PubkeyToAddress(privkey.PublicKey)
	return address.String()
}
