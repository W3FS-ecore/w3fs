package requestvalidationex

import (
	"github.com/filecoin-project/go-fil-markets/storagemarket/impl/requestvalidation"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime"
	"github.com/libp2p/go-libp2p-core/peer"

	datatransfer "github.com/filecoin-project/go-data-transfer"

	"github.com/filecoin-project/go-fil-markets/storagemarket"
)

// PushDeals gets deal states for Push validations
type PushDeals interface {
	Get(cid.Cid) (storagemarket.MinerDeal, error)
}

// PullDeals gets deal states for Pull validations
type PullDeals interface {
	Get(cid.Cid) (storagemarket.ClientDeal, error)
}

// UnifiedRequestValidator is a data transfer request validator that validates
// StorageDataTransferVoucher from the given state store
// It can be made to only accept push requests (Provider) or pull requests (Client)
// by passing nil for the statestore value for pushes or pulls
type UnifiedRequestValidatorEx struct {
	pushDeals PushDeals
	pullDeals PullDeals
}

// NewUnifiedRequestValidator returns a new instance of UnifiedRequestValidator
func NewUnifiedRequestValidatorEx(pushDeals PushDeals, pullDeals PullDeals) *UnifiedRequestValidatorEx {
	return &UnifiedRequestValidatorEx{
		pushDeals: pushDeals,
		pullDeals: pullDeals,
	}
}

// SetPushDeals sets the store to look up push deals with
func (v *UnifiedRequestValidatorEx) SetPushDeals(pushDeals PushDeals) {
	v.pushDeals = pushDeals
}

// SetPullDeals sets the store to look up pull deals with
func (v *UnifiedRequestValidatorEx) SetPullDeals(pullDeals PullDeals) {
	v.pullDeals = pullDeals
}

// ValidatePush implements the ValidatePush method of a data transfer request validator.
// If no pushStore exists, it rejects the request
// Otherwise, it calls the ValidatePush function to validate the deal
func (v *UnifiedRequestValidatorEx) ValidatePush(isRestart bool, _ datatransfer.ChannelID, sender peer.ID, voucher datatransfer.Voucher, baseCid cid.Cid, selector ipld.Node) (datatransfer.VoucherResult, error) {
	if v.pushDeals == nil {
		return nil, requestvalidation.ErrNoPushAccepted
	}

	return nil, nil
}

// ValidatePull implements the ValidatePull method of a data transfer request validator.
// If no pullStore exists, it rejects the request
// Otherwise, it calls the ValidatePull function to validate the deal
func (v *UnifiedRequestValidatorEx) ValidatePull(isRestart bool, _ datatransfer.ChannelID, receiver peer.ID, voucher datatransfer.Voucher, baseCid cid.Cid, selector ipld.Node) (datatransfer.VoucherResult, error) {
	if v.pullDeals == nil {
		return nil, requestvalidation.ErrNoPullAccepted
	}

	return nil, nil
}

var _ datatransfer.RequestValidator = &UnifiedRequestValidatorEx{}
