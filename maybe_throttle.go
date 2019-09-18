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
	"log"
	"time"

	"github.com/asig/duplikator/edam"
)

func maybeThrottle(err error) bool {
	if e, ok := err.(*edam.EDAMSystemException); ok {
		if e.ErrorCode == edam.EDAMErrorCode_RATE_LIMIT_REACHED {
			log.Printf("Rate limit reached: Sleeping for %d seconds.", e.GetRateLimitDuration())
			time.Sleep(time.Duration(e.GetRateLimitDuration()+1) * time.Second)
			return true
		}
	}
	return false
}

