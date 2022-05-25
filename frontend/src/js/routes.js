export default [
  {
    path: '/login/',
    asyncComponent: () => import('@/pages/login.f7'),
    beforeEnter: checkAuth,
    options: {
      auth: false,
    },
  },
  {
    path: '/accounts/',
    asyncComponent: () => import('@/pages/accounts.f7'),
    beforeEnter: checkAuth,
    options: {
      auth: true,
    },
  },
  {
    path: '/persons/',
    asyncComponent: () => import('@/pages/persons.f7'),
    beforeEnter: checkAuth,
    options: {
      auth: true,
    },
  },
  {
    path: '/objects/',
    asyncComponent: () => import('@/pages/objects.f7'),
    beforeEnter: checkAuth,
    options: {
      auth: true,
    },
  },
  {
    path: '/targets/',
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
