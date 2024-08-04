#!/usr/bin/env bash

sudo apt-get update -y && sudo apt-get install pipx -y
pipx install pre-commit
pipx ensurepath
