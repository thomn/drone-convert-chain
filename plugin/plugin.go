// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"fmt"

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
func New(targets ...Target) converter.Plugin {
	plugins := make(map[string]converter.Plugin, len(targets))
	for _, target := range targets {
		plugins[target.Endpoint] = converter.Client(
			target.Endpoint,
			target.Signer,
			target.SkipVerify,
		)
	}
	return &plugin{targets: plugins}
}

type plugin struct {
	targets map[string]converter.Plugin
}

func (p *plugin) Convert(ctx context.Context, req *converter.Request) (*drone.Config, error) {
	// get the configuration file from the request.
	cfg := req.Config.Data

	logger := logrus.WithField("repo", req.Repo.Slug)
	for endpoint, target := range p.targets {
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

		cfg = config.Data
	}

	// returns the modified configuration file.
	return &drone.Config{
		Data: cfg,
	}, nil
}
