// Copyright 2016 Andrew O'Neill

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mongo

import (
	"encoding/hex"
	"sync"

	"github.com/foolusion/choices"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Mongo struct {
	namespaces    []choices.Namespace
	sess          *mgo.Session
	url, db, coll string
	mu            sync.RWMutex
}

func WithMongoStorage(url, db, collection string) func(*choices.Config) error {
	return func(ec *choices.Config) error {
		m := &Mongo{url: url, db: db, coll: collection}
		sess, err := mgo.Dial(url)
		if err != nil {
			return err
		}
		m.sess = sess
		ec.Storage = m
		return nil
	}
}

type MongoNamespace struct {
	Name        string
	Segments    string
	TeamID      []string
	Experiments []MongoExperiment
}

type MongoExperiment struct {
	Name     string
	Segments string
	Params   []MongoParam
}

type MongoParam struct {
	Name  string
	Type  choices.ValueType
	Value bson.Raw
}

func (m *Mongo) Ready() error {
	return m.sess.Ping()
}

func (m *Mongo) Update() error {
	c := m.sess.DB(m.db).C(m.coll)
	iter := c.Find(bson.M{}).Iter()
	var mongoNamespaces []MongoNamespace
	err := iter.All(&mongoNamespaces)
	if err != nil {
		return err
	}

	namespaces := make([]choices.Namespace, len(mongoNamespaces))
	for i, n := range mongoNamespaces {
		namespaces[i] = choices.Namespace{
			Name:        n.Name,
			TeamID:      n.TeamID,
			Experiments: make([]choices.Experiment, len(n.Experiments)),
		}
		nsSeg, err := hex.DecodeString(n.Segments)
		if err != nil {
			return err
		}
		var nss [16]byte
		copy(nss[:], nsSeg[:16])
		namespaces[i].Segments = nss
		for j, e := range n.Experiments {
			namespaces[i].Experiments[j] = choices.Experiment{
				Name:   e.Name,
				Params: make([]choices.Param, len(e.Params)),
			}
			expSeg, err := hex.DecodeString(e.Segments)
			if err != nil {
				return err
			}
			var exps [16]byte
			copy(exps[:], expSeg[:16])
			namespaces[i].Experiments[j].Segments = exps

			for k, p := range e.Params {
				namespaces[i].Experiments[j].Params[k] = choices.Param{Name: p.Name}
				switch p.Type {
				case choices.ValueTypeUniform:
					var uniform choices.Uniform
					p.Value.Unmarshal(&uniform)
					namespaces[i].Experiments[j].Params[k].Value = &uniform
				case choices.ValueTypeWeighted:
					var weighted choices.Weighted
					p.Value.Unmarshal(&weighted)
					namespaces[i].Experiments[j].Params[k].Value = &weighted
				}
			}
		}
	}

	m.mu.Lock()
	m.namespaces = namespaces
	m.mu.Unlock()
	return nil
}

func (m *Mongo) Read() []choices.Namespace {
	m.mu.RLock()
	ns := m.namespaces
	m.mu.RUnlock()
	return ns
}
