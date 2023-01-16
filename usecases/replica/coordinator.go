//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2022 SeMI Technologies B.V. All rights reserved.
//
//  CONTACT: hello@semi.technology
//

package replica

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

// readyOp asks a replica to be read to second phase commit
type readyOp func(ctx context.Context, host, requestID string) error

// readyOp asks a replica to execute the actual operation
type commitOp[T any] func(ctx context.Context, host, requestID string) (T, error)

// coordinator coordinates replication of write request
type coordinator[T any] struct {
	Client              // needed to commit and abort operation
	resolver  *resolver // host names of replicas
	class     string
	shard     string
	requestID string
	// responses collect all responses of batch job
	responses []T
	// nodes     []string
}

func newCoordinator[T any](r *Replicator, shard, requestID string) *coordinator[T] {
	return &coordinator[T]{
		Client: r.client,
		resolver: &resolver{
			schema:       r.stateGetter,
			nodeResolver: r.resolver,
			class:        r.class,
		},
		class:     r.class,
		shard:     shard,
		requestID: requestID,
	}
}

// TODO add level

// broadcast sends write request to all replicas (first phase of a two-phase commit)
func (c *coordinator[T]) broadcast(ctx context.Context, replicas []string, op readyOp) error {
	errs := make([]error, len(replicas))
	var g errgroup.Group
	for i, replica := range replicas {
		i, replica := i, replica
		g.Go(func() error {
			errs[i] = op(ctx, replica, c.requestID)
			return nil
		})
	}
	g.Wait()
	var err error
	for _, err = range errs {
		if err != nil {
			break
		}
	}

	if err != nil {
		for _, node := range replicas {
			c.Abort(ctx, node, c.class, c.shard, c.requestID)
		}
	}

	return err
}

// commitAll tells replicas to commit pending updates related to a specific request
// (second phase of a two-phase commit)
func (c *coordinator[T]) commitAll(ctx context.Context, replicas []string, op commitOp[T]) error {
	var g errgroup.Group
	c.responses = make([]T, len(replicas))
	errs := make([]error, len(replicas))
	for i, replica := range replicas {
		i, replica := i, replica
		g.Go(func() error {
			resp, err := op(ctx, replica, c.requestID)
			c.responses[i], errs[i] = resp, err
			return nil
		})
	}
	g.Wait()
	var err error
	for _, err = range errs {
		if err != nil {
			return err
		}
	}

	return nil
}

// Replicate writes on all replicas of specific shard
func (c *coordinator[T]) Replicate(ctx context.Context, ask readyOp, com commitOp[T]) error {
	state, err := c.resolver.State(c.shard)
	if err == nil {
		_, err = state.ConsistencyLevel(All)
	}
	const msg = "replication with consistency level 'ALL'"
	if err != nil {
		return fmt.Errorf("%s: %w : class %q shard %q", msg, err, c.class, c.shard)
	}
	if err := c.broadcast(ctx, state.Hosts, ask); err != nil {
		return fmt.Errorf("%s: broadcast: %w", msg, err)
	}
	if err := c.commitAll(context.Background(), state.Hosts, com); err != nil {
		return fmt.Errorf("%s commit: %w", msg, err)
	}
	return nil
}

// broadcast sends write request to all replicas (first phase of a two-phase commit)
func (c *coordinator[T]) broadcast2(ctx context.Context, replicas []string, op readyOp, level int) ([]string, error) {
	errs := make([]error, len(replicas))
	activeReplicas := make([]string, 0, len(replicas))
	var g errgroup.Group
	for i, replica := range replicas {
		i, replica := i, replica
		g.Go(func() error {
			errs[i] = op(ctx, replica, c.requestID)
			return errs[i]
		})
	}
	firstErr := g.Wait()
	for i, err := range errs {
		if err == nil {
			activeReplicas = append(activeReplicas, replicas[i])
		}
	}
	if len(activeReplicas) < level {
		firstErr = fmt.Errorf("no enough active replicas found: %w", firstErr)
	} else {
		firstErr = nil
	}

	if firstErr != nil {
		for _, node := range replicas {
			c.Abort(ctx, node, c.class, c.shard, c.requestID)
		}
	}

	return activeReplicas, firstErr
}

// commitAll tells replicas to commit pending updates related to a specific request
// (second phase of a two-phase commit)
func (c *coordinator[T]) commitAll2(ctx context.Context, replicas []string, op commitOp[T]) <-chan simpleResult[T] {
	replyCh := make(chan simpleResult[T], len(replicas))
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(len(replicas))
		for _, replica := range replicas {
			go func(replica string) {
				defer wg.Done()
				resp, err := op(ctx, replica, c.requestID)
				replyCh <- simpleResult[T]{resp, err}
			}(replica)
		}
		wg.Wait()
		close(replyCh)
	}()

	return replyCh
}

// Replicate writes on all replicas of specific shard
func (c *coordinator[T]) Replicate2(ctx context.Context, cl ConsistencyLevel, ask readyOp, com commitOp[T]) (<-chan simpleResult[T], int, error) {
	state, err := c.resolver.State(c.shard)
	level := 0
	if err == nil {
		level, err = state.ConsistencyLevel(cl)
	}
	const msg = "replication with consistency level 'ALL'" // TODO maybe move to upper level
	if err != nil {
		return nil, level, fmt.Errorf("%s: %w : class %q shard %q", msg, err, c.class, c.shard)
	}
	nodes, err := c.broadcast2(ctx, state.Hosts, ask, level)
	if err != nil {
		return nil, level, fmt.Errorf("%s: broadcast: %w", msg, err)
	}
	return c.commitAll2(context.Background(), nodes, com), level, nil
}
