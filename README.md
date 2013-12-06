raptr
=====

A simple, shell-based APT Repository builder and updater.

Example usage:

raptr init \
  --archives "production staging testing" \
  --sections "public private"

raptr add \
  --section "public" \
  --package "my-package-name" \
  --dir "/location/to/some/package-dir"

raptr update \
  --gpg="mykey@domain.com" # optional key
  --store-key # used to put the store the key *in* the repository for convenience


