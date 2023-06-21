package app

import (
	"fmt"
	"io"

	"github.com/jackc/pgproto3/v2"
	"github.com/pg-sharding/spqr/pkg/config"
	"github.com/pg-sharding/spqr/pkg/conn"
	"github.com/pg-sharding/spqr/pkg/spqrlog"
	"github.com/pg-sharding/spqr/router/client"
	"github.com/pg-sharding/spqr/router/parser"
	"github.com/pg-sharding/spqr/router/qrouter"
	"github.com/pg-sharding/spqr/router/rulerouter"
	"github.com/pg-sharding/spqr/router/server"
	"github.com/spaolacci/murmur3"
)

type Qinteractor interface{}

type QinteractorImpl struct{}

func AdvancedPoolModeNeeded(rst rulerouter.RelayStateMgr) bool {
	return rst.Client().Rule().PoolMode == config.PoolModeTransaction && rst.Client().Rule().PoolPreparedStatement || rst.RouterMode() == config.ProxyMode
}

func procQuery(rst rulerouter.RelayStateMgr, query string, msg pgproto3.FrontendMessage, cmngr rulerouter.PoolMgr) error {
	spqrlog.Logger.Printf(spqrlog.DEBUG1, "received query '%v' from %p", query, rst.Client())
	state, comment, err := rst.Parse(query)
	if err != nil {
		return err
	}

	mp, err := parser.ParseComment(comment)
	if err == nil {
		// if val, ok := mp["sharding_key"]; ok {
		// 	ds, err := qr.deparseKeyWithRangesInternal(ctx, val)
		// 	if err != nil {
		// 		return SkipRoutingState{}, err
		// 	}
		// 	return ShardMatchState{
		// 		Routes: []*DataShardRoute{ds},
		// 	}, nil
		// }
		if val, ok := mp["target-session-attrs"]; ok {
			// TBD: validate

			spqrlog.Logger.Printf(spqrlog.DEBUG2, "parse tsa %s from comment", val)

			rst.Client().SetTsa(val)
		}
	}

	switch st := state.(type) {
	case parser.ParseStateTXBegin:
		rst.AddSilentQuery(msg)
		rst.Client().StartTx()

		if err := rst.Client().Send(&pgproto3.CommandComplete{
			CommandTag: []byte("BEGIN"),
		}); err != nil {
			return err
		}

		return rst.Client().Send(&pgproto3.ReadyForQuery{
			TxStatus: byte(conn.TXACT),
		})
	case parser.ParseStateTXCommit:
		if !cmngr.ConnectionActive(rst) {
			return fmt.Errorf("no connection to shards")
		}
		rst.AddQuery(msg)
		ok, err := rst.ProcessMessageBuf(true, true, cmngr)
		if ok {
			rst.Client().CommitActiveSet()
		}
		return err
	case parser.ParseStateTXRollback:
		rst.AddQuery(msg)
		ok, err := rst.ProcessMessageBuf(true, true, cmngr)
		if ok {
			rst.Client().Rollback()
		}
		return err
	case parser.ParseStateEmptyQuery:
		if err := rst.Client().Send(&pgproto3.EmptyQueryResponse{}); err != nil {
			return err
		}
		return rst.Client().Send(&pgproto3.ReadyForQuery{
			TxStatus: byte(rst.TxStatus()),
		})
	// with tx pooling we might have no active connection while processing set x to y
	case parser.ParseStateSetStmt:
		rst.AddQuery(msg)
		if ok, err := rst.ProcessMessageBuf(true, true, cmngr); err != nil {
			return err
		} else if ok {
			rst.Client().SetParam(st.Name, st.Value)
		}
		return nil
	case parser.ParseStateResetStmt:
		rst.Client().ResetParam(st.Name)

		if cmngr.ConnectionActive(rst) {
			if err := rst.ProcessMessage(rst.Client().ConstructClientParams(), true, false, cmngr); err != nil {
				return err
			}
		}

		if err := rst.Client().Send(&pgproto3.CommandComplete{CommandTag: []byte("RESET")}); err != nil {
			return err
		}

		return rst.Client().Send(&pgproto3.ReadyForQuery{
			TxStatus: byte(rst.TxStatus()),
		})
	case parser.ParseStateResetMetadataStmt:
		if cmngr.ConnectionActive(rst) {
			rst.AddQuery(msg)
			_, err := rst.ProcessMessageBuf(true, true, cmngr)
			if err != nil {
				return err
			}

			rst.Client().ResetParam(st.Setting)
			if st.Setting == "session_authorization" {
				rst.Client().ResetParam("role")
			}

			return nil
		}

		rst.Client().ResetParam(st.Setting)
		if st.Setting == "session_authorization" {
			rst.Client().ResetParam("role")
		}

		if err := rst.Client().Send(&pgproto3.CommandComplete{CommandTag: []byte("RESET")}); err != nil {
			return err
		}

		return rst.Client().Send(&pgproto3.ReadyForQuery{
			TxStatus: byte(rst.TxStatus()),
		})
	case parser.ParseStateResetAllStmt:
		rst.Client().ResetAll()

		if cmngr.ConnectionActive(rst) {
			if err := rst.Client().Send(&pgproto3.CommandComplete{CommandTag: []byte("RESET")}); err != nil {
				return err
			}
		}

		return rst.Client().Send(&pgproto3.ReadyForQuery{
			TxStatus: byte(rst.TxStatus()),
		})
	case parser.ParseStateSetLocalStmt:
		if cmngr.ConnectionActive(rst) {
			rst.AddQuery(msg)
			_, err := rst.ProcessMessageBuf(true, true, cmngr)
			return err
		}
		if err := rst.Client().Send(&pgproto3.CommandComplete{CommandTag: []byte("SET")}); err != nil {
			return err
		}
		return rst.Client().Send(&pgproto3.ReadyForQuery{
			TxStatus: byte(rst.TxStatus()),
		})
	case parser.ParseStatePrepareStmt:
		// sql level prepares stmt pooling
		if AdvancedPoolModeNeeded(rst) {
			spqrlog.Logger.Printf(spqrlog.DEBUG1, "sql level prep statement pooling support is on")
			rst.Client().StorePreparedStatement(st.Name, st.Query)
			return rst.Client().ReplyParseComplete()
		} else {
			rst.AddQuery(msg)
			_, err := rst.ProcessMessageBuf(true, true, cmngr)
			return err
		}
	case parser.ParseStateExecute:
		if AdvancedPoolModeNeeded(rst) {
			// do nothing
			rst.Client().PreparedStatementQueryByName(st.Name)
			return nil
		} else {
			rst.AddQuery(msg)
			_, err := rst.ProcessMessageBuf(true, true, cmngr)
			return err
		}
	case parser.ParseStateExplain:
		_ = rst.Client().ReplyErrMsg("parse state explain is not implemented")
		return nil
	default:
		rst.AddQuery(msg)
		_, err := rst.ProcessMessageBuf(true, true, cmngr)
		return err
	}
}

