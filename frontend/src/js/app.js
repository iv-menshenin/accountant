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
          if (parameters.url !== '/auth/refresh') {
            parameters.headers['X-Auth-Token'] = window.localStorage.getItem('devalio_token');
          }
        },
        statusCode: {
          401: async (xhr) => {
            try {
              const { method, url } = xhr.requestParameters;
              const { data } = await this.request.postJSON('/auth/refresh', { token: window.localStorage.getItem('devalio_refresh') });
              window.localStorage.setItem('devalio_token', data.data.jwt_token);
              window.localStorage.setItem('devalio_refresh', data.data.refresh_token);
              xhr.open(method, url);
              xhr.setRequestHeader('X-Auth-Token', data.data.jwt_token);
              xhr.send();
            } catch (error) {
              console.error('401', error);
            }
          },
        }
      });
    },
  },
});
