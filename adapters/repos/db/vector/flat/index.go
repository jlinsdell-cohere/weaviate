//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2023 Weaviate B.V. All rights reserved.
//
//  CONTACT: hello@weaviate.io
//

package flat

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/weaviate/weaviate/adapters/repos/db/helpers"
	"github.com/weaviate/weaviate/adapters/repos/db/lsmkv"
	"github.com/weaviate/weaviate/adapters/repos/db/vector/common"
	"github.com/weaviate/weaviate/adapters/repos/db/vector/hnsw"
	"github.com/weaviate/weaviate/adapters/repos/db/vector/hnsw/distancer"
	"github.com/weaviate/weaviate/adapters/repos/db/vector/ssdhelpers"
	"github.com/weaviate/weaviate/entities/cyclemanager"
	"github.com/weaviate/weaviate/entities/schema"
	"github.com/weaviate/weaviate/entities/storobj"
	flatent "github.com/weaviate/weaviate/entities/vectorindex/flat"
	"github.com/weaviate/weaviate/usecases/floatcomp"
)

type flat struct {
	id                  string
	dims                int32
	store               *lsmkv.Store
	logger              logrus.FieldLogger
	distancerProvider   distancer.Provider
	shardName           string
	trackDimensionsOnce sync.Once
	ef                  int64

	TempVectorForIDThunk common.TempVectorForID
	vectorForID          common.VectorForID[float32]
	bq                   ssdhelpers.BinaryQuantizer

	tempVectors *common.TempVectorsPool
	pqResults   *common.PqMaxPool
	pool        *pools

	shardCompactionCallbacks cyclemanager.CycleCallbackGroup
	shardFlushCallbacks      cyclemanager.CycleCallbackGroup

	compression string
	flatStore   func() *lsmkv.CursorReplace
}

func New(
	cfg hnsw.Config, uc flatent.UserConfig,
	shardCompactionCallbacks, shardFlushCallbacks cyclemanager.CycleCallbackGroup,
	flatStore func() *lsmkv.CursorReplace,
) (*flat, error) {
	if err := cfg.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid config")
	}

	if cfg.Logger == nil {
		logger := logrus.New()
		logger.Out = io.Discard
		cfg.Logger = logger
	}

	index := &flat{
		id:                       cfg.ID,
		logger:                   cfg.Logger,
		distancerProvider:        cfg.DistanceProvider,
		ef:                       int64(uc.EF),
		shardName:                cfg.ShardName,
		TempVectorForIDThunk:     cfg.TempVectorForIDThunk,
		vectorForID:              cfg.VectorForIDThunk,
		tempVectors:              common.NewTempVectorsPool(),
		pqResults:                common.NewPqMaxPool(100),
		compression:              uc.Compression,
		shardCompactionCallbacks: shardCompactionCallbacks,
		shardFlushCallbacks:      shardFlushCallbacks,
		flatStore:                flatStore,
	}
	index.initStore(cfg.RootPath, cfg.ClassName)

	return index, nil
}

func (h *flat) storeCompressedVector(index uint64, vector []byte) {
	Id := make([]byte, 8)
	binary.LittleEndian.PutUint64(Id, index)
	h.store.Bucket(helpers.CompressedObjectsBucketLSM).Put(Id, vector)
}

func (index *flat) initStore(rootPath, className string) error {
	if index.compression == flatent.CompressionNone {
		return nil
	}

	store, err := lsmkv.New(fmt.Sprintf("%s/%s/%s", rootPath, className, index.shardName), "", index.logger, nil,
		index.shardCompactionCallbacks, index.shardFlushCallbacks)
	if err != nil {
		return errors.Wrap(err, "Init lsmkv (compressed vectors store)")
	}
	err = store.CreateOrLoadBucket(context.Background(), helpers.CompressedObjectsBucketLSM)
	if err != nil {
		return errors.Wrapf(err, "Create or load bucket (compressed vectors store)")
	}
	index.store = store
	return nil
}

