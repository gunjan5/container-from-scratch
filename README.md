# container-from-scratch
Make container from scratch using Golang

## 3 main components of a container:
1. **Namespaces:**
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

2. **cgroups:**
  ```
  Control group  is a Linux kernel feature that limits, accounts for, and isolates
  the resource usage (CPU, memory, disk I/O, network, etc.) of a collection of processes.

  If Namespace is for isolation, Cgroup is for sharing resources fairly.
  ```

3. **Layered Filesystem:**
  ```
  Multilayer filesystem, so some base layers (e.g linux base layers) can be shared by multiple containers.
  Layers are usually read-only, and containers make copy if they have to modify the layer.
  ```


### Credit:
- This is based on @doctor_julz's blog post about containers

