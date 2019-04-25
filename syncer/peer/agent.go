package peer

import (
	"context"
	"io"
	"github.com/apache/servicecomb-service-center/pkg/log"
	"os"

	"github.com/apache/servicecomb-service-center/pkg/gopool"
	"github.com/hashicorp/serf/cmd/serf/command/agent"
	"github.com/hashicorp/serf/serf"
)

type Agent struct {
	*agent.Agent
	conf *Config
}

func Create(peerConf *Config, logOutput io.Writer) (*Agent, error) {
	if logOutput == nil {
		logOutput = os.Stderr
	}

	serConf, err := peerConf.convertToSerf()
	if err != nil {
		return nil, err
	}

	serfAgent, err := agent.Create(peerConf.Config, serConf, logOutput)
	if err != nil {
		return nil, err
	}
	return &Agent{Agent: serfAgent, conf: peerConf}, nil
}

func (a *Agent) Start(ctx context.Context) {
	err := a.Agent.Start()
	if err != nil {
		log.Errorf(err, "start peer agent failed")
	}

	gopool.Go(func(ctx context.Context) {
		select {
		case <-ctx.Done():
			a.Shutdown()
		}
	})
}

func (a *Agent) Leave() error {
	return a.Agent.Leave()
}

func (a *Agent) Shutdown() error {
	return a.Agent.Leave()
}

func (a *Agent) ShutdownCh() <-chan struct{} {
	return a.Agent.ShutdownCh()
}

func (a *Agent) LocalMember() *serf.Member {
	serfAgent := a.Agent.Serf()
	if serfAgent != nil {
		member := serfAgent.LocalMember()
		return &member
	}
	serfAgent.State()
	return nil
}

func (a *Agent) Member(node string) *serf.Member {
	serfAgent := a.Agent.Serf()
	if serfAgent != nil {
		ms := serfAgent.Members()
		for _, m := range ms {
			if m.Name == node {
				return &m
			}
		}
	}
	return nil
}

func (a *Agent) SerfConfig() *serf.Config {
	return a.Agent.SerfConfig()
}

func (a *Agent) Join(addrs []string, replay bool) (n int, err error) {
	return a.Agent.Join(addrs, replay)
}

func (a *Agent) ForceLeave(node string) error {
	return a.Agent.ForceLeave(node)
}

func (a *Agent) UserEvent(name string, payload []byte, coalesce bool) error {
	return a.Agent.UserEvent(name, payload, coalesce)
}

func (a *Agent) Query(name string, payload []byte, params *serf.QueryParam) (*serf.QueryResponse, error) {
	return a.Agent.Query(name, payload, params)
}
