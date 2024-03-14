//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2024 Weaviate B.V. All rights reserved.
//
//  CONTACT: hello@weaviate.io
//

package composer_test

import (
	"context"
	"math"
	"sync"
	"testing"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate/adapters/repos/db/vector/composer"
	"github.com/weaviate/weaviate/adapters/repos/db/vector/compressionhelpers"
	"github.com/weaviate/weaviate/adapters/repos/db/vector/hnsw"
	"github.com/weaviate/weaviate/adapters/repos/db/vector/hnsw/distancer"
	"github.com/weaviate/weaviate/adapters/repos/db/vector/testinghelpers"
	"github.com/weaviate/weaviate/entities/cyclemanager"
	"github.com/weaviate/weaviate/entities/storobj"
	ent "github.com/weaviate/weaviate/entities/vectorindex/composer"
	flatent "github.com/weaviate/weaviate/entities/vectorindex/flat"
	hnswent "github.com/weaviate/weaviate/entities/vectorindex/hnsw"
)

func TestBackup_Integration(t *testing.T) {
	dimensions := 20
	vectors_size := 10_000
	queries_size := 10
	k := 10

	vectors, queries := testinghelpers.RandomVecs(vectors_size, queries_size, dimensions)
	truths := make([][]uint64, queries_size)
	distancer := distancer.NewL2SquaredProvider()
	compressionhelpers.Concurrently(uint64(len(queries)), func(i uint64) {
		truths[i], _ = testinghelpers.BruteForce(vectors, queries[i], k, distanceWrapper(distancer))
	})
	logger, _ := test.NewNullLogger()

	dirName := t.TempDir()
	indexID := "restore-integration-test"
	noopCallback := cyclemanager.NewCallbackGroupNoop()
	fuc := flatent.UserConfig{}
	fuc.SetDefaults()
	hnswuc := hnswent.UserConfig{
		MaxConnections:        30,
		EFConstruction:        64,
		EF:                    32,
		VectorCacheMaxObjects: 1_000_000,
	}

	config := composer.Config{
		RootPath:         dirName,
		ID:               indexID,
		Logger:           logger,
		DistanceProvider: distancer,
		MakeCommitLoggerThunk: func() (hnsw.CommitLogger, error) {
			return hnsw.NewCommitLogger(dirName, indexID, logger, noopCallback)
		},
		VectorForIDThunk: func(ctx context.Context, id uint64) ([]float32, error) {
			vec := vectors[int(id)]
			if vec == nil {
				return nil, storobj.NewErrNotFoundf(id, "nil vec")
			}
			return vec, nil
		},
		TempVectorForIDThunk:     TempVectorForIDThunk(vectors),
		TombstoneCallbacks:       noopCallback,
		ShardCompactionCallbacks: noopCallback,
		ShardFlushCallbacks:      noopCallback,
	}

	uc := ent.UserConfig{
		Threeshold: uint64(vectors_size),
		Distance:   distancer.Type(),
		HnswUC:     hnswuc,
		FlatUC:     fuc,
	}

	store := testinghelpers.NewDummyStore(t)

	idx, err := composer.New(config, uc, store)
	require.Nil(t, err)
	idx.PostStartup()

	compressionhelpers.Concurrently(uint64(vectors_size), func(i uint64) {
		idx.Add(i, vectors[i])
	})

	wg := sync.WaitGroup{}
	wg.Add(1)
	idx.Upgrade(func() {
		wg.Done()
	})
	wg.Wait()
	recall1, _ := recallAndLatency(queries, k, idx, truths)
	assert.True(t, recall1 > 0.9)

	assert.Nil(t, idx.Shutdown(context.Background()))
	idx, err = composer.New(config, uc, store)
	require.Nil(t, err)
	idx.PostStartup()

	recall2, _ := recallAndLatency(queries, k, idx, truths)
	assert.True(t, math.Abs(float64(recall1-recall2)) <= 0.1)
}
