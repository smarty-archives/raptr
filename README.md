raptr
==========

A simple, shell-based APT Repository builder and updater. In other words, [R]epository builder and [APT] update[R]: R.A.P.T.R.


####Example Usage


##### One-time Repository Initialization (akin to `git init`--it should be run once)

```
raptr init --archives "production staging testing" --sections "public private" --cpus "amd64 i386 source" 
```

##### Adding Binary Packages

```
raptr add --target "*/public" --package "/location/to/some/package-dir/package.deb"
```

##### Adding Source Packages

```
raptr add --target "*/public" --package "/location/to/some/package-dir/package.dsc"
```

##### Updating Indexes (after adding packages)

```
raptr update --gpg="mykey@domain.com" --export-public-key
```
