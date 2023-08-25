import cookies from 'react-cookies';

/**
 * 操作cookie
 */
const cookie = {
  getToken(tokenKey: string) {
    return cookies.load(tokenKey);
  },
  clearToken(tokenKey: string) {
    cookies.remove(tokenKey, {
      domain: document.domain,
      path: '/'
    });
  },
  saveToken(tokenKey: string, tokenValue: string) {
    cookies.save(tokenKey, tokenValue, {
      domain: document.domain,
    });
  },
  updateToken(tokenKey: string, tokenValue: string) {
    this.clearToken(tokenKey);
    this.saveToken(tokenKey, tokenValue);
  },
};
export default cookie;
