// pqsd is an agent that connects to a postgresql cluster and manages stream emission.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"
	"encoding/json"

	_ "net/http/pprof"

	_ "golang.org/x/net/trace"

	"github.com/golang/protobuf/jsonpb"
	"github.com/google/gops/agent"
	_ "github.com/kardianos/minwinsvc" // import minwinsvc for windows service support
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	bigquery "cloud.google.com/go/bigquery"

	ctxutil "g.ghn.vn/logistic/bi/streaming/pg2-big-query/ctxutil"
	"g.ghn.vn/logistic/bi/streaming/pg2-big-query/pqs"
	"g.ghn.vn/logistic/bi/streaming/pg2-big-query"
	
)

var (
	verbose         = flag.Bool("v", false, "be verbose")
	postgresCluster = flag.String("connect", "", "postgresql cluster address")
	tableRegexp     = flag.String("tables", ".*", "regexp of tables to manage")
	remove          = flag.Bool("remove", false, "if true, remove triggers and exit")
	grpcAddr        = flag.String("addr", ":7000", "listen addr")
	debugAddr       = flag.String("debugaddr", ":7001", "listen debug addr")
	redactions      = flag.String("redactions", "", "details of fields to redact in JSON format i.e '{\"public\":{\"users\":[\"password\",\"ssn\"]}}'")
)

const (
	gracefulStopMaxWait = 10 * time.Second
)

func main() {
	flag.Parse()
	if err := run(ctxutil.BackgroundWithSignals()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	// starts the gops diagnostics agent
	if err := agent.Listen(agent.Options{
		ShutdownCleanup: true,
	}); err != nil {
		return err
	}

	tableRe, err := regexp.Compile(*tableRegexp)
	if err != nil {
		return err
	}

	opts := []pqstream.ServerOption{
		pqstream.WithTableRegexp(tableRe),
	}

	if (len(*redactions)) > 0 {
		rfields, err := pqstream.DecodeRedactions(*redactions)
		if err != nil {
			return errors.Wrap(err, "decoding redactions")
		}

		if len(rfields) > 0 {
			opts = append(opts, pqstream.WithFieldRedactions(rfields))
		}
	}
	if *verbose {
		l := logrus.New()
		l.Level = logrus.DebugLevel
		opts = append(opts, pqstream.WithLogger(l))
	}

	server, err := pqstream.NewServer(*postgresCluster, opts...)
	if err != nil {
		return err
	}

	err = errors.Wrap(server.RemoveTriggers(), "RemoveTriggers")
	if err != nil || *remove {
		return err
	}

	if err = server.InstallTriggers(); err != nil {
		return errors.Wrap(err, "InstallTriggers")
	}

	sb := func(e *pqs.Event) {
		m := &jsonpb.Marshaler{}
		sbd, err := m.MarshalToString(e.GetPayload())
		if err != nil {
			log.Fatalln(err)
		}
		byt := []byte(sbd)
		var datampp map[string]bigquery.Value
		if err := json.Unmarshal(byt, &datampp); err != nil {
			log.Fatalln(err)
		}
		fmt.Println("datampp:", datampp)
	}

	if err = server.HandleEvents(ctx, sb); err != nil {
		log.Fatalln(err)
	}
	return err
}

func bqInsert(djson map[string]interface{}){

}