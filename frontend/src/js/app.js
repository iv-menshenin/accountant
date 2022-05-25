import $ from 'dom7';
import Framework7 from './framework7-custom.js';

// Import F7 Styles
import '../css/framework7-custom.less';

// Import Icons and App Custom Styles
import '../css/icons.css';
import '../css/app.less';


// Import Routes
import routes from './routes.js';
// Import Store
import store from './store.js';

// Import main app component
import App from '../app.f7';


var app = new Framework7({
  name: 'Accountant', // App name
  theme: 'aurora', // Automatic theme detection
  el: '#app', // App root element
  component: App, // App main component

  // App store
  store: store,
  // App routes
  routes: routes,

  on: {
    init() {
      this.request.setup({
        beforeCreate: (parameters) => {
          parameters.url = 'https://victoria.devaliada.ru' + parameters.url;
          parameters.headers.authorization = window.localStorage.getItem('devalio_token');
        },
      });
    },
  },
});
