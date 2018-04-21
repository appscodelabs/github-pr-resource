#!/bin/bash
cd $1/$2
git rev-parse --abbrev-ref HEAD
