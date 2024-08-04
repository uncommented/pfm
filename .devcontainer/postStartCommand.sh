#!/usr/bin/env bash

git config --global safe.directory /pfm

cd /pfm
pre-commit install
