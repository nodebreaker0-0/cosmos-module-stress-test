package tx_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/test-go/testify/require"

	"github.com/b-harvest/liquidity-stress-test/client"
	"github.com/b-harvest/liquidity-stress-test/tx"
	"github.com/b-harvest/liquidity-stress-test/wallet"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

var (
	c *client.Client

	rpcAddress  = "http://localhost:26657"
	grpcAddress = "localhost:9090"
)

func TestMain(m *testing.M) {
	c = client.NewClient(rpcAddress, grpcAddress)

	os.Exit(m.Run())
}

func TestFindAllPairs(t *testing.T) {
	pairs := []struct {
		pairs []string
	}{
		{
			[]string{
				"uakt",
				"uatom",
				"ubtsg",
				"udvpn",
				"ugcyb",
				"uiris",
				"uluna",
				"ungm",
				"uxprt",
				"uxrn",
				"xrun",
			},
		},
	}

	for _, p := range pairs {
		for i := 0; i < len(p.pairs)-1; i++ {
			for j := i + 1; j < len(p.pairs); j++ {
				t.Log(p.pairs[i], p.pairs[j])
			}
		}
	}
}

func TestDepositWithinBatch(t *testing.T) {
	mnemonic := "guard cream sadness conduct invite crumble clock pudding hole grit liar hotel maid produce squeeze return argue turtle know drive eight casino maze host"

	accAddr, privKey, err := wallet.RecoverAccountFromMnemonic(mnemonic, "")
	require.NoError(t, err)

	chainID, err := c.RPC.GetNetworkChainID(context.Background())
	require.NoError(t, err)

	testCases := []struct {
		name         string
		accAddr      string
		privKey      *secp256k1.PrivKey
		poolId       uint64
		depositCoins sdktypes.Coins
	}{
		{
			"uakt/uatom",
			accAddr,
			privKey,
			uint64(1),
			sdktypes.NewCoins(sdktypes.NewCoin("uakt", sdktypes.NewInt(33333)), sdktypes.NewCoin("uatom", sdktypes.NewInt(33333))),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			msg, err := tx.MsgDeposit(tc.accAddr, tc.poolId, tc.depositCoins)
			require.NoError(t, err)

			msgs := []sdktypes.Msg{msg}

			tx := tx.NewTransaction(c, chainID, tx.DefaultGasLimit, tx.DefaultFees, tx.DefaultMemo)

			resp, err := tx.SignAndBroadcast(context.Background(), accAddr, privKey, msgs...)
			require.NoError(t, err)

			fmt.Println("Code: ", resp.TxResponse.Code)
			fmt.Println("Height: ", resp.TxResponse.Height)
			fmt.Println("TxHash: ", resp.TxResponse.TxHash)
		})
	}
}

func TestWithdrawWithinBatch(t *testing.T) {
	mnemonic := "guard cream sadness conduct invite crumble clock pudding hole grit liar hotel maid produce squeeze return argue turtle know drive eight casino maze host"

	accAddr, privKey, err := wallet.RecoverAccountFromMnemonic(mnemonic, "")
	require.NoError(t, err)

	chainID, err := c.RPC.GetNetworkChainID(context.Background())
	require.NoError(t, err)

	testCases := []struct {
		name     string
		accAddr  string
		privKey  *secp256k1.PrivKey
		poolId   uint64
		poolCoin sdktypes.Coin
	}{
		{
			"uakt/uatom",
			accAddr,
			privKey,
			uint64(1),
			sdktypes.NewCoin("pool94720F40B38D6DD93DCE184D264D4BE089EDF124A9C0658CDBED6CA18CF27752", sdktypes.NewInt(5)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			msg, err := tx.MsgWithdraw(tc.accAddr, tc.poolId, tc.poolCoin)
			require.NoError(t, err)

			msgs := []sdktypes.Msg{msg}

			tx := tx.NewTransaction(c, chainID, tx.DefaultGasLimit, tx.DefaultFees, tx.DefaultMemo)

			resp, err := tx.SignAndBroadcast(context.Background(), accAddr, privKey, msgs...)
			require.NoError(t, err)

			fmt.Println("Code: ", resp.TxResponse.Code)
			fmt.Println("Height: ", resp.TxResponse.Height)
			fmt.Println("TxHash: ", resp.TxResponse.TxHash)
		})
	}
}
