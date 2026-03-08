#!/bin/bash

echo "deb https://hawkhackerf.github.io/yhujin-apt/repo stable main" \
| sudo tee /etc/apt/sources.list.d/yhujin.list

sudo apt update
sudo apt install -y yhujin
