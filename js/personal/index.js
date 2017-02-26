import React from 'react';
import { render } from 'react-dom';
import App from 'js/personal/app';

import Store from './store';
import API from './api';

const api = new API(window.ENV.API_BASE_URL);
const store = new Store(api);

render(
    <App store={ store } />,
    document.getElementById('react-doc')
);
