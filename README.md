raptr
=====

A simple, shell-based APT Repository builder and updater.  
In other words, Repository builder and APT updateR: RAPTR.  .

Example usage:

raptr init \  
  --archives "production staging testing" \  
  --sections "public private" \  
  --cpus "amd64 i386 source"  
  
raptr add \  
  --target "*/public" \  
  --package "/location/to/some/package-dir/package.dsc" # or deb    
  
raptr update \  
  --gpg="mykey@domain.com" \  
  --export-public-key  



