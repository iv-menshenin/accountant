# Accountant

## Framework7 CLI Options

Framework7 app created with following options:

```
{
  "cwd": "/Users/sergeivoroshilov/Development/accountant/frontend",
  "type": [
    "web"
  ],
  "name": "Accountant",
  "framework": "core",
  "template": "blank",
  "bundler": "vite",
  "cssPreProcessor": "less",
  "theming": {
    "customColor": true,
    "color": "#3e6b83",
    "darkTheme": false,
    "iconFonts": true,
    "fillBars": true
  },
  "customBuild": true,
  "customBuildConfig": {
    "rtl": false,
    "darkTheme": false,
    "lightTheme": true,
    "themes": [
      "aurora"
    ],
    "components": [
      "calendar",
      "card",
      "checkbox",
      "contacts-list",
      "chip",
      "data-table",
      "dialog",
      "form",
      "grid",
      "infinite-scroll",
      "input",
      "login-screen",
      "panel",
      "popover",
      "popup",
      "preloader",
      "progressbar",
      "radio",
      "searchbar",
      "skeleton",
      "toast",
      "toggle",
      "tooltip",
      "typography",
      "virtual-list"
    ]
  }
}
```

## Install Dependencies

First of all we need to install dependencies, run in terminal
```
npm install
```

## NPM Scripts

* 🔥 `start` - run development server
* 🔧 `dev` - run development server
* 🔧 `build` - build web app for production

## Vite

There is a [Vite](https://vitejs.dev) bundler setup. It compiles and bundles all "front-end" resources. You should work only with files located in `/src` folder. Vite config located in `vite.config.js`.
## Assets

Assets (icons, splash screens) source images located in `assets-src` folder. To generate your own icons and splash screen images, you will need to replace all assets in this directory with your own images (pay attention to image size and format), and run the following command in the project directory:

```
framework7 assets
```

Or launch UI where you will be able to change icons and splash screens:

```
framework7 assets --ui
```



## Documentation & Resources

* [Framework7 Core Documentation](https://framework7.io/docs/)



* [Framework7 Icons Reference](https://framework7.io/icons/)
* [Community Forum](https://forum.framework7.io)

## Support Framework7

Love Framework7? Support project by donating or pledging on:
- Patreon: https://patreon.com/framework7
- OpenCollective: https://opencollective.com/framework7