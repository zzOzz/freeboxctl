# Freeboxctl
freeboxctl is command line interface for managing freebox server remote downloading.

1. launch app
    ~~~
    freeboxctl get downloads
    ~~~
    
    ~~~
    INFO[0000] Authorize: Token="*************************" TrackID=0
    ~~~
2. validate app on the freebox display (use arrows)
3. set env (you can change freebox url to enable access from outside if you have open freebox interface to the internet)
    ~~~
    export GOFBX_TOKEN=*************************
    export GOFBX_URL="http://mafreebox.free.fr/"
    ~~~
4. let's play (download a file)
    ~~~
    freeboxctl add download http://www.ovh.net/files/1Mio.dat
    ~~~

### Completion

for zsh

~~~
source <(freeboxctl completion -z)
~~~

for bash (without remote files completion)

~~~
source <(freeboxctl completion -b)
~~~


### Add remote download

~~~
freeboxctl add download http://www.ovh.net/files/1Mio.dat
~~~
~~~
{"result":{"id":5119},"success":true}
~~~



### Get current downloads

~~~
freeboxctl get downloads | jq .
~~~
~~~
[
  {
    "rx_bytes": 1040000,
    "tx_bytes": 0,
    "download_dir": "L0ZyZWVib3gvZG93bmxvYWRz",
    "archive_password": "",
    "eta": 0,
    "status": "done",
    "io_priority": "normal",
    "type": "http",
    "piece_length": 0,
    "queue_pos": 1,
    "id": 5119,
    "info_hash": "",
    "created_ts": 1545136386,
    "stop_ratio": 0,
    "tx_rate": 0,
    "name": "1Mio.dat",
    "tx_pct": 10000,
    "rx_pct": 10000,
    "rx_rate": 0,
    "error": "none",
    "size": 1040000
  }
]
~~~

nice display with current status and percentage
~~~
freeboxctl get downloads|jq -r '.[]|.name + " " + ((.rx_bytes/.size*100|floor|tostring)+"%")+ " " + .status'
~~~
~~~
1Mio.dat 100% done
~~~

### Clean up downloads

~~~
freeboxctl delete downloads
~~~

### Downloading files locally

~~~
freeboxctl get files "//Freebox/downloads/1Mio.dat" -d > 1Mio.dat
~~~

### List remote files

just directory/file name
~~~
freeboxctl get files "/Freebox/"  | jq ".[]|.name"
~~~
~~~
"."
".."
".fbxgrabberd"
".fbxtimeshifting"
"downloads"
"Enregistrements"
"Musiques"
"Photos"
"Téléchargements"
"Vidéos"
".apdisk"
".DS_Store"
~~~

full version
~~~
freeboxctl get files "/Freebox/"  | jq .~
~~~

