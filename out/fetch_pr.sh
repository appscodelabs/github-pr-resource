#!/bin/bash
cd $1/$2
git ls-remote origin 'pull/*/head' | grep -F -f <(git rev-parse HEAD) | awk -F'/' '{print $3}'
