package rrouter

import (
	"github.com/jackc/pgproto3"
	"github.com/pg-sharding/spqr/internal/qdb"
	"github.com/pg-sharding/spqr/internal/qrouter"
	"github.com/wal-g/tracelog"
	"golang.org/x/xerrors"
)

type RelayState struct {
	TxActive bool

	ActiveShards []qdb.ShardKey

	traceMsgs bool

	qr      qrouter.Qrouter
	cl      Client
	manager ConnManager

	msgBuf []*pgproto3.BackendMessage
}

func NewRelayState(qr qrouter.Qrouter, client Client, manager ConnManager) *RelayState {
	return &RelayState{
		ActiveShards: nil,
		TxActive:     false,
		msgBuf:       nil,
		traceMsgs:    false,
		qr:           qr,
		cl:           client,
		manager:      manager,
	}
}

func (rst *RelayState) Flush() {
	rst.msgBuf = nil
	rst.traceMsgs = false
}

func (rst *RelayState) Reroute(q *pgproto3.Query) error {

	shards := rst.qr.Route(q.String)

	if len(shards) == 0 {
		_ = rst.manager.UnRouteWithError(rst.cl, nil, "failed to match shard")
		return xerrors.New("failed to match shard")
	}

	if err := rst.manager.UnRouteCB(rst.cl, rst.ActiveShards); err != nil {
		tracelog.ErrorLogger.PrintError(err)
	}
	rst.ActiveShards = shards

	var serv Server
	var err error

	if len(shards) > 1 {
		serv, err = NewMultiShardServer(rst.cl.Route().beRule, rst.cl.Route().servPool)
		if err != nil {
			return err
		}
	} else {
		serv = NewShardServer(rst.cl.Route().beRule, rst.cl.Route().servPool)
	}

	if err := rst.cl.AssignServerConn(serv); err != nil {
		return err
	}

	tracelog.InfoLogger.Printf("route cl %s:%s to %v", rst.cl.Usr(), rst.cl.DB(), shards)

	if err := rst.manager.RouteCB(rst.cl, rst.ActiveShards); err != nil {
		tracelog.ErrorLogger.PrintError(err)
	}

	return nil
}

const TXREL = 73

func (rst *RelayState) RelayStep() error {

	if !rst.TxActive {
		if err := rst.manager.TXBeginCB(rst.cl, rst); err != nil {
			return err
		}
		rst.TxActive = true
	}

	return nil
}
func (rst *RelayState) ShouldRetry() bool {
	return false
}

func (rst *RelayState) CompleteRelay(txst byte) error {

	tracelog.InfoLogger.Printf("complete relay iter with TX status %v", txst)

	if err := rst.cl.Send(&pgproto3.ReadyForQuery{}); err != nil {
		return err
	}

	if txst == TXREL {
		if rst.TxActive {
			if err := rst.manager.TXEndCB(rst.cl, rst); err != nil {
				return err
			}
			rst.TxActive = false
		}
	} else {
		if !rst.TxActive {
			if err := rst.manager.TXBeginCB(rst.cl, rst); err != nil {
				return err
			}
			rst.TxActive = true
		}
	}

	return nil
}