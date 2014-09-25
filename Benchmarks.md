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
cs@fram:~/go/src/github.com/cornelius/gum> time sudo ./gum

real    0m7.973s
user    0m7.130s
sys     0m1.639s
cs@fram:~/go/src/github.com/cornelius/gum> time sudo ./gum

real    0m7.792s
user    0m6.956s
sys     0m1.647s
```

## Desktop

Inspecting a remote machine

### Current inspector

### Go inspector

