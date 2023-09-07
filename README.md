# quics-client

quics-client is a client for the QUIC-S. It is **continuous file synchornization tool** based on the QUIC protocol. 

 **NOTICE**  If you want to use this tool, you should use the server of QUIC-S. You can find the server in [here](https://github.com/quic-s/quics.git) 


## Features

#### 1. App Controller
Start, Reboot, Stop management of the app. And also, user can check the status of the app. All the management is done by the command line made by **cobra** library.

#### 2. App Config
By using **viper** library, environment variables is managed. Mostly it is used for the configuration of the Root Directory Path for sync. And also, it is used for the configuration of the server IP address and port number.

#### 3. Connection Manager
* Server Connection : User can connect with access token and IP address, then dial to server to send first message with user information. 
* Root Directory Connection : User can select the directory that user want to manage. Choose Local Root Directory or Remote Root Directory that any folder user wants. And then, user can start to sync.

#### 4. File Manager
By using **fsnotify** library and goroutine, user can monitor the file changes in the directory. And then, user can send the file changes to the server. 


#### 5. File Sync
For make sync in real time, the Sync Metadata is saved in the local db, the **BadgerDB**. This is for the comparison of the file changes. And also, user can check the file changes in the local db.

#### 6. Extension
**HTTP/3** is used for the each cobra command, except the controller command. And also, the API is used for the communication between the client and server. By using the HTTP/3 API, developers can easily integrate the GUI or other functions.

## Getting Started

### 1. Docker
    
```
    docker run -it --rm -v /path/to/your/dir:/data chromato99/quics-client
```

### 2. Local Install

- 1. Download the latest release from the [releases page]()
- 2. Unpack the archive.
- 3. Run `./qic` in the unpacked directory.



### 3. Build from source

- 1. Install Go 1.21 or later.
- 2. Clone this repository.
    ```
    https://github.com/quic-s/quics-client.git
    ```
- 3. Run the command in the root of the repository.
    ```
    go build ./cmd
    ```



## How to use

| Tag | Command | Options |     Description     | HTTP Method | Endpoint |
| --- | --- | --- | --- | --- | --- |
|controller	| `qic start`| `--hPort`| start client |    |	|
|controller	| `qic reboot`| | reboot client |    |	|
|controller	| `qic shutdown`| | shutdown client |    |	|
|controller	| `qic --help`| | show help |    |	|
|config	| `qic config show`| | read .qis.env |    |	|
|config	| `qic config server`|` --host {serverIp} --port  {port}`| set main server config |    |	|
|config	| `qic config root`|` --abspath {dirpath}  --name {dir-NN}`| set root dir config |    |	|
|config	| `qic config delete`|` --key {key}`| delete root dir config |    |	|
|connect	| `qic connect server`|` --password {password} --host {serverIp} --port {port}`| connect to server | POST | `/api/v1/connect/server` |
|disconnect	| `qic disconnect server`|` --password {ClientPassword}`| disconnect to server | POST | `/api/v1/disconnect/server`|
|disconnect	| `qic disconnect root`|` --root {root dir/NN} --password {pw}`| disconnect to root dir | POST | `/api/v1/disconnect/root`|
|connect	| `qic connect root`|` --local  {root-dir/NN} -- password {pw}`| connect to local root dir | POST | `/api/v1/connect/root/local`|
|connect	| `qic connect root`|` --remote {root dir/NN} --password {pw}`| connect to remote root dir | POST | `/api/v1/connect/root/remote`|
|connect	| `qic connect list-remote`| | get rootdir list | GET | `/api/v1/connect/list/remote`|
|sync	| `qic sync status`|` --pick {root dir/NN}`| get root dir status | POST | `/api/v1/status/root/`|
|sync	| `qic sync rescan`| | ask server to rescan all root dir | POST | `/api/v1/rescan`|



## Contribute

- To report bugs or request features, please use the issue tracker. Before you do so, make sure you are running the latest version, and please do a quick search to see if the issue has already been reported.

- For more discussion, please join the [quics discord](https://discord.gg/HRtY7pNZz2)

