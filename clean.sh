#!/bin/sh

echo Deleting: files/\*
rm files/* 2&>/dev/null

echo Clearing: Database
dbmate drop
dbmate up
