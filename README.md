golapa
------

**Instructions**

1. Setup a new AppEngine project.

2. Edit app.yaml at the base of this project and change `<PROJECT-ID-HERE>` to your AppEngine Project ID.

3. Edit launch/email.go and replace `<APPENGINE-EMAIL-SENDER>` with the email address of the owner of your AppEngine project. Then, replace `<RECIPIENT>` with the email address that should receive new signup details.

4. Download and extract the Go AppEngine SDK from https://developers.google.com/appengine/downloads.

5. Download and install Sass to compile the CSS file: `sudo gem install sass`. If you're on Windows or another OS, follow the instructions here: http://sass-lang.com/install

6. Compile the CSS file: `cd css; sass ../sass/launch.sass launch.css`. 

7. To test the app locally, enter the base directory of this project and run `/path/to/appengine/sdk/goapp serve .`

8. Browse to `http://localhost:8080`

9. Tweak the templates in `templates/` as desired, update CSS in `sass/launch.sass` (rebuild via the sass instructions above), and add any images to `images/` 

10. Deploy to AppEngine. Enter the base directory of this project and run `/path/to/appengine/sdk/goapp deploy .`

**Turret.IO Support**

Golapa now supports Turret.IO as an alternative to sending an email when a user signs up.

**By default, users that sign up will have an attribute *signedup* set to 1 making it easy to create a target for those users**

1. Set your `GOPATH` to the top-level of the project
```
> export GOPATH=~/golapa
```

2. Provide your API key and API secret (see (https://tws.turret.io/apidoc)) in `launch/launch.go`
```
const api_key = string("YOUR TURRET.IO API KEY")
const api_secret = string("YOUR TURRET.IO API SECRET")
```  

3. Edit `launch/launch.go` and change `http.HandleFunc("/worker", EmailSubmitter)` to `http.HandleFunc("/worker", TurretIOSubmitter)` 

Re-deploy your app and check for errors. Be sure to test that users are successfully being added to whatever targets you create with Turret.IO

