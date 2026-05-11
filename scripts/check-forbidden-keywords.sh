#!/usr/bin/env bash
set -euo pipefail

if [[ -z "${BILLTAP_FORBIDDEN_KEYWORDS:-}" ]]; then
  echo "::error title=Forbidden keyword scan not configured::BILLTAP_FORBIDDEN_KEYWORDS is empty."
  exit 1
fi

keywords_file="$(mktemp)"
trap 'rm -f "$keywords_file"' EXIT

printf '%s\n' "$BILLTAP_FORBIDDEN_KEYWORDS" \
  | tr ',' '\n' \
  | sed 's/^[[:space:]]*//;s/[[:space:]]*$//' \
  | awk 'length > 0 { print }' > "$keywords_file"

if [[ ! -s "$keywords_file" ]]; then
  echo "::error title=Forbidden keyword scan not configured::BILLTAP_FORBIDDEN_KEYWORDS did not contain any non-empty keyword."
  exit 1
fi

found=0

scan_file_for_keyword() {
  local file="$1"
  local keyword="$2"

  awk -v kw="$keyword" '
    BEGIN {
      kw = tolower(kw)
      matched = 0
    }
    index(tolower($0), kw) > 0 {
      printf "%s:%d\n", FILENAME, FNR
      matched = 1
    }
    END {
      exit matched ? 1 : 0
    }
  ' "$file"
}

while IFS= read -r -d '' file; do
  [[ -f "$file" ]] || continue
  if ! grep -Iq . "$file"; then
    continue
  fi

  while IFS= read -r keyword; do
    matches="$(scan_file_for_keyword "$file" "$keyword" || true)"
    if [[ -z "$matches" ]]; then
      continue
    fi

    found=1
    while IFS= read -r match; do
      [[ -n "$match" ]] || continue
      match_file="${match%:*}"
      match_line="${match##*:}"
      echo "::error file=${match_file},line=${match_line}::Forbidden keyword matched. Remove or generalize this repository content."
    done <<< "$matches"
  done < "$keywords_file"
done < <(git ls-files -z -- \
  ':!:dist/**' \
  ':!:node_modules/**' \
  ':!:.private/**' \
  ':!:*.log')

if [[ "$found" -ne 0 ]]; then
  exit 1
fi

echo "Forbidden keyword scan passed."
