# container-from-scratch
Make container from scratch using Golang


## Where do containers come from? :bird: & :honeybee:
![how is it made](https://raw.githubusercontent.com/gunjan5/container-from-scratch/master/container.png)

## How to run this thing from CLI(tested on Ubuntu):
- Download and build for your platform:
```bash
$ git clone https://github.com/gunjan5/container-from-scratch.git
$ cd container-from-scratch
$ make get # to download dependancies or make restore to restore Godeps saved dependancies
$ make build
```
- `run` command format is: `sudo ./cfs <action_command> <OS_image> <command_to_run_inside_the_container>`
- Supported OS_images provided with this repo: `BusyBox`, `SlitazOS`, `TinyCore` 
```bash
$ whoami
gunjan
$ pwd
/home/gunjan/go/src/github.com/gunjan5/container-from-scratch
$ sudo ./cfs run SlitazOS sh
/ # whoami
root
/ # pwd
/
/ # 
```
```bash
$ sudo ./cfs run TinyCore ls
[./cfs run TinyCore ls]
[/proc/self/exe newroot TinyCore ls]
[TinyCore ls]
bin      core.gz  dev      etc      init     lib      linuxrc  opt      root     sbin     tmp      usr      var
```

## REST server:
- Start the CFS server `sudo ./cfs server`
- REST calls: 
`(GET) 127.0.0.1:1337/containers`
`(GET) 127.0.0.1:1337/history`
`(POST) 127.0.0.1:1337/run`

- JSON structure examples:

..- Run a new container
  ```
  {
    "state": "run",
    "image": "BusyBox",
    "command": "pwd"
  }
  ```

..* Stop a running container with it's Container ID
```json
  {
    "id": "e7887770-da8e-43db-9ca1-69526d144d7c",
    "state": "stop"
  }
```
    
- CURL call examples: 

`curl -H "Content-Type: application/json" -X POST -d '{"state":"run","image":"TinyCore","command":"ls"}' http://localhost:1337/run`
`curl -H "Content-Type: application/json" -X POST -d '{"id":"d78347b9-d7c1-4e22-b2fc-782c8111cfcb","state":"stop"}' http://localhost:1337/run`



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

```
  ___  ____  ___       ___  _____  _  _  ____   __    ____  _  _  ____  ____    ____  ____  _____  __  __ 
 / __)( ___)/ __)()   / __)(  _  )( \( )(_  _) /__\  (_  _)( \( )( ___)(  _ \  ( ___)(  _ \(  _  )(  \/  )
( (__  )__) \__ \    ( (__  )(_)(  )  (   )(  /(__)\  _)(_  )  (  )__)  )   /   )__)  )   / )(_)(  )    ( 
 \___)(__)  (___/()   \___)(_____)(_)\_) (__)(__)(__)(____)(_)\_)(____)(_)\_)  (__)  (_)\_)(_____)(_/\/\_)
 ___   ___  ____    __   ____  ___  _   _ 
/ __) / __)(  _ \  /__\ (_  _)/ __)( )_( )
\__ \( (__  )   / /(__)\  )( ( (__  ) _ ( 
(___/ \___)(_)\_)(__)(__)(__) \___)(_) (_)
```

### Credit:
- This is based on @doctor_julz's blog post about containers
- Container image from @onsijoe's talk at CF Summit



