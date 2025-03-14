#!/bin/bash
echo "Enter countdown time in seconds:"
read seconds
while [ $seconds -gt 0 ]; do
  echo "$seconds seconds remaining..."
  sleep 1
  ((seconds--))
done
echo "Time's up!"
