#!/bin/bash
#
# Copyright (c) 2019 Andreas Signer <asigner@gmail.com>
#
# This file is part of Duplikator.
#
# Duplikator is free software: you can redistribute it and/or
# modify it under the terms of the GNU General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.
#
# Duplikator is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with Duplikator.  If not, see <http://www.gnu.org/licenses/>.

SRCDIR=/tmp/evernote-sdk
rm -rf ${SRCDIR}
git clone https://github.com/evernote/evernote-thrift.git ${SRCDIR}

for f in UserStore.thrift NoteStore.thrift; do
  thrift \
    -strict \
    -nowarn \
    --allow-64bit-consts \
    --allow-neg-keys \
    --gen go:package_prefix=github.com/asig/evernote-backup/,thrift_import=github.com/apache/thrift/lib/go/thrift \
    -r \
    -I ${SRCDIR} \
    --out .  \
    ${SRCDIR}/src/${f}
done

