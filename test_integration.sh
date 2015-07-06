#!/bin/bash

directory=$(dirname $0)

python3 -m unittest discover -s $directory/test
