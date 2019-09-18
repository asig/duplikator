/*
 * Copyright (c) 2019 Andreas Signer <asigner@gmail.com>
 *
 * This file is part of Duplikator.
 *
 * Duplikator is free software: you can redistribute it and/or
 * modify it under the terms of the GNU General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * Duplikator is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Duplikator.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"testing"
)

func TestMakeFilename(t *testing.T) {
	tests := []struct {
		name string
		raw        string
		expected string
	}{
		{"clean name","hello", "hello"},
		{"needs cleaning", "he\\llo", "he_llo"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := makeFilename(test.raw);
			if got != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, got);
			}
		})
	}
}