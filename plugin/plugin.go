// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"fmt"
	"os"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/converter"
	"github.com/sirupsen/logrus"
)

type Target struct {
	Endpoint   string
	Signer     string
	SkipVerify bool
}

// New returns a new conversion plugin.
func New(debug bool, targets ...Target) converter.Plugin {
	plugins := make(map[string]converter.Plugin, len(targets))
	for _, target := range targets {
		plugins[target.Endpoint] = converter.Client(
			target.Endpoint,
			target.Signer,
			target.SkipVerify,
		)
	}
	return &plugin{
		targets: plugins,
		debug: debug,
	}
}

type plugin struct {
	targets map[string]converter.Plugin
	debug bool
}

func (p *plugin) Convert(ctx context.Context, req *converter.Request) (*drone.Config, error) {
	// get the configuration file from the request.
	cfg := req.Config.Data

	logger := logrus.WithField("repo", req.Repo.Slug)
	conversion := 0
	for endpoint, target := range p.targets {
		conversion++

		logger.WithField("endpoint", endpoint).Info("calling")
		config, err := target.Convert(ctx, &converter.Request{
			Repo:   req.Repo,
			Build:  req.Build,
			Config: drone.Config{Data: cfg},
		})
		if err != nil {
			logger.WithError(err).Error("unable to build")
			return nil, fmt.Errorf("unable to call %w", err)
		}

		dumpfile := fmt.Sprintf("/tmp/%s-%d-%d", req.Repo.Name, req.Build.ID, conversion)
		df, ferr := os.Create(dumpfile)
		if ferr == nil {
			_, _ = df.WriteString(cfg)
			_ = df.Close()
		}

		cfg = config.Data
	}

	// returns the modified configuration file.
	return &drone.Config{
		Data: cfg,
	}, nil
}
