# Popup Manager #

Popup Manager is an open-source tool designed for creating and managing popups on your web pages without writing code. Create and manage informative popups for your web pages from one central place.

### Features ###

* Time based / event based informative popups
* Super performant, developed in Go
* No database required; it uses SQLite as an embedded database
* Support GeoIP2 for a geo location
* Simple analytics

### How do I get set up? ###

Clone the main branch, compile with the latest Go version, set the environment variables in the .env file, and you are good to go.

Because SQLite is used as an embedded database engine, which is a C library, if you want to use the Popup Manager backend on another system, you have to compile it using CGo with the following command:

 CGO_ENABLED=1 CC=musl-gcc go build --ldflags '-linkmode=external -extldflags=-static'

After the initial start, the migration will be automatically executed, and the SQLite database will be created in the same folder as the binary file. The Popup Manager backend requires the assets folder to be in the same directory as the binary file. Inside the assets folder, the GeoIP2 base can be saved, and the upload folder needs to be created where images for the popups will be saved. If the assets and upload folders don't exist, they will be automatically created.

### Website integration ###
Create Website Token and copy the following code line into your Website body.

<script type="text/javascript" src="HOSTNAME_API/api/v1/js/YOUR-WEB-SITE-TOKEN-UUID"></script>

### Contribution guidelines ###

* Writing tests
* Pull requests
* Code review
* Other guidelines

### Who do I talk to? ###

* Repo owner or admin
* Other community or team contact

### Is there a Demo? ###

* Demo can be found on the following url: [https://demo.popupeasy.com](https://demo.popupeasy.com) and can be tested with the following credentials: demo@demo.com / demo123#


### Is there a Cloud-Hosted version? ###

* Yes. The Cloud-Hosted version can be found on the following url: [https://popupeasy.com](https://popupeasy.com)  

### Is the Cloud-Hosted version different from the Self-Hosted version? ###

* Yes. The Cloud-Hosted version has many features not available in the Self-Hosted version. A full list of these features can be found on the official website: [https://popupeasy.com](https://popupeasy.com)
