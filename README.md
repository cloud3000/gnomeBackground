# gnomeBackground

Similar to Variety, but very primitive in comparison.

`gnomeBackground` is a desktop background changer, specifically designed for the Gnome Desktop.

> Configurable using a config file `gnomeBackground.json` that allows you to set the file system directory/directories for images. As well as the screen size, the delay between background changes, and misc stuff.

```json
{
  "screenSize": "5760x1080",
  "fsPath": [
    "/home/michael/Downloads/backgrounds/*.jpg",
    "/home/michael/Downloads/backgrounds/*.png",
    "/home/michael/Downloads/backgrounds/*.jpeg",
    "/home/michael/Downloads/backgrounds/*.webp"
  ],
  "delay": 60,
  "dateStamp": true,
  "fontSize": 96.0
}
```

#### ToDo's that might be fun when I have time:

1. **Image size.** Needs to be able to scale, crop, zoom in/out based on the size of the image and the size of the screen.
2. I may want to **add http get** logic to retrieve image online.
3. Another nice feature to add would be to automatically **reload the images array** whenever changes occur in the config or any of the directories. Using **fsNotify** to watch the config file, and the directory/directories referenced in the config files for changes.
