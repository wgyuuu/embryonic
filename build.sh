#!/bin/bash

if [ $# -lt 1 ]; then
    exit
fi

goPath=$GOPATH
fullPath="$goPath/src/$1"

mkdir -p $fullPath
cp -rf ./ $fullPath
rm -rf $fullPath/.git
find $fullPath -name '.gitignore' -delete

orgPath='github.com\/wgyuuu\/embryonic'

OLD_IFS="$IFS"
IFS='/'
arr=($1)
IFS="$OLD_IFS"
for a in ${arr[@]}; do
    if [ -z $programPath ]; then
        programPath=$a
    else
        programPath="$programPath\/$a"
    fi
done

find $fullPath -name '*.go' | xargs sed -i '' "s/$orgPath/$programPath/g"
