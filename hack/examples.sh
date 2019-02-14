#! /bin/bash
examples=(examples/*.sh)

clear
"${examples[0]}"

for example in "${examples[@]:1}"; do
  echo ""
  read -p "Press ENTER to continue"
  echo ""
  echo "----------------------------------------------------------------------------------------------------"
  echo ""

  # clear the
  for i in "$(seq 2 "$(tput lines)")"; do
    echo "$i"
  done

  clear

  "$example"
done