func (index *flat) AddBatch(ids []uint64, vectors [][]float32) error {
	/*if len(ids) != len(vectors) {
		return errors.Errorf("ids and vectors sizes does not match")
	}
	if len(ids) == 0 {
		return errors.Errorf("insertBatch called with empty lists")
	}
	index.trackDimensionsOnce.Do(func() {
		atomic.StoreInt32(&index.dims, int32(len(vectors[0])))
		fmt.Println("dimensions: " + string(index.dims))
		index.bq = *ssdhelpers.NewBinaryQuantizer()
		index.bq.Fit(vectors)
	})
	for idx := range ids {
		if err := index.Add(ids[idx], vectors[idx]); err != nil {
			return err
		}
	}*/
	return nil
}

func byteSliceFromUint64Slice(x []uint64, slice []byte) []byte {
	for i := range x {
		binary.LittleEndian.PutUint64(slice[i*8:], x[i])
	}
	return slice
}

func uint64SliceFromByteSlice(x []byte, slice []uint64) []uint64 {
	len := len(x) / 8
	for i := 0; i < len; i++ {
		slice[i] = binary.LittleEndian.Uint64(x[i*8:])
	}
	return slice
}

func (index *flat) Add(id uint64, vector []float32) error {
	index.trackDimensionsOnce.Do(func() {
		atomic.StoreInt32(&index.dims, int32(len(vector)))
		index.bq = ssdhelpers.NewBinaryQuantizer()
		index.pool = newPools()
		if index.compression == flatent.CompressionNone {
			return
		}
		atomic.StoreInt32(&index.dims, int32(len(vector)))
		index.bq = ssdhelpers.NewBinaryQuantizer()
	})
	if index.compression == flatent.CompressionNone {
		return nil
	}
	if len(vector) != int(index.dims) {
		return errors.Errorf("insert called with a vector of the wrong size")
	}
	if index.distancerProvider.Type() == "cosine-dot" {
		// cosine-dot requires normalized vectors, as the dot product and cosine
		// similarity are only identical if the vector is normalized
		vector = distancer.Normalize(vector)
	}
	vec, err := index.bq.Encode(vector)
	if err != nil {
		return err
	}
	slice := make([]byte, len(vec)*8)
	index.storeCompressedVector(id, byteSliceFromUint64Slice(vec, slice))

	return nil
}

func (index *flat) Delete(ids ...uint64) error {
	return nil
}

func (index *flat) searchTimeEF(k int) int {
	if index.compression == flatent.CompressionNone {
		return k
	}
	// load atomically, so we can get away with concurrent updates of the
	// userconfig without having to set a lock each time we try to read - which
	// can be so common that it would cause considerable overhead
	ef := int(atomic.LoadInt64(&index.ef))
	if ef < k {
		ef = k
	}
	return ef
}

