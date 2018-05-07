import fetch from 'dva/fetch';
import { notification } from 'antd';

function checkStatus(response) {
  if (response.status >= 200 && response.status < 300) {
    return response;
  }
  return response.json().then((body) => {
    const { message } = body;
    const error = new Error(message);
    error.response = response;
    throw error;
  });
}

function baseURL() {
  let { pathname } = window.location;
  if (pathname.length > 1) {
    if (pathname[pathname.length - 1] !== '/') {
      pathname += '/';
    }
  }
  return `${pathname}api`;
}

/**
 * Requests a URL, returning a promise.
 *
 * @param  {string} url       The URL we want to request
 * @param  {object} [options] The options we want to pass to "fetch"
 * @return {object}           An object containing either "data" or "err"
 */
export default function request(url, options) {
  const base = baseURL();

  const defaultOptions = {
    credentials: 'include',
  };
  const newOptions = {
    ...defaultOptions,
    ...options,
  };

  if (newOptions.method === 'POST' || newOptions.method === 'PUT') {
    newOptions.headers = {
      Accept: 'application/json',
      'Content-Type': 'application/json; charset=utf-8',
      ...newOptions.headers,
    };
    newOptions.body = JSON.stringify(newOptions.body);
  }

  return fetch(`${base}${url}`, newOptions)
    .then(checkStatus)
    .then(response => response.json())
    .catch((error) => {
      if ('stack' in error && 'message' in error) {
        notification.error({
          message: `请求错误: ${url}`,
          description: error.message,
        });
      }
      return error;
    });
}
