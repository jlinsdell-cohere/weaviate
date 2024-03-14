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
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/weaviate/weaviate/adapters/repos/db/vector/common"
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

func TestComposer(t *testing.T) {
	dimensions := 20
	vectors_size := 10_000
	queries_size := 10
	k := 10

	vectors, queries := testinghelpers.RandomVecs(vectors_size, queries_size, dimensions)
	rootPath := t.TempDir()
	distancer := distancer.NewL2SquaredProvider()
	truths := make([][]uint64, queries_size)
	compressionhelpers.Concurrently(uint64(len(queries)), func(i uint64) {
		truths[i], _ = testinghelpers.BruteForce(vectors, queries[i], k, distanceWrapper(distancer))
	})
	noopCallback := cyclemanager.NewCallbackGroupNoop()
	fuc := flatent.UserConfig{}
	fuc.SetDefaults()
	hnswuc := hnswent.UserConfig{
		MaxConnections:        30,
		EFConstruction:        64,
		EF:                    32,
		VectorCacheMaxObjects: 1_000_000,
	}
	composer, err := composer.New(composer.Config{
		RootPath:              rootPath,
		ID:                    "nil-vector-test",
		MakeCommitLoggerThunk: hnsw.MakeNoopCommitLogger,
		DistanceProvider:      distancer,
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
	}, ent.UserConfig{
		Threeshold: uint64(vectors_size),
		Distance:   distancer.Type(),
		HnswUC:     hnswuc,
		FlatUC:     fuc,
	}, testinghelpers.NewDummyStore(t))
	assert.Nil(t, err)

	compressionhelpers.Concurrently(uint64(vectors_size), func(i uint64) {
		composer.Add(i, vectors[i])
	})
	shouldUpgrade, at := composer.ShouldUpgrade()
	assert.True(t, shouldUpgrade)
	assert.Equal(t, vectors_size, at)
	assert.False(t, composer.Upgraded())
	recall1, latency1 := recallAndLatency(queries, k, composer, truths)
	fmt.Println(recall1, latency1)
	assert.True(t, recall1 > 0.99)
	wg := sync.WaitGroup{}
	wg.Add(1)
	composer.Upgrade(func() {
		wg.Done()
	})
	wg.Wait()
	shouldUpgrade, _ = composer.ShouldUpgrade()
	assert.False(t, shouldUpgrade)
	recall2, latency2 := recallAndLatency(queries, k, composer, truths)
	fmt.Println(recall2, latency2)
	assert.True(t, recall2 > 0.9)
	assert.True(t, latency1 > latency2)
}

func recallAndLatency(queries [][]float32, k int, index composer.VectorIndex, truths [][]uint64) (float32, float32) {
	var relevant uint64
	var retrieved int

	var querying time.Duration = 0
	compressionhelpers.Concurrently(uint64(len(queries)), func(i uint64) {
		before := time.Now()
		results, _, _ := index.SearchByVector(queries[i], k, nil)
		querying += time.Since(before)
		retrieved += k
		relevant += testinghelpers.MatchesInLists(truths[i], results)
	})

	recall := float32(relevant) / float32(retrieved)
	latency := float32(querying.Microseconds()) / float32(len(queries))
	return recall, latency
}

func TempVectorForIDThunk(vectors [][]float32) func(context.Context, uint64, *common.VectorSlice) ([]float32, error) {
	return func(ctx context.Context, id uint64, container *common.VectorSlice) ([]float32, error) {
		copy(container.Slice, vectors[int(id)])
		return vectors[int(id)], nil
	}
}

func distanceWrapper(provider distancer.Provider) func(x, y []float32) float32 {
	return func(x, y []float32) float32 {
		dist, _, _ := provider.SingleDist(x, y)
		return dist
	}
}