func (index *flat) SearchByVector(vector []float32, k int, allow helpers.AllowList) ([]uint64, []float32, error) {
	if index.distancerProvider.Type() == "cosine-dot" {
		// cosine-dot requires normalized vectors, as the dot product and cosine
		// similarity are only identical if the vector is normalized
		vector = distancer.Normalize(vector)
	}

	ef := index.searchTimeEF(k)
	heap := index.pqResults.GetMax(ef)
	defer index.pqResults.Put(heap)

	firstId := uint64(0)
	alreadyFound := 0
	if allow != nil {
		firstId, _ = allow.Iterator().Next()
	}

	if index.compression != flatent.CompressionNone {
		query, err := index.bq.Encode(vector)
		if err != nil {
			return nil, nil, err
		}
		cursor := index.store.Bucket(helpers.CompressedObjectsBucketLSM).Cursor()
		var key []byte
		var v []byte
		if allow != nil {
			buff := make([]byte, 16)
			binary.BigEndian.PutUint64(buff[8:], firstId)
			key, v = cursor.Seek(buff)
		} else {
			key, v = cursor.First()
		}
		for key != nil && (allow == nil || alreadyFound < allow.Len()) {
			alreadyFound++
			id := binary.LittleEndian.Uint64(key)
			if allow != nil && !allow.Contains(id) {
				continue
			}
			t := index.pool.uint64SlicePool.Get(len(v) / 8)
			candidate := uint64SliceFromByteSlice(v, t.slice)
			d, _ := index.bq.DistanceBetweenCompressedVectors(candidate, query)

			index.pool.uint64SlicePool.Put(t)
			if heap.Len() < ef || heap.Top().Dist > d {
				if heap.Len() == ef {
					heap.Pop()
				}
				heap.Insert(id, d)
			}
			key, v = cursor.Next()
		}
		cursor.Close()

		ids := make([]uint64, ef)
		for j := range ids {
			ids[j] = heap.Pop().ID
		}

		for _, id := range ids {
			candidate, err := index.vectorForID(context.Background(), id)
			if err != nil {
				return nil, nil, err
			}
			d, _, _ := index.distancerProvider.SingleDist(candidate, vector)
			if heap.Len() < ef || heap.Top().Dist > d {
				if heap.Len() == ef {
					heap.Pop()
				}
				heap.Insert(uint64(id), d)
			}
		}
		for heap.Len() > k {
			heap.Pop()
		}
	} else {
		cursor := index.flatStore()
		var key []byte
		var v []byte
		if allow != nil {
			buff := make([]byte, 16)
			binary.BigEndian.PutUint64(buff[8:], firstId)
			key, v = cursor.Seek(buff)
		} else {
			key, v = cursor.First()
		}

		for key != nil && (allow == nil || alreadyFound < allow.Len()) {
			alreadyFound++
			obj, err := storobj.FromBinary(v)
			id := obj.DocID()
			if allow != nil && !allow.Contains(id) {
				continue
			}
			if err != nil {
				return nil, nil, errors.Wrapf(err, "unmarhsal item %d", id)
			}
			candidate := obj.Vector
			if index.distancerProvider.Type() == "cosine-dot" {
				// cosine-dot requires normalized vectors, as the dot product and cosine
				// similarity are only identical if the vector is normalized
				candidate = distancer.Normalize(candidate)
			}

			d, _, _ := index.distancerProvider.SingleDist(vector, candidate)

			if heap.Len() < ef || heap.Top().Dist > d {
				if heap.Len() == ef {
					heap.Pop()
				}
				heap.Insert(id, d)
			}
			key, v = cursor.Next()
		}
		cursor.Close()
	}

	if heap.Len() < k {
		k = heap.Len()
	}
	ids := make([]uint64, k)
	dists := make([]float32, k)
	for j := k - 1; j >= 0; j-- {
		elem := heap.Pop()
		ids[j] = elem.ID
		dists[j] = elem.Dist
	}

	return ids, dists, nil
}

func (index *flat) SearchByVectorDistance(vector []float32, targetDistance float32, maxLimit int64, allow helpers.AllowList) ([]uint64, []float32, error) {
	var (
		searchParams = newSearchByDistParams(maxLimit)

		resultIDs  []uint64
		resultDist []float32
	)

	recursiveSearch := func() (bool, error) {
		shouldContinue := false

		ids, dist, err := index.SearchByVector(vector, searchParams.TotalLimit(), allow)
		if err != nil {
			return false, errors.Wrap(err, "vector search")
		}

		// ensures the indexers aren't out of range
		offsetCap := searchParams.OffsetCapacity(ids)
		totalLimitCap := searchParams.TotalLimitCapacity(ids)

		ids, dist = ids[offsetCap:totalLimitCap], dist[offsetCap:totalLimitCap]

		if len(ids) == 0 {
			return false, nil
		}

		lastFound := dist[len(dist)-1]
		shouldContinue = lastFound <= targetDistance

		for i := range ids {
			if aboveThresh := dist[i] <= targetDistance; aboveThresh ||
				floatcomp.InDelta(float64(dist[i]), float64(targetDistance), 1e-6) {
				resultIDs = append(resultIDs, ids[i])
				resultDist = append(resultDist, dist[i])
			} else {
				// as soon as we encounter a certainty which
				// is below threshold, we can stop searching
				break
			}
		}

		return shouldContinue, nil
	}

	shouldContinue, err := recursiveSearch()
	if err != nil {
		return nil, nil, err
	}

	for shouldContinue {
		searchParams.Iterate()
		if searchParams.MaxLimitReached() {
			index.logger.
				WithField("action", "unlimited_vector_search").
				Warnf("maximum search limit of %d results has been reached",
					searchParams.MaximumSearchLimit())
			break
		}

		shouldContinue, err = recursiveSearch()
		if err != nil {
			return nil, nil, err
		}
	}

	return resultIDs, resultDist, nil
}

