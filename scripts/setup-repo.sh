#!/bin/bash

mkdir -p ./tests && cd ./tests

# setup repository
git config --global user.email "dev@branchswapper.com"
git config --global user.name "Dev Branchswapper"
git init --initial-branch=main .
touch a.txt b.txt c.txt
git add .
git commit -m "Initial commit"
git branch first
git branch second
git branch third
