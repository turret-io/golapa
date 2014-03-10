golapa
------

**Instructions**

1. Setup a new AppEngine project.

2. Edit app.yaml at the base of this project and change `<PROJECT-ID-HERE>` to your AppEngine Project ID.

3. Edit launch/mail.go and replace `<APPENGINE-EMAIL-SENDER>` with the email address of the owner of your AppEngine project. Then, replace `<RECIPIENT>` with the email address that should receive new signup details.

4. Download the Go AppEngine SDK from https://developers.google.com/appengine/downloads.

5. Download and install Sass to compile the CSS file: `sudo gem install sass`. If you're on Windows or another OS, follow the instructions here: http://sass-lang.com/install

6. After extracting the archive, add the directory your `PATH` or you'll need to use the full path to the `goapp` program.

7. Compile the CSS file: `cd css; sass ../sass/launch.sass launch.css`. 

8. To test the app locally, enter the base directory of this project and run `goapp serve .`

9. Browse to `http://localhost:8080`

10. Tweak the templates in `templates/` as desired, update CSS in `sass/launch.sass` (rebuild via the sass instructions above), and add any images to `images/` 

11. Deploy to AppEngine. Enter the base directory of this project and run `goapp deploy .`
