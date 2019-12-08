import request from '@/utils/request'

export function requestList(query) {
  return request({
    url: '/userSet/list',
    method: 'get',
    params: query
  })
}

export function requestDetail(id) {
  return request({
    url: '/userSet/detail',
    method: 'get',
    params: { id }
  })
}

export function requestUpdate(data) {
  return request({
    url: '/userSet/update',
    method: 'post',
    data
  })
}

export function requestCreate(data) {
  return request({
    url: '/userSet/create',
    method: 'post',
    data
  })
}

export function requestDelete(data) {
  return request({
    url: '/userSet/delete',
    method: 'post',
    data
  })
}

export function requestAdminsRoleIDList(adminsid) {
  return request({
    url: '/userSet/adminsroleidlist',
    method: 'get',
    params: { adminsid }
  })
}

export function requestSetRole(adminsid, data) {
  return request({
    url: '/userSet/setrole',
    method: 'post',
    params: { adminsid },
    data
  })
}

