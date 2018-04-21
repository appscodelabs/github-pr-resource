#!/bin/bash
#pwd
#echo $1 #url
#echo $2 #dir
#echo $3 #pull_id
git clone $1 $2
cd $2
git checkout master
git fetch origin pull/$3/head:$3
git checkout $3
