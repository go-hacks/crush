# crush
64Bit file hashing algorithm 

For when security isn't the main concern and we just want to verify that our files were copied properly. 

Approximately 3X+ the speed of md5sum. 

Benchmarked on Artix Linux(x64) inside a ramdisk.

Build with
```
chmod +x && ./build
```

Run with
```
./crush fileName
```

Note: ./build dumps compilation information so you can see what
escapes to heap. Only values outside the core hasher should. 
