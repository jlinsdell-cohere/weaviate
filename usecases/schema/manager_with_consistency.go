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

package schema

import (
	"context"

	"github.com/weaviate/weaviate/entities/models"
)

// ManagerWithConsistency expose the same interface as Manager but with the consistency flag.
// This is used to ensure that internal users will not miss-use the flag and it doesn't need to be set to a default
// value everytime we use the Manager.
type ManagerWithConsistency struct {
	*Manager
}

func NewManagerWithConsistency(handler *Manager) ManagerWithConsistency {
	return ManagerWithConsistency{Manager: handler}
}

// GetClass overrides the default implementation to consider the consistency flag
func (m *ManagerWithConsistency) GetClass(ctx context.Context, principal *models.Principal,
	name string, consistency bool,
) (*models.Class, error) {
	if err := m.Authorizer.Authorize(principal, "list", "schema/*"); err != nil {
		return nil, err
	}
	if consistency {
		class, err := m.metaWriter.QueryReadOnlyClass(name)
		return &class, err
	} else {
		return m.metaReader.ReadOnlyClass(name), nil
	}
}