// ProcessMessage
func ProcessMessage(qr qrouter.QueryRouter, cl client.RouterClient, cmngr rulerouter.PoolMgr, rst rulerouter.RelayStateMgr, msg pgproto3.FrontendMessage) error {
	if cl.Rule().PoolMode == config.PoolModeTransaction && !cl.Rule().PoolPreparedStatement {
		switch q := msg.(type) {
		case *pgproto3.Terminate:
			return nil
		case *pgproto3.Sync:
			// copy interface
			cpQ := *q
			q = &cpQ
			return rst.ProcessMessage(q, true, true, cmngr)
		case *pgproto3.FunctionCall:
			// copy interface
			cpQ := *q
			q = &cpQ
			return rst.ProcessMessage(q, true, true, cmngr)
		case *pgproto3.Parse:
			// copy interface
			cpQ := *q
			q = &cpQ
			return procQuery(rst, q.Query, q, cmngr)
		case *pgproto3.Execute:
			// copy interface
			cpQ := *q
			q = &cpQ
			return rst.ProcessMessage(q, false, true, cmngr)
		case *pgproto3.Bind:
			// copy interface
			cpQ := *q
			q = &cpQ
			return rst.ProcessMessage(q, false, true, cmngr)
		case *pgproto3.Describe:
			// copy interface
			cpQ := *q
			q = &cpQ
			return rst.ProcessMessage(q, false, true, cmngr)
		case *pgproto3.Query:
			// copy interface
			cpQ := *q
			q = &cpQ
			return procQuery(rst, q.String, q, cmngr)
		default:
			return nil
		}
	}

	switch q := msg.(type) {
	case *pgproto3.Terminate:
		return nil
	case *pgproto3.Sync:
		return rst.Sync(true, true, cmngr)
	case *pgproto3.Parse:
		// copy interface
		cpQ := *q
		q = &cpQ

		hash := murmur3.Sum64([]byte(q.Query))
		spqrlog.Logger.Printf(spqrlog.DEBUG1, "name %v, query %v, hash %d", q.Name, q.Query, hash)
		if err := cl.ReplyDebugNoticef("name %v, query %v, hash %d", q.Name, q.Query, hash); err != nil {
			return err
		}
		cl.StorePreparedStatement(q.Name, q.Query)
		// simply reply witch ok parse complete
		return cl.ReplyParseComplete()
	case *pgproto3.Describe:
		// copy interface
		cpQ := *q
		q = &cpQ

		if q.ObjectType == 'P' {
			if err := rst.ProcessMessage(q, true, true, cmngr); err != nil {
				return err
			}
			return nil
		}
		query := cl.PreparedStatementQueryByName(q.Name)
		hash := murmur3.Sum64([]byte(query))

		if err := rst.PrepareRelayStep(cl, cmngr); err != nil {
			return err
		}

		q.Name = fmt.Sprintf("%d", hash)
		if err := rst.PrepareStatement(hash, server.PrepStmtDesc{
			Name:  q.Name,
			Query: query,
		}); err != nil {
			return err
		}

		var err error
		if err = rst.RelayRunCommand(q, false, false); err != nil {
			if rst.ShouldRetry(err) {
				return fmt.Errorf("retry logic for prepared statements is not implemented")
			}
		}
		return err
	case *pgproto3.FunctionCall:
		// copy interface
		cpQ := *q
		q = &cpQ
		spqrlog.Logger.Printf(spqrlog.DEBUG1, "client %p function call: simply fire parse stmt to connection", cl)
		return rst.ProcessMessage(q, false, true, cmngr)
	case *pgproto3.Execute:
		// copy interface
		cpQ := *q
		q = &cpQ
		spqrlog.Logger.Printf(spqrlog.DEBUG1, "client %p execute prepared statement: simply fire parse stmt to connection", cl)
		return rst.ProcessMessage(q, true, true, cmngr)
	case *pgproto3.Bind:
		// copy interface
		cpQ := *q
		q = &cpQ
		query := cl.PreparedStatementQueryByName(q.PreparedStatement)
		hash := murmur3.Sum64([]byte(query))

		if err := rst.PrepareRelayStep(cl, cmngr); err != nil {
			return err
		}

		if err := rst.PrepareStatement(hash, server.PrepStmtDesc{
			Name:  fmt.Sprintf("%d", hash),
			Query: query,
		}); err != nil {
			return err
		}

		q.PreparedStatement = fmt.Sprintf("%d", hash)

		return rst.RelayRunCommand(q, false, true)
	case *pgproto3.Query:
		// copy interface
		cpQ := *q
		q = &cpQ
		return procQuery(rst, q.String, q, cmngr)
	default:
		return nil
	}
}

func Frontend(qr qrouter.QueryRouter, cl client.RouterClient, cmngr rulerouter.PoolMgr, rcfg *config.Router) error {
	spqrlog.Logger.Printf(spqrlog.INFO, "process frontend for route %s %s (client %p)", cl.Usr(), cl.DB(), cl)

	_ = cl.ReplyDebugNoticef("process frontend for route %s %s", cl.Usr(), cl.DB())
	rst := rulerouter.NewRelayState(qr, cl, cmngr, rcfg)

	defer rst.Close()

	var msg pgproto3.FrontendMessage
	var err error

	for {
		msg, err = cl.Receive()
		if err != nil {
			switch err {
			case io.ErrUnexpectedEOF:
				fallthrough
			case io.EOF:
				return nil
				// ok
			default:
				return rst.UnRouteWithError(rst.ActiveShards(), err)
			}
		}

		if err := ProcessMessage(qr, cl, cmngr, rst, msg); err != nil {
			switch err {
			case io.ErrUnexpectedEOF:
				fallthrough
			case io.EOF:
				return nil
				// ok
			default:
				// fix all reply err to client to be here
				spqrlog.Logger.Printf(spqrlog.DEBUG5, "client %p iter done with error: %v", rst.Client(), err)
			}
		}
	}
}