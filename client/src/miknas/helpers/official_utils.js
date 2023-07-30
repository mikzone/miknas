import axios from 'axios';
import Qs from 'qs';
import CryptoJS from 'crypto-js';
import {
  Notify,
  LocalStorage,
  Dialog,
  date,
  copyToClipboard,
  extend,
} from 'quasar';

const MsgTypeConverts = {
  success: 'positive',
  error: 'negative',
};

function sendTips(message, msgType = 'info') {
  // 支持msgType: positive, negative, warning, info, ongoing
  if (MsgTypeConverts[msgType]) msgType = MsgTypeConverts[msgType];
  Notify.create({
    // message: `<pre class="q-ma-none">${message}</pre>`,
    message: `${message}`,
    type: msgType,
    position: 'top',
    // html: true,
    progress: true,
  });
}

export function formatFloat(num) {
  return num.toLocaleString(undefined, {
    maximumFractionDigits: 2,
  });
}

var CACHE_DATA = {};

export var gutil = {
  isEmptyObj: function (obj) {
    return Object.keys(obj).length === 0;
  },
  jsonCopy: function (obj) {
    // 通过json转换的方式来拷贝对象
    return JSON.parse(JSON.stringify(obj));
  },
  deepClone(obj) {
    // 注意只能对于object来说，不能用于数组或其它
    return extend(true, null, obj);
  },
  shadowClone(obj) {
    // 注意只能对于object来说，不能用于数组或其它
    return extend(null, obj);
  },
  getNowFormatDate: function () {
    var date = new Date();
    // var year = date.getFullYear();
    var month = date.getMonth() + 1;
    var d = date.getDate();
    var hour = date.getHours();
    var minute = date.getMinutes();
    var second = date.getSeconds();
    if (month < 10) {
      month = '0' + month;
    }
    if (d < 10) {
      d = '0' + d;
    }
    if (hour < 10) {
      hour = '0' + hour;
    }
    if (minute < 10) {
      minute = '0' + hour;
    }
    if (second < 10) {
      second = '0' + second;
    }
    // return year + "-" + month + "-" + d + " " +hour + ":" + minute + ":" + second;
    return hour + ':' + minute + ':' + second;
  },
  getStoreItem: function (key, defaultValue = undefined) {
    // 获取本地存档的数据
    if (!LocalStorage.has(key)) return defaultValue;
    return LocalStorage.getItem(key);
  },
  setStoreItem: function (key, value) {
    // 设置本地存档的数据，目前都放localstorage，后续换到IndexDb中
    LocalStorage.set(key, value);
  },
  refreshTableSelectedRows: function (tableRef, newDataList) {
    // 需要传newDataList, 如果用table_ref.data的话，并不会实时生效
    let ret = [];
    for (let newData of newDataList) {
      let id = tableRef.getRowKey(newData);
      if (tableRef.isRowSelected(id)) ret.push(newData);
    }
    return ret;
  },
  setCacheData: function (key, value) {
    CACHE_DATA[key] = value;
  },
  getCacheData: function (key) {
    return CACHE_DATA[key];
  },
  getDictValueByKeys: function (data, keys, defv) {
    let ret = data;
    for (let k of keys) {
      if (!(k in ret)) return defv;
      ret = ret[k];
    }
    return ret;
  },
  setDictValueByKeys: function (data, keys, value) {
    let ret = data;
    let num = keys.length;
    for (let i = 0; i < num - 1; i++) {
      let k = keys[i];
      if (!(k in ret)) {
        ret[k] = {};
      }
      ret = ret[k];
    }
    ret[keys[num - 1]] = value;
  },
  formatTs(timeStamp) {
    return date.formatDate(timeStamp, 'YYYY-MM-DD HH:mm:ss');
  },
  copyToClipboard(text) {
    return copyToClipboard(text);
  },
  routeFullUrl(routeLocate) {
    let router = gutil.getCacheData('router');
    let href = router.resolve(routeLocate).href;
    let origin = window.location.origin;
    if (href.startsWith(origin)) return href;
    else if (href.startsWith('/')) return `${origin}${href}`;
    else return `${origin}/${href}`;
  },
  strip(x, characters, left, right) {
    var start = 0;
    while (left && characters.indexOf(x[start]) >= 0) {
      start += 1;
    }
    var end = x.length - 1;
    while (right && characters.indexOf(x[end]) >= 0) {
      end -= 1;
    }
    return x.substr(start, end - start + 1);
  },
  authCheck(resource, authData) {
    let ret = authData[resource];
    if (ret === undefined) {
      // 为了减少发送量，客户端并不要求权限一定要存在
      // 未指定的，使用资源的默认值
      ret = resource.split(':')[0].endsWith('@w');
    }
    return ret;
  },
};

// -------------------- aes 加密 ---------------------------

export class MyAes {
  constructor(stringKey) {
    let hashVal = CryptoJS.SHA256(stringKey).toString();
    let key = hashVal.substring(8, 24);
    let iv = hashVal.substring(40, 56);
    this.key = CryptoJS.enc.Utf8.parse(key);
    this.iv = CryptoJS.enc.Utf8.parse(iv);
  }

