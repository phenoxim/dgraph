/*
 * Copyright 2019 Dgraph Labs, Inc. and Contributors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package alpha

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// contains checks whether given element is contained
// in any of the elements of the given list of strings.
func contains(ps []string, p string) bool {
	var res bool
	for _, v := range ps {
		res = res || strings.Contains(v, p)
	}

	return res
}

func TestUpsertExample0(t *testing.T) {
	require.NoError(t, dropAll())
	require.NoError(t, alterSchema(`email: string @index(exact) .`))

	// Mutation with wrong name
	m1 := `
upsert {
  mutation {
    set {
      uid(v) <name> "Wrong" .
      uid(v) <email> "ashish@dgraph.io" .
    }
  }

  query {
    me(func: eq(email, "ashish@dgraph.io")) {
      v as uid
    }
  }
}`
	keys, preds, _, err := mutationWithTs(m1, "application/rdf", false, true, true, 0)
	require.NoError(t, err)
	require.True(t, len(keys) == 0)
	require.True(t, contains(preds, "email"))
	require.True(t, contains(preds, "name"))

	// query should return the wrong name
	q1 := `
{
  q(func: has(email)) {
    uid
    name
    email
  }
}`
	res, _, err := queryWithTs(q1, "application/graphql+-", 0)
	require.NoError(t, err)
	require.Contains(t, res, "Wrong")

	// mutation with correct name
	m2 := `
upsert {
  mutation {
    set {
      uid(v) <name> "Ashish" .
    }
  }

  query {
    me(func: eq(email, "ashish@dgraph.io")) {
      v as uid
    }
  }
}`
	keys, preds, _, err = mutationWithTs(m2, "application/rdf", false, true, true, 0)
	require.NoError(t, err)
	require.True(t, len(keys) == 0)
	require.True(t, contains(preds, "name"))

	// query should return correct name
	res, _, err = queryWithTs(q1, "application/graphql+-", 0)
	require.NoError(t, err)
	require.Contains(t, res, "Ashish")
}