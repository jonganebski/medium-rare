<h1 align="center">🥩 Medium-rare</h1>
<p align="center">Cloning <a href="https://medium.com/">Medium</a>. WYSIWYG editor is relying on editorjs library.</p>
<div align="center">
    <img src="https://img.shields.io/github/languages/count/jonganebski/medium-rare?style=flat"/>
    <img src="https://img.shields.io/github/languages/top/jonganebski/medium-rare?style=flat"/>
    <img src="https://img.shields.io/github/languages/code-size/jonganebski/medium-rare?style=flat"/>
    <img src="https://img.shields.io/github/last-commit/jonganebski/medium-rare?style=flat"/>
</div>

---

## Structure

|     Dir    |          Description         |
|:-----------|:-----------------------------|
|  .vscode   |  Vscode settings             |
|  assets    |  Sass & Typescript files     |
|  aws       |  Connect to aws s3           |
|  config    |  Read .env file              |
|  database  |  Connect to MongoDB          |
|  helper    |  Various functions           |
|  image     |  Images frontend consumes    |
|  middleware|  Fiber middlewares           |
|  model     |  Data models definition      |
|  package   |          |
|  router    |  Router                      |
|  routes    |  API end points              |
|  static    |  Generated static files      |
|  util      |  Cookies and password        |
|  views     |  Pug files                   |


## Run in local

```
// .env

APP_ENV=DEV

PORT=4000

DB_NAME=
COLLECTION_USER=
COLLECTION_STORY=
COLLECTION_COMMENT=
MONGO_URI_DEV=
MONGO_URI_PROD=

EDITORS=example@example.com example2@example.com

AWS_REGION=
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
BUCKET_NAME=
```

```console
$ npm install 

// frontend 
$ npm run start:dev

// backend
$ air
```
