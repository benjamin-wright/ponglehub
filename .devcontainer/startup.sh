#!/bin/bash

if ! cat ~/.bashrc | grep -q "direnv hook bash"; then
    echo "adding direnv to bashrc"
    echo "eval $(direnv hook bash)" >> ~/.bashrc

    . ~/.bashrc
else
    echo "direnv already in bashrc"
fi