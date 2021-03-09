# filepool

Sync files between server and client.

Files are hashed and encrypted for performance and privacy purposes.

We currently support stream AES, SHA256. Local files

# Mode

- `upload` upload files from client to server

- `download` download files from server to client

- `clean` delete files exist at server but client

# TODO

- Extend to distributed file system. Each file is stored in at least 3 machines at all time.
