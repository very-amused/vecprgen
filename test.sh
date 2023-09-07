#!/bin/bash

go build -o vecprgen || exit 1

n=100

total=0

while true; do
	./vecprgen -debug -n $n || exit 1
	((total++))
	echo "Generated $total sets of $n vectors without error"
done

