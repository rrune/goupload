## GoUpload

Upload and host files and text. Comes with file, user and permission management. Files and texts can be either public or only visible to users. Undocumented personal projected, use [file browser](https://github.com/filebrowser/filebrowser) instead.

<img src="screenshots/screen1.png" alt="Upload a file" width="800">
<img src="screenshots/screen2.png" alt="Paste a text" width="800">
<img src="screenshots/screen3.png" alt="Dashboard to manage everything" width="800">

# Host it
- copy `config/config.yml.sample` to `config/config.yml`
- configure the config
- use `config/goupload.sql` to create the database
- configure `docker-compose.yml`
- `docker compose up`
