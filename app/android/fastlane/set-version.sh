#!/bin/bash

version_code=$1
version_name=$2

sed -i "s/versionCode = .*/versionCode = ${version_code}/g" ../variables.gradle
sed -i "s/versionName = .*/versionName = \"${version_name}\"/g" ../variables.gradle
