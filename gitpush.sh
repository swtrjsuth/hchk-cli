#!/bin/bash

MSG=""
if [[ "$1" ]] ; then MSG="$1" ; else MSG="Commit $(date +%s)" ; fi
git add --all
git commit -am "$MSG"
git push git@github.com:swtrjsuth/hchk-cli.git main

