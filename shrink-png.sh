#!/bin/sh
#
# NOTE: this only works on OSX due to the syntax of the stat command
#
# brew install pngcrush optipng

inFile=$1

if [ -z $inFile ]; then
    echo "Syntax: $0 filename"
    exit
fi

inFileSize=$(stat -f%z "$inFile")

tempFile=$inFile.tempFile.png
#echo "Shrinking $inFile ..."

if [ ! -e $inFile ]; then
    echo "$inFile dont exist"
    exit
fi

pngcrush -s -brute -rem alla $inFile $tempFile 2> /dev/null

if [ $? -ne 0 ]; then
    echo "Error 1 occured while processing $inFile"
    exit
fi

optipng -o7 $tempFile 2> /dev/null

if [ $? -ne 0 ]; then
    echo "Error 2 occured while processing $inFile"
    exit
fi

outFileSize=$(stat -f%z "$tempFile")

diffSize=$(($inFileSize - $outFileSize))

if [ $diffSize -lt 1 ]; then
    # echo "Error: resulting $inFile is not smaller than the original ($diffSize)"
    rm $tempFile
    exit
fi

echo "$inFile was $inFileSize bytes, shrank to $outFileSize. saved $diffSize bytes"

rm $inFile
mv $tempFile $inFile
