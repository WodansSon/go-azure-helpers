// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package identity

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var _ json.Marshaler = &SystemOrSingleUserAssignedMap{}

type SystemOrSingleUserAssignedMap struct {
	Type        Type                                   `json:"type" tfschema:"type"`
	PrincipalId string                                 `json:"principalId" tfschema:"principal_id"`
	TenantId    string                                 `json:"tenantId" tfschema:"tenant_id"`
	IdentityIds map[string]UserAssignedIdentityDetails `json:"userAssignedIdentities"`
}

func (s *SystemOrSingleUserAssignedMap) MarshalJSON() ([]byte, error) {
	// we use a custom marshal function here since we can only send the Type / UserAssignedIdentities field
	identityType := TypeNone
	userAssignedIdentityIds := map[string]UserAssignedIdentityDetails{}

	if s != nil {
		if s.Type == TypeSystemAssigned {
			identityType = TypeSystemAssigned
		}
		if s.Type == TypeUserAssigned {
			identityType = TypeUserAssigned
		}

		if identityType != TypeNone {
			userAssignedIdentityIds = s.IdentityIds
		}
	}

	out := map[string]interface{}{
		"type":                   string(identityType),
		"userAssignedIdentities": nil,
	}
	if len(userAssignedIdentityIds) > 0 {
		out["userAssignedIdentities"] = userAssignedIdentityIds
	}
	return json.Marshal(out)
}

// ExpandSystemOrSingleUserAssignedMap expands the schema input into a SystemOrSingleUserAssignedMap struct
func ExpandSystemOrSingleUserAssignedMap(input []interface{}) (*SystemOrSingleUserAssignedMap, error) {
	identityType := TypeNone
	identityIds := make(map[string]UserAssignedIdentityDetails, 0)

	if len(input) > 0 {
		raw := input[0].(map[string]interface{})
		typeRaw := raw["type"].(string)
		if typeRaw == string(TypeSystemAssigned) {
			identityType = TypeSystemAssigned
		}
		if typeRaw == string(TypeUserAssigned) {
			identityType = TypeUserAssigned
		}

		identityIdsRaw := raw["identity_ids"].(*schema.Set).List()
		for _, v := range identityIdsRaw {
			identityIds[v.(string)] = UserAssignedIdentityDetails{
				// intentionally empty since the expand shouldn't send these values
			}
		}
	}

	if identityType == TypeUserAssigned {
		if len(identityIds) == 0 {
			return nil, fmt.Errorf("`identity_ids` must be specified when `type` is set to %q", string(TypeUserAssigned))
		}

		if len(identityIds) > 1 {
			return nil, fmt.Errorf("`identity_ids` can only contain a single identity ID when `type` is set to %q", string(TypeUserAssigned))
		}
	}

	if len(identityIds) > 0 && identityType == TypeSystemAssigned {
		return nil, fmt.Errorf("`identity_ids` can only be specified when `type` is set to %q", string(TypeUserAssigned))
	}

	identity := &SystemOrSingleUserAssignedMap{
		Type:        identityType,
		IdentityIds: identityIds,
	}

	return identity, nil
}

// FlattenSystemOrSingleUserAssignedMap turns a SystemOrSingleUserAssignedMap into a []interface{}
func FlattenSystemOrSingleUserAssignedMap(input *SystemOrSingleUserAssignedMap) (*[]interface{}, error) {
	if input == nil {
		return &[]interface{}{}, nil
	}

	input.Type = normalizeType(input.Type)
	if input.Type != TypeSystemAssigned && input.Type != TypeUserAssigned {
		return &[]interface{}{}, nil
	}

	identityIds := make([]string, 0)
	for raw := range input.IdentityIds {
		id, err := commonids.ParseUserAssignedIdentityIDInsensitively(raw)
		if err != nil {
			return nil, fmt.Errorf("parsing %q as a User Assigned Identity ID: %+v", raw, err)
		}
		identityIds = append(identityIds, id.ID())
	}

	return &[]interface{}{
		map[string]interface{}{
			"type":         string(input.Type),
			"identity_ids": identityIds,
			"principal_id": input.PrincipalId,
			"tenant_id":    input.TenantId,
		},
	}, nil
}

// ExpandSystemOrSingleUserAssignedMapFromModel expands the typed schema input into a SystemOrSingleUserAssignedMap struct
func ExpandSystemOrSingleUserAssignedMapFromModel(input []ModelSystemAssignedUserAssigned) (*SystemOrSingleUserAssignedMap, error) {
	if len(input) == 0 {
		return &SystemOrSingleUserAssignedMap{
			Type:        TypeNone,
			IdentityIds: nil,
		}, nil
	}

	identity := input[0]

	identityIds := make(map[string]UserAssignedIdentityDetails, len(identity.IdentityIds))
	for _, v := range identity.IdentityIds {
		identityIds[v] = UserAssignedIdentityDetails{
			// intentionally empty since the expand shouldn't send these values
		}
	}

	if identity.Type == TypeUserAssigned {
		if len(identityIds) == 0 {
			return nil, fmt.Errorf("`identity_ids` must be specified when `type` is set to %q", string(TypeUserAssigned))
		}

		if len(identityIds) > 1 {
			return nil, fmt.Errorf("`identity_ids` can only contain a single identity ID when `type` is set to %q", string(TypeUserAssigned))
		}
	}

	if len(identityIds) > 0 && identity.Type == TypeSystemAssigned {
		return nil, fmt.Errorf("`identity_ids` can only be specified when `type` is set to %q", string(TypeUserAssigned))
	}

	return &SystemOrSingleUserAssignedMap{
		Type:        identity.Type,
		IdentityIds: identityIds,
	}, nil
}

// FlattenSystemOrSingleUserAssignedMapToModel turns a SystemOrSingleUserAssignedMap into a typed schema model
func FlattenSystemOrSingleUserAssignedMapToModel(input *SystemOrSingleUserAssignedMap) (*[]ModelSystemAssignedUserAssigned, error) {
	if input == nil {
		return &[]ModelSystemAssignedUserAssigned{}, nil
	}

	input.Type = normalizeType(input.Type)
	if input.Type != TypeSystemAssigned && input.Type != TypeUserAssigned {
		return &[]ModelSystemAssignedUserAssigned{}, nil
	}

	identityIds := make([]string, 0)
	for raw := range input.IdentityIds {
		id, err := commonids.ParseUserAssignedIdentityIDInsensitively(raw)
		if err != nil {
			return nil, fmt.Errorf("parsing %q as a User Assigned Identity ID: %+v", raw, err)
		}
		identityIds = append(identityIds, id.ID())
	}

	return &[]ModelSystemAssignedUserAssigned{
		{
			Type:        input.Type,
			IdentityIds: identityIds,
			PrincipalId: input.PrincipalId,
			TenantId:    input.TenantId,
		},
	}, nil
}
