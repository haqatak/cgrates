/*
Real-time Online/Offline Charging System (OCS) for Telecom & ISP environments
Copyright (C) ITsysCOM GmbH

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>
*/

package engine

import (
	"testing"

	"github.com/cgrates/cgrates/utils"
)

func TestPostgressGetStorageType(t *testing.T) {
	poS := &PostgresStorage{}
	storageType := poS.GetStorageType()
	if storageType != utils.MetaPostgres {
		t.Errorf("expected %s, got %s", utils.MetaPostgres, storageType)
	}
}

func TestExtraFieldsQueries(t *testing.T) {
	poS := &PostgresStorage{}
	field := "Subject"
	value := "1001"
	expectedExistsQuery := " extra_fields ?? ?"
	expectedValueQuery := " (extra_fields ->> ?) = ?"
	existsQuery, existsParams := poS.extraFieldsExistsQry(field)
	valueQuery, valueParams := poS.extraFieldsValueQry(field, value)
	if existsQuery != expectedExistsQuery {
		t.Errorf("extraFieldsExistsQry: expected query to be %s, but got %s", expectedExistsQuery, existsQuery)
	}
	if len(existsParams) != 1 || existsParams[0] != field {
		t.Errorf("extraFieldsExistsQry params: expected [%s], got %v", field, existsParams)
	}
	if valueQuery != expectedValueQuery {
		t.Errorf("extraFieldsValueQry: expected query to be %s, but got %s", expectedValueQuery, valueQuery)
	}
	if len(valueParams) != 2 || valueParams[0] != field || valueParams[1] != value {
		t.Errorf("extraFieldsValueQry params: expected [%s, %s], got %v", field, value, valueParams)
	}
}

func TestPostgresNotExtraFieldsValueQry(t *testing.T) {
	poS := &PostgresStorage{}
	field := "Tor"
	value := "voice"
	expectedQuery := " NOT (extra_fields ?? ? AND (extra_fields ->> ?) = ?)"
	query, params := poS.notExtraFieldsValueQry(field, value)
	if query != expectedQuery {
		t.Errorf("expected query to be %s, but got %s", expectedQuery, query)
	}
	if len(params) != 3 || params[0] != field || params[1] != field || params[2] != value {
		t.Errorf("notExtraFieldsValueQry params mismatch: got %v", params)
	}
}

func TestPostgresNotExtraFieldsExistsQry(t *testing.T) {
	poS := &PostgresStorage{}
	field := "tor"
	expectedQuery := " NOT extra_fields ?? ?"
	query, params := poS.notExtraFieldsExistsQry(field)
	if query != expectedQuery {
		t.Errorf("expected query to be %s, but got %s", expectedQuery, query)
	}
	if len(params) != 1 || params[0] != field {
		t.Errorf("notExtraFieldsExistsQry params mismatch: got %v", params)
	}
}
