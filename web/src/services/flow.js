import { stringify } from 'qs';
import request from '../utils/request';

function isError(v) {
  return Object.prototype.toString.call(v) === '[object Error]';
}

export async function queryPage(params) {
  return request(`/flow/page?${stringify(params)}`).then((response) => {
    if (isError(response)) {
      return {};
    }
    return response;
  });
}

export async function get(params) {
  return request(`/flow/${params.record_id}`).then((response) => {
    if (isError(response)) {
      return {};
    }
    return response;
  });
}

export async function save(params) {
  return request('/flow', {
    method: 'POST',
    body: params,
  });
}

export async function del(params) {
  return request(`/flow/${params.record_id}`, {
    method: 'DELETE',
  });
}
