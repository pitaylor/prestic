#!/usr/bin/env bash

mkdir -p tmp

echo "I'm mock restic! I was run with this command:"

(
  IFS=$'\n' read -r -d '' -a vars < <((env | grep ^RESTIC | sort) && printf '\0')

  # shellcheck disable=SC2199
  if [[ " ${@} " =~ " --stdin " ]] || [[ " ${@} " =~ " - " ]]; then
    stdin=$(timeout 2 cat || printf '<n/a>')
    vars+=("STDIN=${stdin}")
  fi

  # print env vars and stdin
  if ((${#vars[@]})); then
    printf "%q " "${vars[@]}"
  fi

  # print command
  printf "%s" "$(basename "${0}")"

  # print arguments
  if ((${#@})); then
    printf " %q" "${@}"
  fi

  printf "\n"
) | tee -a tmp/commands.log

# shellcheck disable=SC2199
if [[ " ${@} " =~ " --fail " ]]; then
  echo "simulated failure"
  exit 1
fi

echo "snapshot bedabb1e saved"
