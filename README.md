# container-from-scratch
make container from scratch using Golang

## 3 main concepts:
1. Namespaces
  - PID: 
  ```
  Process isolation, so every PID namespace can have their PIDs start from 1 or 2
  and can have overlap. They don't have visiblity into other PID namespaces
  ```
  
  - MNT:
  ```
  Mount namespace gives processes in the namespace their own mount table.
  Mounting/unmounting doesn't affect others.
  ```
  
  - NET:
  ```
  Provides network isolation.
  Network namespaces can bt connected to each other, and can talk with routing.
  ```
  
  - UTS:
  ```
  Isolation for hostname and domain name
  ```
  
  - IPC:
  ```
  Inter Process Communication isolation
  ```
  
  - USER:
  ```
  It isolates container UIDs from host UIDs/GIDs, so root user ID 0, will be mapped to a random user ID on the host.
  This makes the namespace/container with root ID think it has root access (it actually does, but only limited to container respurces), but it doesn't have that on the host
  ```
  
2. cgroups 
3. Layered Filesystem

