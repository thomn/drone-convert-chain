// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package main

import (
	"io"
	"net/http"

	"github.com/drone/drone-go/plugin/converter"
	"github.com/thomn/drone-convert-chain/plugin"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// spec provides the plugin settings.
type spec struct {
	Bind   string `envconfig:"DRONE_BIND"`
	Debug  bool   `envconfig:"DRONE_DEBUG"`
	Secret string `envconfig:"DRONE_SECRET"`

	TargetEndpoints    []string `envconfig:"TARGET_ENDPOINTS"`
	TargetSecrets      []string `envconfig:"TARGET_SECRETS"`
	TargetSkipVerifies []bool   `envconfig:"TARGET_SKIP_VERIFIES"`
}

func main() {
	spec := new(spec)
	err := envconfig.Process("", spec)
	if err != nil {
		logrus.Fatal(err)
	}

	if spec.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if spec.Secret == "" {
		logrus.Fatalln("missing secret key")
	}
	if spec.Bind == "" {
		spec.Bind = ":3000"
	}

	targets := make([]plugin.Target, len(spec.TargetEndpoints))
	for i := range spec.TargetEndpoints {
		logrus.WithField("target", spec.TargetEndpoints[i]).Info("added target")
		targets[i] = plugin.Target{
			Endpoint:   spec.TargetEndpoints[i],
			Signer:     spec.TargetSecrets[i],
			SkipVerify: spec.TargetSkipVerifies[i],
		}
	}

	handler := converter.Handler(
		plugin.New(spec.Debug, targets...),
		spec.Secret,
		logrus.StandardLogger(),
	)

	logrus.Infof("server listening on address %s", spec.Bind)

	http.Handle("/", handler)
	http.HandleFunc("/healthz", healthz)
	logrus.Fatal(http.ListenAndServe(spec.Bind, nil))
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, "OK")
}
