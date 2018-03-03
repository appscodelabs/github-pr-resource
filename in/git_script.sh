#!/bin/bash
#pwd
#echo $1 #url
#echo $2 #dir
#echo $3 #pull_id
#echo $4 #ref
git clone $1 $2
cd $2
git checkout master
git pull
git fetch origin pull/$3/head:$4
git checkout $4
