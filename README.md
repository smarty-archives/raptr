raptr
==========

A simple, shell-based APT Repository builder and updater. In other words, [R]epository builder and [APT] update[R]: R.A.P.T.R.


####Example Usage


##### Initialization

```
raptr init --archives "production staging testing" --sections "public private" --cpus "amd64 i386 source" 
```

##### Adding Packages

```
raptr add --target "*/public" --package "/location/to/some/package-dir/package.dsc" # or deb
```

##### Updating Indexes (after adding packages)

```
raptr update --gpg="mykey@domain.com" --export-public-key
```