~~~
[
  {
    "path": "L0ZyZWVib3gv",
    "filecount": 0,
    "link": false,
    "modification": 1544558820,
    "foldercount": 0,
    "name": ".",
    "index": 0,
    "mimetype": "inode/directory",
    "hidden": true,
    "type": "dir",
    "size": 4096
  },
  {
    "path": "Lw==",
    "filecount": 0,
    "link": false,
    "modification": 1544558819,
    "foldercount": 0,
    "name": "..",
    "index": 1,
    "mimetype": "inode/directory",
    "hidden": true,
    "type": "dir",
    "size": 60
  },
  {
    "path": "L0ZyZWVib3gvLmZieGdyYWJiZXJk",
    "filecount": 0,
    "link": false,
    "modification": 1543425554,
    "foldercount": 0,
    "name": ".fbxgrabberd",
    "index": 2,
    "mimetype": "inode/directory",
    "hidden": true,
    "type": "dir",
    "size": 4096
  },
  {
    "path": "L0ZyZWVib3gvLmZieHRpbWVzaGlmdGluZw==",
    "filecount": 0,
    "link": false,
    "modification": 1538508715,
    "foldercount": 0,
    "name": ".fbxtimeshifting",
    "index": 3,
    "mimetype": "inode/directory",
    "hidden": true,
    "type": "dir",
    "size": 4096
  },
  {
    "path": "L0ZyZWVib3gvZG93bmxvYWRz",
    "filecount": 0,
    "link": false,
    "modification": 1545134955,
    "foldercount": 0,
    "name": "downloads",
    "index": 4,
    "mimetype": "inode/directory",
    "hidden": false,
    "type": "dir",
    "size": 28672
  },
  {
    "path": "L0ZyZWVib3gvRW5yZWdpc3RyZW1lbnRz",
    "filecount": 0,
    "link": false,
    "modification": 1544558820,
    "foldercount": 0,
    "name": "Enregistrements",
    "index": 5,
    "mimetype": "inode/directory",
    "hidden": false,
    "type": "dir",
    "size": 4096
  },
  {
    "path": "L0ZyZWVib3gvTXVzaXF1ZXM=",
    "filecount": 0,
    "link": false,
    "modification": 1544558820,
    "foldercount": 0,
    "name": "Musiques",
    "index": 6,
    "mimetype": "inode/directory",
    "hidden": false,
    "type": "dir",
    "size": 4096
  },
  {
    "path": "L0ZyZWVib3gvUGhvdG9z",
    "filecount": 0,
    "link": false,
    "modification": 1544558820,
    "foldercount": 0,
    "name": "Photos",
    "index": 7,
    "mimetype": "inode/directory",
    "hidden": false,
    "type": "dir",
    "size": 4096
  },
  {
    "path": "L0ZyZWVib3gvVMOpbMOpY2hhcmdlbWVudHM=",
    "filecount": 0,
    "link": false,
    "modification": 1544558820,
    "foldercount": 0,
    "name": "Téléchargements",
    "index": 8,
    "mimetype": "inode/directory",
    "hidden": false,
    "type": "dir",
    "size": 4096
  },
  {
    "path": "L0ZyZWVib3gvVmlkw6lvcw==",
    "filecount": 0,
    "link": false,
    "modification": 1544558820,
    "foldercount": 0,
    "name": "Vidéos",
    "index": 9,
    "mimetype": "inode/directory",
    "hidden": false,
    "type": "dir",
    "size": 4096
  },
  {
    "path": "L0ZyZWVib3gvLmFwZGlzaw==",
    "filecount": 0,
    "link": false,
    "modification": 1427452486,
    "foldercount": 0,
    "name": ".apdisk",
    "index": 10,
    "mimetype": "application/octet-stream",
    "hidden": true,
    "type": "file",
    "size": 293
  },
  {
    "path": "L0ZyZWVib3gvLkRTX1N0b3Jl",
    "filecount": 0,
    "link": false,
    "modification": 1544957510,
    "foldercount": 0,
    "name": ".DS_Store",
    "index": 11,
    "mimetype": "application/octet-stream",
    "hidden": true,
    "type": "file",
    "size": 10244
  }
]
~~~

## Extract archive

~~~
freeboxctl extract "/Freebox/downloads/1Mio.rar"| jq .
~~~
~~~
{
  "result": {
    "created_ts": 1545136925,
    "curr_bytes": 0,
    "curr_bytes_done": 0,
    "done_ts": 0,
    "duration": 0,
    "error": "none",
    "eta": 0,
    "from": "/Freebox/downloads/1Mio.rar",
    "id": 45,
    "nfiles": 0,
    "nfiles_done": 0,
    "progress": 0,
    "rate": 0,
    "started_ts": 1545136925,
    "state": "running",
    "to": "/Freebox/downloads",
    "total_bytes": 0,
    "total_bytes_done": 0,
    "type": "extract"
  },
  "success": true
}
~~~

## List FS tasks
~~~
freeboxctl get tasks| jq .
~~~
~~~
{
  "success": true,
  "result": [
    {
      "id": 43,
      "from": "/Freebox/downloads/1Mio.rar",
      "to": "/Freebox/downloads/1Mio.dat",
      "type": "extract",
      "state": "done"
    },
    {
      "id": 44,
      "from": "/Freebox/downloads/1Mio.rar",
      "to": "/Freebox/downloads/1Mio.dat",
      "type": "extract",
      "state": "done"
    },
    {
      "id": 45,
      "from": "/Freebox/downloads/1Mio.rar",
      "to": "/Freebox/downloads/1Mio.dat",
      "type": "extract",
      "state": "done"
    }
  ]
}
~~~

## Clean FS tasks
~~~
freeboxctl delete tasks
~~~