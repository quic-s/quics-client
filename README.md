# quics-client

quics-client is a client for the QUIC-S. It is **continuous file synchornization tool** based on the QUIC protocol. 

 **NOTICE**  If you want to use this tool, you should use the server of QUIC-S. You can find the server in [here](https://github.com/quic-s/quics.git) 

[Features](#features) | [Getting Started](#getting-started) | [How to use](#how-to-use) | [Documentation](#documentation) | [Contribute](#contribute)

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

> For more detail logic and implementation, please check [QUIC-S Docs](https://github.com/quic-s/quics/tree/main/docs)
## Getting Started

### 1. Docker
    
```
docker run -it -d  -v /path/to/your/dir:/dirs --name quics-client -p 6120:6120 -p 6121:6121/udp  quics/quics-client
```

### 2. Local Install

- 1. Download the latest release from the [releases page](https://github.com/quic-s/quics-client/releases)
- 2. Unpack the archive.
- 3. Run `mv ./qic /usr/local/bin/qic`



### 3. Build from source

- 1. Install Go 1.21 or later.
- 2. Clone this repository.
    ```
    https://github.com/quic-s/quics-client.git
    ```
- 3. Run the command in the root of the repository.
    ```
    go build -o qic ./cmd
    ```



## How to use

> **Check the video for how to use** (click)
>
> [![quics-video](https://img.youtube.com/vi/B5u2qiNFiV8/0.jpg)](https://youtu.be/B5u2qiNFiV8)

| Tag | Command | Options |     Description     | HTTP Method | Endpoint |
| --- | --- | --- | --- | --- | --- |
| controller | `qic start` | `--hport`  | start the program | X |  |
| controller | `qic` | `--help` | show the help message | X |  |
| config | `qic config show` |  | read .qis.env | GET | /api/v1/config/show |
| config | `qic config server` | `--host` `--port`  | main server config | POST | /api/v1/config/server |
| connect | `qic connect server` | `--host` `--port` `--password`  | connect to the server | POST | /api/v1/connect/server |
| connect | `qic connect root` |`--local` `--password` |   connect to the local root directory | POST | /api/v1/connect/root/local |
| connect | `qic connect root` | `--remote` `--password`  |   connect to the remote root directory | POST | /api/v1/connect/root/remote |
| connect | `qic connect list-remote` |  | get the list of remote root directory | GET | /api/v1/connect/list/remote |
| disconnect | `qic disconnect root` | `--path` | disconnect to the root directory | POST | /api/v1/disconnect/root |
| sync | `qic sync status` |y  `--path` | get the status of the root directory | POST | /api/v1/sync/status |
| sync | `qic sync rescan` |  | rescan the all of root directory | POST | /api/v1/sync/rescan |
| conflict | `qic conflict list` |   | get the list of the root directory | GET | /api/v1/conflict/list |
| conflict | `qic conflict choose` | `--path` `--candidate` | choose the file of the root directory | POST | /api/v1/conflict/choose |
| conflict | `qic conflict download` | `--path` | download the file of the root directory | POST | /api/v1/conflict/download |
| sharing | `qic share file` | `--path` `--cnt` | share the file of the root directory | POST | /api/v1/share/download |
| sharing | `qic share stop` | `--link` | stop the sharing of the file | POST | /api/v1/share/stop |
| history | `qic history rollback` | `--path` `--version`  | rollback the file to certain version | POST | /api/v1/history/rollback |
| history | qic` history show` | `--path` `--from-head` | show the history of the file | POST | /api/v1/history/show |
| history | `qic history download` | `--path`  `--version` | download the file of the root directory | POST | /api/v1/history/download |



## Documentation
For more detail logic and implementation, please check [QUIC-S Docs](https://github.com/quic-s/quics/tree/main/docs)
Also you can check [quics](https://github.com/quic-s/quics) for server side and [quics-protocol](https://github.com/quic-s/quics-protocol) for protocol.

## Contribute
QUIC-S is an open source project, and contributions of any kind are welcome and appreciated.

We also have a awesome plan to make QUIC-S better. Check [ROADMAP.md](https://github.com/quic-s/quics/blob/main/ROADMAP.md) will be helpful to understand our project's direction.

To contribute, please read [CONTRIBUTING.md](https://github.com/quic-s/quics/blob/main/CONTRIBUTING.md)
- To report bugs or request features, please use the issue tracker. Before you do so, make sure you are running the latest version, and please do a quick search to see if the issue has already been reported.

- For more discussion, please join the [quics discord](https://discord.gg/HRtY7pNZz2)




