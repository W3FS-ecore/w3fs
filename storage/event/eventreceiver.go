package event

import (
	"context"
	"fmt"
	datatransfer "github.com/filecoin-project/go-data-transfer"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/go-fil-markets/storagemarket/impl/requestvalidation"
	"github.com/filecoin-project/go-statemachine/fsm"
	"github.com/filecoin-project/go-statestore"
	logging "github.com/ipfs/go-log/v2"
	"golang.org/x/xerrors"
	"sync"
)

var log = logging.Logger("event")

// EventReceiver is any thing that can receive FSM events
type EventReceiver interface {
	Open(id interface{}, args ...interface{}) (err error)
	Cancel(id interface{}, args ...interface{}) (err error)
	Restart(id interface{}, args ...interface{}) (err error)
	Disconnected(id interface{}, args ...interface{}) (err error)
	Error(id interface{}, args ...interface{}) (err error)
	Completed(id interface{}, args ...interface{}) (err error)
	TransferRequestQueued(id interface{}, args ...interface{}) (err error)
	Accept(id interface{}, args ...interface{}) (err error)
}

type DataTransferEventReceiver struct {
	Sms      sync.Map
	Notifs   sync.Map
	Dispatch fsm.Notifier
}

func NewDataTransferEventReceiver() DataTransferEventReceiver {
	return DataTransferEventReceiver{
		Sms:    sync.Map{},
		Notifs: sync.Map{},
	}
}

func (d *DataTransferEventReceiver) Open(id interface{}, args ...interface{}) (err error) {
	log.Info("data transfer channel open")
	return nil
}
func (d *DataTransferEventReceiver) Cancel(id interface{}, args ...interface{}) (err error) {
	log.Info("data transfer channel cancel")
	a, ok := d.Notifs.Load(statestore.ToKey(id))
	if ok {
		notifs := a.(chan interface{})
		close(notifs)
	}
	return nil
}
func (d *DataTransferEventReceiver) Restart(id interface{}, args ...interface{}) (err error) {
	log.Info("data transfer channel restart")
	return nil
}
func (d *DataTransferEventReceiver) Disconnected(id interface{}, args ...interface{}) (err error) {
	log.Info("data transfer channel disconnected")
	a, ok := d.Notifs.Load(statestore.ToKey(id))
	if ok {
		notifs := a.(chan interface{})
		close(notifs)
	}
	return nil
}
func (d *DataTransferEventReceiver) Error(id interface{}, args ...interface{}) (err error) {
	log.Error("chid:", id, "data transfer channel error:", args)
	a, ok := d.Notifs.Load(statestore.ToKey(id))
	if ok {
		notifs := a.(chan interface{})
		close(notifs)
	}
	return nil
}
func (d *DataTransferEventReceiver) Completed(id interface{}, args ...interface{}) (err error) {
	log.Info("data transfer completed")
	a, ok := d.Notifs.Load(statestore.ToKey(id))
	if ok {
		notifs := a.(chan interface{})
		if notifs != nil {
			notifs <- args
		}
	}
	b, ok := d.Sms.Load(statestore.ToKey(id))
	if ok && d.Dispatch != nil {
		d.Dispatch(storagemarket.ClientEventDataTransferComplete, b)
	}
	return nil
}
func (d *DataTransferEventReceiver) TransferRequestQueued(id interface{}, args ...interface{}) (err error) {
	log.Info("data transfer channel TransferRequestQueued")
	return nil
}
func (d *DataTransferEventReceiver) Accept(id interface{}, args ...interface{}) (err error) {
	log.Info("data transfer channel Accept")
	return nil
}
func (d *DataTransferEventReceiver) Begin(id interface{}, args interface{}) (err error) {
	d.Sms.Store(statestore.ToKey(id), args)
	return nil
}
func (d *DataTransferEventReceiver) End(id interface{}, args ...interface{}) (err error) {
	d.Sms.Delete(statestore.ToKey(id))
	return nil
}

func (d *DataTransferEventReceiver) Get(id interface{}) (interface{}, bool) {
	a, ok := d.Sms.Load(statestore.ToKey(id))
	return a, ok
}

func (d *DataTransferEventReceiver) Stop(ctx context.Context) (err error) {
	log.Info("data transfer channel Stop")
	return nil
}

func (d *DataTransferEventReceiver) VerifiedData(id interface{}, args ...interface{}) (err error) {
	return nil
}

func (d *DataTransferEventReceiver) List(id interface{}) (err error) {
	a := make([]interface{}, 10)
	d.Sms.Range(func(_, v interface{}) bool {
		a = append(a, v)
		return true
	})
	id = a
	return nil
}

func (d *DataTransferEventReceiver) IsTerminated(id interface{}) bool {
	return true
}

func (d *DataTransferEventReceiver) GetSync(ctx context.Context, id interface{}, a interface{}) (err error) {
	a, _ = d.Sms.Load(statestore.ToKey(id))
	return nil
}

func (d *DataTransferEventReceiver) Wait(ctx context.Context, id interface{}) (err error) {
	notifs := make(chan interface{}, 1)
	d.Notifs.Store(statestore.ToKey(id), notifs)
	select {
	case _, ok := <-notifs:
		if !ok {
			return xerrors.Errorf("data transfer event receiver wait notifs channel closed")
		}
	case <-ctx.Done():
		return xerrors.Errorf("data transfer event receiver context deadline")
	}
	d.Notifs.Delete(statestore.ToKey(id))
	close(notifs)
	return nil
}

func DataTransferSubscriber(deals EventReceiver) datatransfer.Subscriber {
	return func(event datatransfer.Event, channelState datatransfer.ChannelState) {
		voucher, ok := channelState.Voucher().(*requestvalidation.StorageDataTransferVoucher)
		// if this event is for a transfer not related to storage, ignore
		if !ok {
			log.Debugw("ignoring data-transfer event as it's not storage related", "event", datatransfer.Events[event.Code], "channelID",
				channelState.ChannelID())
			return
		}

		log.Debugw("processing storage provider dt event", "event", datatransfer.Events[event.Code], "proposalCid", voucher.Proposal, "channelID",
			channelState.ChannelID(), "channelState", datatransfer.Statuses[channelState.Status()])

		if channelState.Status() == datatransfer.Completed {
			err := deals.Completed(voucher.Proposal, channelState.ChannelID())
			if err != nil {
				log.Errorf("processing dt event: %s", err)
			}
		}

		// Translate from data transfer events to provider FSM events
		// Note: We ignore data transfer progress events (they do not affect deal state)
		err := func() error {
			switch event.Code {
			case datatransfer.Cancel:
				return deals.Cancel(voucher.Proposal, channelState.ChannelID())
			case datatransfer.Restart:
				return deals.Restart(voucher.Proposal, channelState.ChannelID())
			case datatransfer.Disconnected:
				return deals.Disconnected(voucher.Proposal)
			case datatransfer.Open:
				return deals.Open(voucher.Proposal, channelState.ChannelID())
			case datatransfer.TransferRequestQueued:
				return deals.TransferRequestQueued(voucher.Proposal, channelState.ChannelID())
			case datatransfer.Accept:
				return deals.Accept(voucher.Proposal, channelState.ChannelID())
			case datatransfer.Error:
				return deals.Error(voucher.Proposal, fmt.Errorf("deal data transfer failed: %s", event.Message))
			default:
				return nil
			}
		}()
		if err != nil {
			log.Errorw("error processing storage provider dt event", "event", datatransfer.Events[event.Code], "proposalCid", voucher.Proposal, "channelID",
				channelState.ChannelID(), "err", err)
		}
	}
}
