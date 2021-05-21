# gnomeBackground

Similar to Variety, but very primitive in comparison.

`gnomeBackground` is a desktop background changer, specifically designed for the Gnome Desktop.

> Configurable using a config file `gnomeBackground.json` that allows you to set the file system directory/directories for images. As well as the screen size, the delay between background changes, and misc stuff.

```json
{
  "screenSize": "1920x1080",
  "fsPath": [
    "/usr/share/backgrounds/*.svg",
    "/usr/share/backgrounds/*.jpg",
    "/usr/share/backgrounds/*.png",
    "/usr/share/backgrounds/*.jpeg",
    "/usr/share/backgrounds/pop/*.svg",
    "/usr/share/backgrounds/pop/*.jpg",
    "/usr/share/backgrounds/pop/*.png",
    "/usr/share/backgrounds/pop/*.jpeg"
  ],
  "delay": 90,
  "dateStamp": {
    "display": true,
    "font": "Bitstream-Charter",
    "fontSize": 67.0,
    "color": "#00000099",
    "backgroundColor": "#FFFFFF25",
    "position": "                                                    %v",
    "format": "Monday  01/02/2006    03:04 "
  }
}
```

#### ToDo's that might be fun when I have time:

1. I may want to **add http get** logic to retrieve images online.
2. Another nice feature to add would be to automatically **reload the images array** whenever changes occur in the config file or any of the directories. Using **fsNotify** to watch the config file, and the directory/directories referenced in the config files for changes.
