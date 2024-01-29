#!/usr/bin/env sh


echo "Welcome to a short script to migrate all the themes from themes/*.yml to toml with alacritty migrate"


for dirName in themes/*; do
    echo "$dirName"
    alacritty migrate --config-file="$dirName"
    rm "$dirName"
done