  encrypt(data) {
    let encrypted = CryptoJS.AES.encrypt(data, this.key, {
      iv: this.iv,
      mode: CryptoJS.mode.CBC,
      padding: CryptoJS.pad.Pkcs7,
    });
    return encrypted.ciphertext.toString(CryptoJS.enc.Base64);
  }

  decrypt(data) {
    let bytes = CryptoJS.AES.decrypt(data, this.key, {
      iv: this.iv,
      mode: CryptoJS.mode.CBC,
      padding: CryptoJS.pad.Pkcs7,
    });
    return bytes.toString(CryptoJS.enc.Utf8);
  }
}

export const MikCall = {
  // 将常用的http请求封装于此

  sendTips(message, msgType = 'info') {
    sendTips(message, msgType);
  },

  sendSuccTips(message) {
    sendTips(message, 'positive');
  },

  sendErrorTips(message) {
    sendTips(message, 'negative');
  },

  failRet(why, code = '', ext = undefined) {
    // code 错误代码，用固定的字符串会比较好理解一些
    // ext 传其它的信息，一般用dict
    let ret = { why, code, suc: false };
    if (ext) ret['ext'] = ext;
    return ret;
  },

  makeConfirm(message, cb) {
    Dialog.create({
      // title: 'Confirm',
      message: message,
      cancel: true,
      persistent: true,
    }).onOk(() => {
      cb();
    });
  },

  makePrompt(message, defaultValue, cb, title) {
    Dialog.create({
      title: title,
      message: message,
      prompt: {
        model: defaultValue || '',
        type: 'text',
      },
      cancel: true,
      persistent: true,
    }).onOk((data) => {
      cb(data);
    });
  },

  coMakeConfirm(message) {
    return new Promise((resovle) => {
      Dialog.create({
        // title: 'Confirm',
        message: message,
        cancel: true,
        persistent: true,
      })
        .onOk(() => {
          resovle(true);
        })
        .onCancel(() => {
          resovle(false);
        });
    });
  },

  coMakePrompt(message, defaultValue, title) {
    return MikCall.coCreateDialog({
      title: title,
      message: message,
      prompt: {
        model: defaultValue || '',
        type: 'text',
      },
      cancel: true,
      persistent: true,
    });
  },

  coCreateDialog(opts) {
    return new Promise((resovle) => {
      Dialog.create(opts)
        .onOk((data) => {
          resovle([true, data]);
        })
        .onCancel(() => {
          resovle([false, null]);
        });
    });
  },

  sucRet(ret) {
    return { suc: true, ret: ret };
  },

  genUrlWithParam(url, param) {
    if (!param) return url;
    let queryString = Qs.stringify(param);
    if (queryString.length <= 0) return url;
    return `${url}?${queryString}`;
  },

  filterResponseResult(res) {
    // 服务端使用 sucRet, failRet 返回的话，可以使用该函数返回一个promise
    let data = res.data;
    if (data.suc) {
      return Promise.resolve(data);
    } else if (data.suc == undefined) {
      return Promise.reject(
        MikCall.failRet('返回数据格式不正确，必须包含suc，请检查!')
      );
    } else {
      return Promise.reject(data);
    }
  },

  catchResponseError(error) {
    let errMsg = error.message;
    if (error.response) {
      let errData = error.response.data;
      if (errData && errData.suc === false && errData.why) {
        return Promise.reject(errData);
      }
      let statusCode = error.response.status;
      let statusText = error.response.statusText;
      let url = error.response.config.url;
      errMsg = `(${statusCode}:${statusText}) in ${url}`;
    }
    let ret = MikCall.failRet(errMsg);
    return Promise.reject(ret);
  },

  alertRespErrMsg(data) {
    if (data.why) {
      MikCall.sendErrorTips(data.why);
    } else {
      // 否则继续promise调用链出去
      return Promise.reject(data);
    }
  },

  mikaxios(type, url, data, conf) {
    // Post 方法遵循 json请求，返回{suc, ret, why}的格式
    // 非raw请求的话，默认会处理一下http请求错误，根据suc进行处理
    // 返回请求的promise
    conf = conf || {};
    let { raw, convErr } = conf;
    let reqObj = {
      method: type,
      data: data,
      url: url,
      withCredentials: true,
    };
    let ret = axios(reqObj);
    if (!raw) {
      ret = ret.then(
        MikCall.filterResponseResult,
        MikCall.catchResponseError
      );
    }
    if (convErr) {
      ret = convErr(ret);
    }
    return ret;
  },

  mcpost(url, data, extraConf) {
    extraConf = extraConf || {};
    if (!extraConf.convErr) extraConf.convErr = MikCall.filterPromise;
    return MikCall.mikaxios('POST', url, data, extraConf);
  },

  mcget(url, param, extraConf) {
    extraConf = extraConf || {};
    if (!extraConf.convErr) extraConf.convErr = MikCall.filterPromise;
    return MikCall.mikaxios('GET', url, param, extraConf);
  },

  filterPromise(promise) {
    // 转换成 err, data的形式
    return promise.then((data) => data).catch((err) => err);
  },

  coDelay(t) {
    return new Promise((resovle) => {
      setTimeout(() => {
        resovle(true);
      }, t);
    });
  },
};
