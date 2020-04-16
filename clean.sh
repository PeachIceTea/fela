#!/bin/sh

echo Deleting: executables
rm fela __debug_bin 2&>/dev/null

echo Deleting: files
rm -rf files/audio/* files/cover/* 2&>/dev/null

echo Clearing: Database
dbmate drop
dbmate up
