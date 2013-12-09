raptr
=====

A simple, shell-based APT Repository builder and updater.

Example usage:

raptr init
  --archives "production staging testing" 
  --sections "public private"

raptr add 
  --section "public" 
  --package "my-package-name" 
  --directory "/location/to/some/package-dir"

raptr update 
  --gpg="mykey@domain.com"
  --export-public-key


