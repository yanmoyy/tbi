package synchronizer

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yanmoyy/tbi/internal/indexer"
)

func TestProcessEvent(t *testing.T) {
	type input struct {
		event indexer.GnoEvent
	}
	type expected struct {
		valid bool
	}

	tester := func(t *testing.T, input input, expected expected) {
		_, err := processEvent(input.event)
		if expected.valid {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}

	t.Run("Mint", func(t *testing.T) {
		tester(t, input{
			event: indexer.GnoEvent{
				Type:    "Transfer",
				Func:    "Mint",
				PkgPath: "gno.land/r/gnoswap/v1/test_token/bar",
				Attrs: []indexer.GnoEventAttr{
					{
						Key:   "from",
						Value: "",
					},
					{
						Key:   "to",
						Value: "g17290cwvmrapvp869xfnhhawa8sm9edpufzat7d",
					},
					{
						Key:   "value",
						Value: "100000000000000",
					},
				},
			},
		}, expected{
			valid: true,
		})
		tester(t, input{
			event: indexer.GnoEvent{
				Type:    "Transfer",
				Func:    "Mint",
				PkgPath: "gno.land/r/gnoswap/v1/test_token/bar",
				Attrs: []indexer.GnoEventAttr{
					{
						Key:   "from",
						Value: "g17290cwvmrapvp869xfnhhawa8sm9edpufzat7d",
					},
					{
						Key:   "to",
						Value: "g17290cwvmrapvp869xfnhhawa8sm9edpufzat7d",
					},
					{
						Key:   "value",
						Value: "100000000000000",
					},
				},
			},
		}, expected{
			valid: false,
		})
	})

	t.Run("Burn", func(t *testing.T) {
		tester(t, input{
			event: indexer.GnoEvent{
				Type:    "Transfer",
				Func:    "Burn",
				PkgPath: "gno.land/r/gnoswap/v1/test_token/bar",
				Attrs: []indexer.GnoEventAttr{
					{
						Key:   "from",
						Value: "g17290cwvmrapvp869xfnhhawa8sm9edpufzat7d",
					},
					{
						Key:   "to",
						Value: "",
					},
					{
						Key:   "value",
						Value: "100000000000000",
					},
				},
			},
		}, expected{
			valid: true,
		})

		tester(t, input{
			event: indexer.GnoEvent{
				Type:    "Transfer",
				Func:    "Burn",
				PkgPath: "gno.land/r/gnoswap/v1/test_token/bar",
				Attrs: []indexer.GnoEventAttr{
					{
						Key:   "from",
						Value: "g17290cwvmrapvp869xfnhhawa8sm9edpufzat7d",
					},
					{
						Key:   "to",
						Value: "g17290cwvmrapvp869xfnhhawa8sm9edpufzat7d",
					},
					{
						Key:   "value",
						Value: "100000000000000",
					},
				},
			},
		}, expected{
			valid: false,
		})
	})

	t.Run("Transfer", func(t *testing.T) {
		tester(t, input{
			event: indexer.GnoEvent{
				Type:    "Transfer",
				Func:    "Transfer",
				PkgPath: "gno.land/r/gnoswap/v1/test_token/bar",
				Attrs: []indexer.GnoEventAttr{
					{
						Key:   "from",
						Value: "g16a7etgm9z2r653ucl36rj0l2yqcxgrz2jyegzx",
					},
					{
						Key:   "to",
						Value: "g17290cwvmrapvp869xfnhhawa8sm9edpufzat7d",
					},
					{
						Key:   "value",
						Value: "100000000000000",
					},
				},
			},
		}, expected{
			valid: true,
		})
		tester(t, input{
			event: indexer.GnoEvent{
				Type:    "Transfer",
				Func:    "Transfer",
				PkgPath: "gno.land/r/gnoswap/v1/test_token/bar",
				Attrs: []indexer.GnoEventAttr{
					{
						Key:   "from",
						Value: "g17290cwvmrapvp869xfnhhawa8sm9edpufzat7d",
					},
					{
						Key:   "to",
						Value: "g17290cwvmrapvp869xfnhhawa8sm9edpufzat7d",
					},
					{
						Key:   "value",
						Value: "",
					},
				},
			},
		}, expected{
			valid: false,
		})
	})

}
