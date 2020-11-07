#!/bin/bash
cd frontend
yarn

if [ "$1" == "--dev" ]; then
  yarn dev
else
  yarn build;
  yarn start;
fi