func (index *flat) UpdateUserConfig(updated schema.VectorIndexConfig, callback func()) error {
	parsed, ok := updated.(flatent.UserConfig)
	if !ok {
		callback()
		return errors.Errorf("config is not UserConfig, but %T", updated)
	}

	// Store automatically as a lock here would be very expensive, this value is
	// read on every single user-facing search, which can be highly concurrent
	atomic.StoreInt64(&index.ef, int64(parsed.EF))

	callback()
	return nil
}

func (index *flat) Drop(ctx context.Context) error {
	return nil
}

func (index *flat) Flush() error {
	return nil
}

func (index *flat) Shutdown(ctx context.Context) error {
	if err := index.store.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "flat shutdown")
	}
	return index.Drop(ctx)
}

func (index *flat) SwitchCommitLogs(context.Context) error {
	return nil
}

func (index *flat) ListFiles(context.Context) ([]string, error) {
	return nil, nil
}

func (i *flat) ValidateBeforeInsert(vector []float32) error {
	return nil
}

func (index *flat) PostStartup() {
	index.prefillCache()
}

func (h *flat) prefillCache() {
	/*limit := 0
	limit = int(h.cache.CopyMaxSize())

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
		defer cancel()

		var err error
		err = newVectorCachePrefiller(h.cache, h, h.logger).Prefill(ctx, limit)

		if err != nil {
			h.logger.WithError(err).Error("prefill vector cache")
		}
	}()*/
}

func (index *flat) Dump(labels ...string) {
	if len(labels) > 0 {
		fmt.Printf("--------------------------------------------------\n")
		fmt.Printf("--  %s\n", strings.Join(labels, ", "))
	}
	fmt.Printf("--------------------------------------------------\n")
	fmt.Printf("ID: %s\n", index.id)
	fmt.Printf("--------------------------------------------------\n")
}

func (index *flat) DistanceBetweenVectors(x, y []float32) (float32, bool, error) {
	return index.distancerProvider.SingleDist(x, y)
}

func (index *flat) ContainsNode(id uint64) bool {
	return true
}

func (index *flat) DistancerProvider() distancer.Provider {
	return index.distancerProvider
}

func (index *flat) getCompressedVectorForID(ctx context.Context, id uint64) ([]uint64, error) {
	slice := index.tempVectors.Get(int(index.dims))
	vec, err := index.TempVectorForIDThunk(context.Background(), id, slice)
	index.tempVectors.Put(slice)
	if err != nil {
		return nil, errors.Wrap(err, "Getting vector for id")
	}
	if index.distancerProvider.Type() == "cosine-dot" {
		// cosine-dot requires normalized vectors, as the dot product and cosine
		// similarity are only identical if the vector is normalized
		vec = distancer.Normalize(vec)
	}

	return index.bq.Encode(vec)
}

func newSearchByDistParams(maxLimit int64) *common.SearchByDistParams {
	initialOffset := 0
	initialLimit := common.DefaultSearchByDistInitialLimit

	return common.NewSearchByDistParams(initialOffset, initialLimit, initialOffset+initialLimit, maxLimit)
}

func ValidateUserConfigUpdate(initial, updated schema.VectorIndexConfig) error {
	/*
		initialParsed, ok := initial.(flatent.UserConfig)
		if !ok {
			return errors.Errorf("initial is not UserConfig, but %T", initial)
		}

		updatedParsed, ok := updated.(flatent.UserConfig)
		if !ok {
			return errors.Errorf("updated is not UserConfig, but %T", updated)
		}

		immutableFields := []immutableParameter{
			{
				name:     "distance",
				accessor: func(c flatent.UserConfig) interface{} { return c.Distance },
			},
		}

		for _, u := range immutableFields {
			if err := validateImmutableField(u, initialParsed, updatedParsed); err != nil {
				return err
			}
		}
	*/
	return nil
}