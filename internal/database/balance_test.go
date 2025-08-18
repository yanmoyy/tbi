package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddTokenBalance(t *testing.T) {
	c := getTestClient(t)

	clear := func() {
		err := c.ClearTokenBalances()
		require.NoError(t, err)
	}
	clear()
	defer clear()

	ctx := context.Background()

	err := c.CreateTokenBalance(ctx, "address1", "token1", 10)
	require.NoError(t, err)
	err = c.CreateTokenBalance(ctx, "address2", "token2", 20)
	require.NoError(t, err)

	type input struct {
		address   string
		token     string
		change    int64 // positive or negative
		increment bool
	}

	type expected struct {
		balance int64
	}

	tester := func(t *testing.T, i input, e expected) {
		err = c.UpdateTokenBalance(ctx, i.address, i.token, i.change, i.increment)
		require.NoError(t, err)

		balance, err := c.GetTokenBalance(ctx, i.address, i.token)
		require.NoError(t, err)

		require.Equal(t, e.balance, balance.Amount)
	}

	tester(t, input{
		address:   "address1",
		token:     "token1",
		change:    100,
		increment: true,
	}, expected{
		balance: 110,
	})
	tester(t, input{
		address:   "address2",
		token:     "token2",
		change:    200,
		increment: true,
	}, expected{
		balance: 220,
	})
}

func TestGetTokenBalanceList(t *testing.T) {
	c := getTestClient(t)

	clear := func() {
		err := c.ClearTokenBalances()
		require.NoError(t, err)
	}
	clear()
	defer clear()

	type input struct {
		data []struct {
			address string
			token   string
			amount  int64
		}
	}

	type expected struct {
		data []struct {
			address string
			token   string
			amount  int64
		}
	}

	tester := func(t *testing.T, i input, e expected) {
		ctx := t.Context()
		for _, d := range i.data {
			err := c.CreateTokenBalance(ctx, d.address, d.token, d.amount)
			require.NoError(t, err)
		}
		balances, err := c.GetTokenBalanceList(ctx)
		require.NoError(t, err)
		require.Len(t, balances, len(e.data))
		for i, b := range balances {
			require.Equal(t, e.data[i].address, b.Address)
			require.Equal(t, e.data[i].token, b.TokenPath)
			require.Equal(t, e.data[i].amount, b.Amount)
		}
	}

	tester(t, input{
		data: []struct {
			address string
			token   string
			amount  int64
		}{
			{"c", "token1", 20},
			{"a", "token1", 20},
			{"a", "token2", 10},
			{"b", "token3", 30},
		},
	}, expected{
		data: []struct {
			address string
			token   string
			amount  int64
		}{
			{"a", "token1", 20},
			{"a", "token2", 10},
			{"b", "token3", 30},
			{"c", "token1", 20},
		},
	})

}
