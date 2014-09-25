# Benchmarks

Inspection of unmanaged files without extracting files

## Laptop

Local inspection running the inspector on the inspected machine

### Current inspector

```
:~/go/src/github.com/cornelius/gum> time mxy inspect w.x.y.z. -s unmanaged-files
Inspecting w.x.y.z for unmanaged-files...
Inspecting unmanaged-files...
 -> Found 3136 unmanaged files and trees.

real    8m55.273s
user    0m53.353s
sys     0m25.035s
```

### Go inspector

```
:~/go/src/github.com/cornelius/gum> time sudo ./gum

real    0m8.127s
user    0m7.182s
sys     0m1.644s
```

## Desktop

Inspecting a remote machine

### Current inspector

```
:~/go/src/github.com/cornelius/gum> time mxy inspect endurance.suse.de -s unmanaged-files
Inspecting endurance.suse.de for unmanaged-files...
Inspecting unmanaged-files...
 -> Found 2976 unmanaged files and trees.

real    13m23.761s
user    0m54.770s
sys     0m25.071s
```

### Go inspector

```
:~/go/src/github.com/cornelius/gum> time scp gum root@x.y.z:
gum                                                                                                       100% 2510KB   2.5MB/s   00:01

real    0m2.934s
user    0m0.024s
sys     0m0.021s
:~/go/src/github.com/cornelius/gum> time ssh root@x.y.z ./gum

real    0m7.726s
user    0m0.018s
sys     0m0.004s
:~/go/src/github.com/cornelius/gum> time scp root@x.y.z:UNMANAGED_FILES .
UNMANAGED_FILES                                                                                           100%  248KB 247.7KB/s   00:00

real    0m0.764s
user    0m0.018s
sys     0m0.009s
```
