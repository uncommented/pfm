#!/usr/bin/env bash

sudo apt-get update -y && sudo apt-get install pre-commit -y
cd /pfm
pre-commit install
