#!/bin/bash

go build -o vecprgen || exit 1

n=1000000

total=0

while [ $? == 0 ]; do
	./vecprgen -n $n
	((total++))
	echo "Generated $total sets of $n vectors without error"
done

