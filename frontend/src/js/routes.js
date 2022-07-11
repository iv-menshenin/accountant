export default [
  {
    path: '/login/',
    name: 'login',
    asyncComponent: () => import('@/pages/login.f7'),
    beforeEnter: checkAuth,
    options: {
      auth: false,
    },
  },
  {
    path: '/accounts/',
    name: 'accounts',
    asyncComponent: () => import('@/pages/accounts.f7'),
    beforeEnter: checkAuth,
    options: {
      auth: true,
    },
    routes: [
      {
        path: 'new/',
        name: 'new_account',
        popup: {
          asyncComponent: () => import('@/pages/account.f7'),
        },
        beforeEnter: checkAuth,
        options: {
          auth: true,
        },
      },
      {
        path: ':accountId/',
        name: 'account',
        popup: {
          asyncComponent: () => import('@/pages/account.f7'),
        },
        beforeEnter: checkAuth,
        options: {
          auth: true,
        },
      },
    ],
  },
  {
    path: '/persons/',
    name: 'persons',
    asyncComponent: () => import('@/pages/persons.f7'),
    beforeEnter: checkAuth,
    options: {
      auth: true,
    },
  },
  {
    path: '/objects/',
    name: 'objects',
    asyncComponent: () => import('@/pages/objects.f7'),
    beforeEnter: checkAuth,
    options: {
      auth: true,
    },
  },
  {
    path: '/targets/',
    name: 'targets',
    asyncComponent: () => import('@/pages/targets.f7'),
    beforeEnter: checkAuth,
    options: {
      auth: true,
    },
  },
  {
    path: '(.*)',
    redirect: function ({ resolve }) {
      resolve(isAuth() ? '/accounts/' : '/login/');
    },
  },
];

function checkAuth({ to, resolve }) {
  if (to.route.options.auth ^ isAuth()) {
    this.navigate({ path: isAuth() ? '/accounts/' : '/login/' });
  } else {
    resolve(to);
  }
}

function isAuth() {
  return Boolean(window.localStorage.getItem('devalio_token'));
}
