# `web`

This directory contains the source files of the Angular web application of the shinpuru web interface.

If you want to compile the web app or want to work with it, you need to download the required packages first.
```
npm install
```

You can start the Angular development server as following.
```
ng serve --port 8081
```

With the following command, you can build the static files of the web app which are located in the `dist` directory afterwards.
```
ng build --configuration production
```