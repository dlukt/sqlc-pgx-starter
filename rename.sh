#!/usr/bin/env bash
#
# Rename this starter to your own project after cloning.
#
#   ./rename.sh <new-module-path> [new-dir-name]
#
# Examples:
#   ./rename.sh github.com/acme/billing
#   ./rename.sh github.com/acme/billing billing-svc
#
# It rewrites the module path (and the old directory name) everywhere in the
# tree, edits go.mod, and optionally renames the project directory.
#
set -euo pipefail

OLD_MODULE="github.com/dlukt/sqlc-pgx-starter"
OLD_NAME="sqlc-pgx-starter"

usage() {
    cat <<EOF
Usage: $(basename "$0") <new-module-path> [new-dir-name]

Replaces ${OLD_MODULE} with your own module path across the repository.

Examples:
  $(basename "$0") github.com/acme/billing
  $(basename "$0") github.com/acme/billing billing-svc
EOF
}

if [[ $# -lt 1 || "$1" == "-h" || "$1" == "--help" ]]; then
    usage
    exit 0
fi

NEW_MODULE="$1"
NEW_NAME="${2:-$(basename "$NEW_MODULE")}"

if [[ "$NEW_MODULE" == "$OLD_MODULE" ]]; then
    echo "Nothing to do — already using ${OLD_MODULE}"
    exit 0
fi

# Validate the module path looks sane (no spaces, no shell metacharacters).
if [[ ! "$NEW_MODULE" =~ ^[a-zA-Z0-9._/-]+$ ]]; then
    echo "error: invalid module path '${NEW_MODULE}'" >&2
    exit 1
fi
if [[ ! "$NEW_NAME" =~ ^[a-zA-Z0-9._-]+$ ]]; then
    echo "error: invalid directory name '${NEW_NAME}'" >&2
    exit 1
fi

echo "Module path : ${OLD_MODULE}  ->  ${NEW_MODULE}"
echo "Directory   : ${OLD_NAME}  ->  ${NEW_NAME}"
echo
read -r -p "Proceed? [y/N] " confirm
[[ "$confirm" =~ ^[yY]$ ]] || { echo "Aborted."; exit 1; }

# Files referencing the module path or directory name, excluding .git.
collect() {
    grep -rlI --exclude-dir=.git "$1" . 2>/dev/null || true
}

# 1. Rewrite the module path.
mapfile -t module_files < <(collect "$OLD_MODULE")
for f in "${module_files[@]}"; do
    sed -i "s|${OLD_MODULE}|${NEW_MODULE}|g" "$f"
    echo "  updated module ref: $f"
done

# 2. Rewrite the old directory name in text files (README, Makefile, compose, ...).
mapfile -t name_files < <(collect "$OLD_NAME")
for f in "${name_files[@]}"; do
    sed -i "s|${OLD_NAME}|${NEW_NAME}|g" "$f"
    echo "  updated name ref:   $f"
done

# 3. Make go.mod authoritative (idempotent with the sed above).
go mod edit -module "$NEW_MODULE"

# 4. Optionally rename the project directory itself.
current_dir="$(basename "$PWD")"
if [[ "$NEW_NAME" != "$OLD_NAME" && "$current_dir" == "$OLD_NAME" ]]; then
    parent="$(dirname "$PWD")"
    target="${parent}/${NEW_NAME}"
    cd "$parent"
    mv "$OLD_NAME" "$NEW_NAME"
    echo
    echo "Renamed project directory -> ${target}"
    echo "Run the following to continue:"
    echo "  cd \"${target}\""
fi

echo
echo "Done. Next steps:"
echo "  go mod tidy"
echo "  make check        # fmt + vet + build"
echo "  git add -A && git commit -m \"chore: rename to ${NEW_MODULE}\""
