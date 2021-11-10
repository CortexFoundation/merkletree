// Copyright 2021 The CortexTheseus Authors
// This file is part of the CortexTheseus library.
//
// The CortexTheseus library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The CortexTheseus library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the CortexTheseus library. If not, see <http://www.gnu.org/licenses/>.

// Package core implements the Cortex consensus protocol

package merkletree

import (
	"bytes"
	"github.com/CortexFoundation/CortexTheseus/common"
	"testing"
)

var (
	testContents = []testContent{
		{common.HexToHash("0000000000000000000000000000000000000000000000000000000000000000")},
		{common.HexToHash("0000000000000000000000000000000000000000000000000000000000000001")},
		{common.HexToHash("0000000000000000000000000000000000000000000000000000000000000002")},
		{common.HexToHash("0000000000000000000000000000000000000000000000000000000000000003")},
		{common.HexToHash("0000000000000000000000000000000000000000000000000000000000000004")},
		{common.HexToHash("0000000000000000000000000000000000000000000000000000000000000005")},
		{common.HexToHash("0000000000000000000000000000000000000000000000000000000000000006")},
		{common.HexToHash("0000000000000000000000000000000000000000000000000000000000000007")},
		{common.HexToHash("0000000000000000000000000000000000000000000000000000000000000008")},
		{common.HexToHash("0000000000000000000000000000000000000000000000000000000000000009")},
	}
)

type testContent struct {
	common.Hash
}

func (tc testContent) CalculateHash() ([]byte, error) {
	return tc.Bytes(), nil
}

func (tc testContent) Equals(other Content) (bool, error) {
	h1, err := tc.CalculateHash()
	if err != nil {
		return false, err
	}
	h2, err := other.CalculateHash()
	if err != nil {
		return false, err
	}
	return bytes.Equal(h1, h2), nil
}

func TestNewTree(t *testing.T) {
	root, err := NewTree([]Content{testContents[0], testContents[1]})
	if err != nil {
		t.Error("new tree error: ", err)
	}
	rootHash := common.BytesToHash(root.Root.Hash).Hex()
	want := "0xa6eef7e35abe7026729641147f7915573c7e97b47efa546f5f6e3230263bcb49"
	if rootHash != want {
		t.Errorf("root unmatched. should be %s, got %s", want, rootHash)
	}
	//t.Log(root.String())
}

func TestMerkleTree_AddNode(t *testing.T) {
	t_rebuild, err := NewTree([]Content{testContents[0]})
	if err != nil {
		t.Fatal("new tree error: ", err)
	}

	t_add, err := NewTree([]Content{testContents[0]})
	if err != nil {
		t.Fatal("new tree error: ", err)
	}
	var tmp []Content
	for i := 1; i < 10; i += 1 {
		tmp = nil
		for j := 0; j <= i; j++ {
			tmp = append(tmp, testContents[j])
		}
		t_rebuild.RebuildTreeWith(tmp)
		rebuild_hash := common.BytesToHash(t_rebuild.merkleRoot)

		t_add.AddNode(testContents[i])
		add_hash := common.BytesToHash(t_add.merkleRoot)

		if v, err := t_add.VerifyTree(); !v || err != nil {
			t.Fatalf("root add not verified, at %d", i)
		}

		if v, err := t_rebuild.VerifyTree(); !v || err != nil {
			t.Fatalf("root add not verified, at %d", i)
		}

		if add_hash != rebuild_hash || len(t_add.Leafs) != len(t_rebuild.Leafs) {
			t.Log("Rebuild:" + t_rebuild.String())
			t.Log("AddNode:" + t_add.String())
			t.Fatalf("root unmatched at %d. rebuild hash is %s, add hash is %s", i, rebuild_hash, add_hash)
		}
		//t.Log(t_add.String())
		prettyPrint(t_rebuild.Root, 0)
	}
	//	t.Log(t_add.String())
	//print2DUtil(t_add.Root, 0)
	//	print2DUtil(t_rebuild.Root, 0)
}
